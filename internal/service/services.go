package service

type Services struct {
	Course     *CourseService
	Lesson     *LessonService
	Enrollment *EnrollmentService
	Auth       *AuthService
}
