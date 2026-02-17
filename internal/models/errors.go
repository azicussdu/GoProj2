package models

import "errors"

var ErrTeacherNotFound = errors.New("no teacher with this ID error")

var ErrCourseNotFound = errors.New("course not found error")
var ErrSlugAlreadyExists = errors.New("this slug is already exists in courses")

var ErrLessonNotFound = errors.New("lesson not found error")
