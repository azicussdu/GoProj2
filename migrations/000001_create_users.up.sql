create type user_role as enum ('student', 'teacher', 'admin');

create table users(
    id              serial primary key,
    full_name       varchar(255) not null,
    email           varchar(255) not null unique,
    password_hash   text not null,
    role            user_role default 'student'::user_role not null,
    is_active       boolean   default true                 not null,
    created_at      timestamp default now()                not null,
    updated_at      timestamp default now()                not null
);