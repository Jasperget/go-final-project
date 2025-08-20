package api

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// DateFormat - это константа для формата даты, используемого во всем API.
const DateFormat = "20060102"

// NextDate вычисляет следующую дату на основе правила повторения.
func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не может быть пустым для вычисления")
	}

	startDate, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %w", err)
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	// Обнуляем время для корректного сравнения дат
	now = now.Truncate(24 * time.Hour)
	next := startDate

	// Бесконечный цикл, который прерывается, когда найдена корректная следующая дата.
	for {
		// Сначала всегда сдвигаем дату на один интервал.
		switch rule {
		case "y":
			next = next.AddDate(1, 0, 0)
		case "d":
			if len(parts) < 2 {
				return "", fmt.Errorf("не указан интервал в днях для правила 'd'")
			}
			days, err := strconv.Atoi(parts[1])
			if err != nil || days <= 0 || days > 400 {
				return "", fmt.Errorf("неверный интервал дней: %s. Должно быть число от 1 до 400", parts[1])
			}
			next = next.AddDate(0, 0, days)
		case "w":
			if len(parts) < 2 {
				return "", fmt.Errorf("не указаны дни недели для правила 'w'")
			}
			daysOfWeekStr := strings.Split(parts[1], ",")
			var daysOfWeek []int
			for _, dayStr := range daysOfWeekStr {
				day, err := strconv.Atoi(dayStr)
				if err != nil || day < 1 || day > 7 {
					return "", fmt.Errorf("неверный день недели: %s", dayStr)
				}
				daysOfWeek = append(daysOfWeek, day)
			}
			sort.Ints(daysOfWeek) // Сортируем для удобного поиска

			// Преобразуем день недели Go (Sun=0) в наш формат (Mon=1..Sun=7)
			currentWeekday := int(next.Weekday())
			if currentWeekday == 0 {
				currentWeekday = 7
			}

			// Ищем следующий подходящий день недели
			nextDayIndex := -1
			for i, day := range daysOfWeek {
				if day > currentWeekday {
					nextDayIndex = i
					break
				}
			}

			var daysToAdd int
			if nextDayIndex != -1 {
				// Следующий день в этой же неделе
				daysToAdd = daysOfWeek[nextDayIndex] - currentWeekday
			} else {
				// Следующий день на следующей неделе
				daysToAdd = (7 - currentWeekday) + daysOfWeek[0]
			}
			next = next.AddDate(0, 0, daysToAdd)

		default:
			return "", fmt.Errorf("неподдерживаемый формат правила: %s", rule)
		}

		// Если полученная дата оказалась после точки отсчета 'now', то мы нашли то, что нужно.
		if next.After(now) {
			return next.Format(DateFormat), nil
		}
	}
}
