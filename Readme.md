# Web-Сервис сотрудников

* POST   /employee :  добавление сотрудника
* DELETE /employee/{id} :  удаление сотрудника по его id
* GET    /company/{id}     :  вывод списĸа сотрудников для указанной компании
* GET    /company/{id}/department/{name} :  вывод списĸа сотрудников для указанного отдела компании
* PATCH  /employee/{id} :  редаĸтирование сотрудника по его ID

Пример команды для запуска контейнера с postgres в Docker:
```bash
docker run --name postges -p 5432:5432 -e POSTGRES_PASSWORD=root -e POSTGRES_USER=root -d postgres:14
```
Запустить миграции с помощью golang-migrate/migrate:
```bash
migrate -path migrations -database "postgres://root:root@localhost/postgres?sslmode=disable" up
migrate -path migrations -database "postgres://root:root@localhost/postgres?sslmode=disable" down
```
Запуск сервиса:
```bash
docker-compose -f docker-compose.yaml up --no-start
docker-compose -f docker-compose.yaml start
```