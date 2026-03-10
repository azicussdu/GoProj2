package models

import (
	"errors"
	"time"
)

type Enrollment struct {
	ID          int        `db:"id" json:"id"`
	UserID      int        `db:"user_id" json:"user_id"`
	CourseID    int        `db:"course_id" json:"course_id"`
	Progress    int        `db:"progress" json:"progress"`
	IsCompleted bool       `db:"is_completed" json:"is_completed"`
	EnrolledAt  time.Time  `db:"enrolled_at" json:"enrolled_at"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at,omitempty"`
}

type MyCourse struct {
	CourseID    int        `db:"course_id" json:"course_id"`
	Title       string     `db:"title" json:"title"`
	Description *string    `db:"description" json:"description,omitempty"`
	Slug        string     `db:"slug" json:"slug"`
	Price       int        `db:"price" json:"price"`
	Duration    int        `db:"duration" json:"duration"`
	Level       *string    `db:"level" json:"level,omitempty"`
	IsActive    bool       `db:"is_active" json:"is_active"`
	TeacherID   int        `db:"teacher_id" json:"teacher_id"`
	Progress    int        `db:"progress" json:"progress"`
	IsCompleted bool       `db:"is_completed" json:"is_completed"`
	EnrolledAt  time.Time  `db:"enrolled_at" json:"enrolled_at"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at,omitempty"`
}

type CreateEnrollment struct {
	UserID      int        `db:"user_id" json:"user_id" binding:"required"`
	CourseID    int        `db:"course_id" json:"course_id" binding:"required"`
	Progress    int        `db:"progress" json:"progress"`
	IsCompleted bool       `db:"is_completed" json:"is_completed"`
	EnrolledAt  time.Time  `db:"enrolled_at" json:"enrolled_at"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at"`
}

func (c *CreateEnrollment) Validate() error {
	if c.UserID <= 0 {
		return errors.New("invalid user id")
	}

	if c.CourseID <= 0 {
		return errors.New("invalid course id")
	}

	if c.Progress < 0 || c.Progress > 100 {
		return errors.New("progress must be between 0 and 100")
	}

	return nil
}

type UpdateEnrollment struct {
	Progress    *int       `db:"progress" json:"progress"`
	IsCompleted *bool      `db:"is_completed" json:"is_completed"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at"`
}

func (u *UpdateEnrollment) Validate() error {
	if u.Progress == nil && u.IsCompleted == nil && u.CompletedAt == nil {
		return errors.New("no fields provided for update")
	}

	if u.Progress != nil && (*u.Progress < 0 || *u.Progress > 100) {
		return errors.New("progress must be between 0 and 100")
	}

	return nil
}
