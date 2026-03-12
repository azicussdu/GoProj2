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