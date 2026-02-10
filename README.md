### ТАБЛИЦА COURSES

```sql
CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    slug VARCHAR(255) UNIQUE NOT NULL,
    price INTEGER NOT NULL DEFAULT 0,
    duration INT NOT NULL DEFAULT 0, -- minutes
    level VARCHAR(50), -- beginner, advanced
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    instructor_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

---

### ЭНДПОЙНТЫ ДЛЯ COURSES

**GET /courses** — Возвращает все курсы

**GET /courses/3** — Возвращает один курс по ID.

**POST /courses** — Создает новый курс.  
*Тело запроса:*
```json
{
  "title": "Python 2 Development",
  "description": "Learn how to build production-ready backend services using Python.",
  "slug": "python-backend-development",
  "price": 9900,
  "duration": 220,
  "level": "beginner",
  "is_active": true,
  "instructor_id": 2
}
```

**PUT /courses/3** — Обновляет курс по ID.  
*Пример тела запроса:*
```json
{
  "title": "Updated Go Backend Course",
  "price": 24900,
  "slug": "go-backend-pro"
}
```

**DELETE /courses/3** — Удаляет курс по ID.