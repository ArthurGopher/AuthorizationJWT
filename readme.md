UPDATE.

1) Добавил генерацию файлов trace и profile через хэндлеры startProfilingHandler и stopProfilingHandler (мне было удобно оставить их в main.go 
и не выносить их в директорию handlers)

1.1. После запуска сервера мы прописываем curl -X POST http://localhost:8080/start-profiling \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIxMjMifQ.E3nZ2_CSpzB_K6yslDT7pGVuuBi9rHQSenFIhv1k-ZE" 
  (доступ защищен, поэтому нужен токен)

  Этот эндпоин запускае процесс профилирования 

1.2 далее, для нагрузки на сервер, мы можем, например, продублировать несколько раз: 

curl -X POST http://localhost:8080/api/address/search \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIxMjMifQ.E3nZ2_CSpzB_K6yslDT7pGVuuBi9rHQSenFIhv1k-ZE" \
  -H "Content-Type: application/json" \
  -d '{"query": "Moscow"}'

1.3 далее мы прерываем процесс профилирования через: 

curl -X POST http://localhost:8080/stop-profiling \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIxMjMifQ.E3nZ2_CSpzB_K6yslDT7pGVuuBi9rHQSenFIhv1k-ZE"

1.4 Для анализа профиля мы пишем команду go tool pprof cpu_profile.prof,
получаем примерно такой ответ: Type: cpu
Time: Jan 14, 2025 at 7:37pm (MSK)
Duration: 64.84s, Total samples = 30ms (0.046%)
Entering interactive mode (type "help" for commands, "o" for options)

1.5  go tool trace trace.out 

2025/01/14 19:40:18 Parsing trace...
2025/01/14 19:40:18 Splitting trace...
2025/01/14 19:40:18 Opening browser. Trace viewer is listening on http://127.0.0.1:55028

Эта команда откроет веб-интерфейс для анализа трассировки.


2) Возможные узкие места или проблемы в коде:
2.1 Глобальная переменная UserStorage может стать узким местом при высоком количестве запросов.
     Нужно использовать базу данных вместо памяти.

2.2 Отсутствует тайм-аут для HTTP-запросов, что может привести к зависанию сервера
Нужно добавить тай-ауты 

2.3 Некоторые ошибки в коде игнорируются и опускаются
Нужно добавить обработку ошибок

2.4 Ключ API жестко вшит в коде
Нужно вынести его в файл .env 




Old Readme
1 Сначало нужно зарегистрировать пользователя через POST http://localhost:8080/api/register и тела запроса, например: {"username": "user01",
    "password": "password321"
    }
2  Далее оставляем тело запроса и меняем маршрут на http://localhost:8080/api/login -> получаем токен. Вот рабочий токен: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXJvayJ9.U4XQ1s67AACewsZu83PzwNWti5EAu9aYyCDAKk7eaGs

3 Отправим нагрузку на сервер, например: curl -X POST http://localhost:8080/api/address/geocode -d '{"lat":"55.7558","lng":"37.6173"}' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXJvayJ9.U4XQ1s67AACewsZu83PzwNWti5EAu9aYyCDAKk7eaGs" -H "Content-Type: application/json"

4 Загружаем профиль. curl -X GET http://localhost:8080/mycustompath/pprof/profile -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXJvayJ9.U4XQ1s67AACewsZu83PzwNWti5EAu9aYyCDAKk7eaGs" -o cpu_profile

5 Открываем профиль (нужно установить Graphviz). go tool pprof -http=:8081 cpu_profile

6 Получаем визуализацию нашего профиля в виде графов вызовов функций


