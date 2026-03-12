alter table enrollments
    add constraint fk_enrollments_user
        foreign key (user_id) references users(id)
            on delete cascade;

alter table enrollments
    add constraint fk_enrollments_course
        foreign key (course_id) references courses(id)
            on delete cascade;

alter table enrollments
    add constraint unique_user_course
        unique (user_id, course_id);
