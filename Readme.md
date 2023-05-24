Поднять postgres в Docker:
```bash
docker run --name postges -p 5432:5432 -e POSTGRES_PASSWORD=root -e POSTGRES_USER=root -d postgres:14
```
Запустить миграции с помощью golang-migrate/migrate:
```bash
migrate -path migrations -database "postgres://root:root@localhost/postgres?sslmode=disable" up
migrate -path migrations -database "postgres://root:root@localhost/postgres?sslmode=disable" down
```
