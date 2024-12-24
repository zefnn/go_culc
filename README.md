#О проекте
Этот проект представляет собой HTTP-сервер, написанный на языке Go, который позволяет вычислять математические выражения, переданные в формате JSON. Сервер поддерживает обработку POST-запросов и возвращает результат вычислений или сообщение об ошибке, если выражение некорректно.
___

Описание проекта:
Проект состоит из следующих компонентов:

HTTP-сервер: Слушает запросы на указанном порту (по умолчанию 8080) и обрабатывает маршрут /api/v1/calculate.

Обработка выражений: Сервер принимает JSON-запросы с полем expression, вычисляет значение выражения и возвращает результат.

Обработка ошибок: Если выражение некорректно, сервер возвращает соответствующий код ошибки и описание проблемы.

Логирование: Все запросы и ошибки логируются в файл log.txt.

ИНСТРУКЦИЯ ПО ЗАПУСКУ:

**Команда для запуска : $env:PORT="8080"; go run ./cmd/main.go**
Вы можете выбрать порт самостоятельно, заменив 8080 на ваше значение

**Примеры использования**
Запустите сервер и скопируйте команду в cmd

1) curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+2*2\"}" http://localhost:8080/api/v1/calculate

{"result":6}
ststus 200 OK

2) curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+2*\"}" http://localhost:8080/api/v1/calculate

{"error":"Expression is not valid: not enough operands"}
status 422 Unprocessable Entity

3) curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+2)\"}" http://localhost:8080/api/v1/calculate

{"error":"Expression is not valid: unbalanced brackets"}
status 422 Unprocessable Entity

4) curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"20+15a\"}" http://localhost:8080/api/v1/calculate

{"error":"Internal server error: expression is not valid: unsupported symbol"}
status 500 Internal Server Error

5) curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"20/0\"}" http://localhost:8080/api/v1/calculate

{"error":"Expression is not valid: division by zero"}
422
Unprocessable Entity
