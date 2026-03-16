package models

import (
	"errors"
	"time"
)

type Enrollment struct {
	ID          int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID      int        `gorm:"column:user_id" json:"user_id"`
	CourseID    int        `gorm:"column:course_id" json:"course_id"`
	Progress    int        `gorm:"column:progress" json:"progress"`
	IsCompleted bool       `gorm:"column:is_completed" json:"is_completed"`
	EnrolledAt  time.Time  `gorm:"column:enrolled_at;autoCreateTime" json:"enrolled_at"`
	CompletedAt *time.Time `gorm:"column:completed_at" json:"completed_at,omitempty"`
}

func (Enrollment) TableName() string {
	return "enrollments"
}

type MyCourse struct {
	CourseID    int        `gorm:"column:course_id" json:"course_id"`
	Title       string     `gorm:"column:title" json:"title"`
	Description *string    `gorm:"column:description" json:"description,omitempty"`
	Slug        string     `gorm:"column:slug" json:"slug"`
	Price       int        `gorm:"column:price" json:"price"`
	Duration    int        `gorm:"column:duration" json:"duration"`
	Level       *string    `gorm:"column:level" json:"level,omitempty"`
	IsActive    bool       `gorm:"column:is_active" json:"is_active"`
	TeacherID   int        `gorm:"column:teacher_id" json:"teacher_id"`
	Progress    int        `gorm:"column:progress" json:"progress"`
	IsCompleted bool       `gorm:"column:is_completed" json:"is_completed"`
	EnrolledAt  time.Time  `gorm:"column:enrolled_at" json:"enrolled_at"`
	CompletedAt *time.Time `gorm:"column:completed_at" json:"completed_at,omitempty"`
}

type CreateEnrollment struct {
	UserID      int        `json:"user_id" binding:"required"`
	CourseID    int        `json:"course_id" binding:"required"`
	Progress    int        `json:"progress"`
	IsCompleted bool       `json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at"`
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
	Progress    *int       `json:"progress"`
	IsCompleted *bool      `json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at"`
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
