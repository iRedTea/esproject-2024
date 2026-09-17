# EasyStartupProject

Небольшой Go-сервис для управления проектами на сайте. Написан в 2024, как pet-проект, код открыт в 2026 (обновлен только Dockerfile и README! Остальной код сохранен так, как он выглядел в 2024! Никакого рефакторинга с тех пор не производилось!).

### Технологии
- Go
- Rest API
- Gin
- JWT
- SQL (MySQL)
- Swagger

### Быстрый старт

Требования: `Go >= 1.22`, Makefile (опционально), Docker (опционально).

Сборка локально:

```bash
go build -o esproject ./cmd/release
./esproject
```


### Переменные окружения

- `ADMIN_TOKEN` - токен администратора (используется в `pkg/service/auth.go`).
- `DB_PASSWORD` - пароль для ДБ.

Все остальные параметры хранятся в configs/config.yml

### API Endpoints
- `GET /api/project/:id` : Получить проект по ID
- `DELETE /api/project/:id` : Удалить проект по ID (только владелец)
- `POST /api/project/:id` : Редактировать проект по ID (только владелец)
- `GET /api/project/all?limit={}&start={}` : Получить список ID проектов (параметры: `limit`, `start`)
- `GET /api/project/count` : Получить количество проектов
- `GET /api/project/my?limit={}` : Получить ID проектов, где пользователь является участником (параметр: `limit`)
- `PUT /api/project` : Создать новый проект (JSON в теле запроса)
- `GET /api/member/invite/:id` : Получить приглашение по ID
- `GET /api/member/invite/my` : Получить приглашения текущего пользователя
- `POST /api/member/invite/accept/:id` : Принять приглашение
- `POST /api/member/invite` : Пригласить пользователя в проект (только владелец)
- `POST /api/member/kick` : Исключить участника из проекта (только владелец)
- `GET /swagger/*any` : Swagger UI и документация

### Конфигурация
- Для запуска нужно прописать конфигурационный файл configs/config.yml. Пример ключевых полей:

```yaml
port: "8011"
db:
	host: "..."
	port: "3306"
	username: "..."
	dbname: "..."
```
