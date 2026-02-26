package models

import "errors"

var (
	ErrCourseNotFound    = errors.New("course not found error")
	ErrSlugAlreadyExists = errors.New("this slug is already exists in courses")
	ErrLessonNotFound    = errors.New("lesson not found error")
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTeacherNotFound    = errors.New("no teacher with this ID error")
	ErrUserNotFound       = errors.New("user with email not found error")
)
