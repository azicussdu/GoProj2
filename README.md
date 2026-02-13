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

---

### ТАБЛИЦЫ

```sql
create type user_role as enum ('student', 'teacher', 'admin');

create table users (
                       id              serial primary key,
                       full_name       varchar(255) not null,
                       email           varchar(255) not null unique,
                       password_hash   text not null,
                       role            user_role not null default 'student',
                       is_active       boolean not null default true,
                       created_at      timestamp not null default now(),
                       updated_at      timestamp not null default now()
);

create table courses (
                         id            serial primary key,
                         title         varchar(255) not null,
                         description   text,
                         slug          varchar(255) not null unique,
                         price         integer not null default 0,
                         duration      integer not null default 0,
                         level         varchar(50),
                         is_active     boolean not null default false,
                         teacher_id    integer not null references users(id) on delete restrict,
                         created_at    timestamp not null default now(),
                         updated_at    timestamp not null default now(),
                         deleted_at    timestamp null
);

create table lessons (
                         id          serial primary key,
                         course_id   integer not null references courses(id) on delete cascade,
                         title       varchar(255) not null,
                         content     text,
                         video_url   text,
                         duration    integer not null default 0,
                         position    integer not null default 0, -- order of lesson in course
                         is_preview  boolean not null default false,
                         created_at  timestamp not null default now(),
                         updated_at  timestamp not null default now(),
                         deleted_at  timestamp null
);

create table enrollments (
                             id          serial primary key,
                             user_id     integer not null references users(id) on delete cascade,
                             course_id   integer not null references courses(id) on delete cascade,
                             enrolled_at timestamp not null default now(),

                             constraint unique_user_course
                                 unique (user_id, course_id)
);
```