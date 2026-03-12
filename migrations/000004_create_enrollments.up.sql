create table enrollments (
    id          serial primary key,
    user_id     integer not null,
    course_id   integer not null,
    progress integer not null default 0,
    enrolled_at timestamp not null default now(),
    completed_at timestamp null
);