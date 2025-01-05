Здравствуйте, проверил работу кода через Curl, все работает!

1 Сначало нужно зарегистрировать пользователя через POST http://localhost:8080/api/register и тела запроса, например: {"username": "user01",
    "password": "password321"
    }
2  Далее оставляем тело запроса и меняем маршрут на http://localhost:8080/api/login -> получаем токен. Вот рабочий токен: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXJvayJ9.U4XQ1s67AACewsZu83PzwNWti5EAu9aYyCDAKk7eaGs

3 Отправим нагрузку на сервер, например: curl -X POST http://localhost:8080/api/address/geocode -d '{"lat":"55.7558","lng":"37.6173"}' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXJvayJ9.U4XQ1s67AACewsZu83PzwNWti5EAu9aYyCDAKk7eaGs" -H "Content-Type: application/json"

4 Загружаем профиль. curl -X GET http://localhost:8080/mycustompath/pprof/profile -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXJvayJ9.U4XQ1s67AACewsZu83PzwNWti5EAu9aYyCDAKk7eaGs" -o cpu_profile

5 Открываем профиль (нужно установить Graphviz). go tool pprof -http=:8081 cpu_profile

6 Получаем визуализацию нашего профиля в виде графов вызовов функций


