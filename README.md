# Калькулятор
Этот проект представляет собой веб-сервис для вычисления математических выражений. Сервис поддерживает базовые арифметические операции (сложение, вычитание, умножение, деление) и обработку выражений с использованием скобок.

### Требования
GO
Git

### Запуск
```cmd
git clone https://github.com/kupetss/calculator_golangV3
cd calculator_golangV3
go mod tidy
cd main
go run main.go
```
##### Приложение поддерживает post запросы с json формата {"expression": "ваше выражение"}
### Пример работы для Power Shell:
```bash
$response = Invoke-WebRequest -Uri http://localhost:8080/api/v1/calculate -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"expression": "2 + 2 * (3 - 1)"}'
$jsonResponse = $response.Content | ConvertFrom-Json
$jsonResponse.id
```
На выходе получите ID
```
ID
```
Для того что бы узнать овет:
```
$response = Invoke-WebRequest -Uri http://localhost:8080/api/v1/expressions/ID
$response.Content
```
```
Результат
```

##### Пример
Запрос
```
PS C:\Users\kupets> $response = Invoke-WebRequest -Uri http://localhost:8080/api/v1/calculate -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"expression": "2 + 2 * (3 - 1)"}'
>> $jsonResponse = $response.Content | ConvertFrom-Json
>> $jsonResponse.id 
```
ID:
```
501445037495581358
```
Для того что бы узнать овет:
```
$response = Invoke-WebRequest -Uri http://localhost:8080/api/v1/expressions/501445037495581358
$response.Content
```
Результат:
```
{"res":{"id":501445037495581358,"status":"ok","result":6}}
```

### Пример работы для cmd:
Запрос:
```
curl -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d "{\"expression\": \"2 + 2 * (3 - 1)\"}"
```
На выходе получаем ID
```
ID
```
Для того что бы получить ответ:
```
curl -X GET http://localhost:8080/api/v1/expressions/ID
```
```
Резутат
```

##### Пример
Запрос:
```
curl -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d "{\"expression\": \"2 + 2 * (3 - 1)\"}"
```
На выходе получаем ID
```
{"id":2559269320730327156}
```
Для того что бы получить ответ:
```
curl -X GET http://localhost:8080/api/v1/expressions/2559269320730327156
```
```
{"res":{"id":2559269320730327156,"status":"ok","result":6}}
```
