package tests

// FullNextDate включает проверку правил для недель и месяцев
// const FullNextDate = true

var Port = 7540
var DBFile = "../scheduler.db"

var FullNextDate = false

// Search включает тесты для поиска задач.
var Search = true

// Token для аутентификации.
// ПЕРЕД ЗАПУСКОМ ТЕСТОВ:
// 1. Запустите сервер с установленным TODO_PASSWORD.
// 2. Выполните POST-запрос на /api/signin с вашим паролем.
// 3. Скопируйте полученный токен и вставьте его сюда.
var Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTU4MjUzODksInBhc3N3b3JkX2hhc2giOiI5Zjg2ZDA4MTg4NGM3ZDY1OWEyZmVhYTBjNTVhZDAxNWEzYmY0ZjFiMmIwYjgyMmNkMTVkNmMxNWIwZjAwYTA4In0.O7poQ8EUe2lMExmGzAjCWWkp0LmgeP9Xr8t3T6QMoCw"
