package models

import "errors"

var (
	ErrCourseNotFound          = errors.New("course not found error")
	ErrSlugAlreadyExists       = errors.New("this slug is already exists in courses")
	ErrLessonNotFound          = errors.New("lesson not found error")
	ErrLessonAlreadyCompleted  = errors.New("lesson already completed")
	ErrEnrollmentAlreadyExists = errors.New("user is already enrolled in this course")
	ErrOnlyStudentsCanEnroll   = errors.New("only students can enroll in courses")
	ErrEnrollmentNotFound      = errors.New("enrollment not found")
	ErrCourseCannotBeActivated = errors.New("course cannot be activated without lessons")
)

var (
	ErrUserAlreadyExists         = errors.New("user already exists")
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrTeacherNotFound           = errors.New("no teacher with this ID error")
	ErrUserNotFound              = errors.New("user with email not found error")
	ErrInvalidRole               = errors.New("invalid role")
	ErrRoleChangeOnlyFromStudent = errors.New("only student role can be changed")
)
