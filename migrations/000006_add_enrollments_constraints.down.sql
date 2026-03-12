alter table enrollments
    drop constraint unique_user_course;

alter table enrollments
    drop constraint fk_enrollments_course;

alter table enrollments
    drop constraint fk_enrollments_user;