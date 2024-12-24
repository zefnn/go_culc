# О проекте
Этот проект представляет собой HTTP-сервер, написанный на языке Go, который позволяет вычислять математические выражения, переданные в формате JSON. Сервер поддерживает обработку POST-запросов и возвращает результат вычислений или сообщение об ошибке, если выражение некорректно.
___

# Описание проекта:
## Проект состоит из следующих компонентов:

+ ### HTTP-сервер:

Слушает запросы на указанном порту (по умолчанию 8080) и обрабатывает маршрут /api/v1/calculate.

+ ### Обработка выражений: 

Сервер принимает JSON-запросы с полем expression, вычисляет значение выражения и возвращает результат.

+ ### Обработка ошибок: 
Если выражение некорректно, сервер возвращает соответствующий код ошибки и описание проблемы.

+ ### Логирование:
Все запросы и ошибки логируются в файл log.txt.

+ ### Ошибки:

Имеется отдельный пакет для возможных ошибок
___
# Инструкции по запуску:

### Команда для запуска:
```
$env:PORT="8080"; go run ./cmd/main.go
```

Вы можете выбрать порт самостоятельно, заменив 8080 на ваше значение
___
# **Примеры использования**
Запустите сервер и скопируйте команду в cmd

1) Status 200 OK
```
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+2*2\"}" http://localhost:8080/api/v1/calculate
```
```
{"result":6}
```

2) Status 422 Unprocessable Entity
```
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+2*\"}" http://localhost:8080/api/v1/calculate
```
```
{"error":"Expression is not valid: not enough operands"}
```

3) Status 422 Unprocessable Entity
```
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+2)\"}" http://localhost:8080/api/v1/calculate
```
```
{"error":"Expression is not valid: unbalanced brackets"}
```


4) Status 500 Internal Server Error
```
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"20+15a\"}" http://localhost:8080/api/v1/calculate
```
```
{"error":"Internal server error: expression is not valid: unsupported symbol"}
```


5) Status 422 Unprocessable Entity
```
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"20/0\"}" http://localhost:8080/api/v1/calculate
```
```
{"error":"Expression is not valid: division by zero"}
```

