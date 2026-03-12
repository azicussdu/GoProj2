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