# tm-backend-trainee-impl-clean-template
Решение [тестового задания](https://github.com/avito-tech/tm-backend-trainee)

## Используемые библиотеки
Роутер Gin  
Postgres pgxpool

## Запуск
Подготовить файл ```.env```  
```docker-compose up --build -d```

## Запросы
Сохранение статистики
```console
curl -v -X POST http://127.0.0.1:8081/v1/statistics/save \ 
-d '{"date":"2006-01-21", "views": 2, "clicks":1, "cost": "5.25"}'
```


Получить статистику (Order параметр необязательный)
```console
curl -v -X POST http://127.0.0.1:8081/v1/statistics/get -d '{"from":"2006-01-21", "to": "2006-01-21"}'
```
Сброс статистики
```console
curl -v --request DELETE http://127.0.0.1:8081/v1/statistics/clear
```