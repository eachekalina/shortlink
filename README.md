# shortlink

[Описание архитектуры](arch/README.md)

## Инструкции

Для запуска используйте docker-compose.yaml. Для создания короткой ссылки отправьте запрос:

```shell
curl -d '{"link": "https://www.google.com/"}' localhost:8080/new
```

Придет ответ вида

```json
{"path":"http://localhost:8080/cAZHEKiXdoxN"}
```

При переходе по ссылке произойдет перенаправление на желаемую страницу.
