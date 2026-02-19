package models

import (
	"errors"
	"strings"
	"time"
)

type Course struct {
	ID          int     `db:"id" json:"id"`
	Title       string  `db:"title" json:"title"`
	Description *string `db:"description" json:"description,omitempty"`
	Slug        string  `db:"slug" json:"slug"`
	Price       int     `db:"price" json:"price"`
	Duration    int     `db:"duration" json:"duration"`
	Level       *string `db:"level" json:"level,omitempty"`
	IsActive    bool    `db:"is_active" json:"is_active"`
	TeacherID   int     `db:"teacher_id" json:"teacher_id"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateCourse struct {
	Title       string  `db:"title" json:"title" binding:"required"`
	Description *string `db:"description" json:"description"`
	Slug        string  `db:"slug" json:"slug" binding:"required"`
	Price       int     `db:"price" json:"price"`
	Duration    int     `db:"duration" json:"duration"`
	Level       *string `db:"level" json:"level"`
	IsActive    bool    `db:"is_active" json:"is_active"`
	TeacherID   int     `db:"teacher_id" json:"teacher_id" binding:"required"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c *CreateCourse) Validate() error {
	if strings.TrimSpace(c.Title) == "" {
		return errors.New("course title is required")
	}

	if c.Price < 0 {
		return errors.New("invalid price")
	}
	return nil
}

type UpdateCourse struct {
	Title       *string `db:"title" json:"title"`
	Description *string `db:"description" json:"description"`
	Slug        *string `db:"slug" json:"slug"`
	Price       *int    `db:"price" json:"price"`
	Duration    *int    `db:"duration" json:"duration"`
	Level       *string `db:"level" json:"level"`
	IsActive    *bool   `db:"is_active" json:"is_active"`
	TeacherID   *int    `db:"teacher_id" json:"teacher_id"`
}

func (c *UpdateCourse) Validate() error {
	if c.Title == nil &&
		c.Description == nil &&
		c.Slug == nil &&
		c.Price == nil &&
		c.Duration == nil &&
		c.Level == nil &&
		c.IsActive == nil {
		return errors.New("no fields provided for update")
	}

	if c.Price != nil && *c.Price < 0 {
		return errors.New("price cannot be negative")
	}

	if c.Duration != nil && *c.Duration < 0 {
		return errors.New("duration cannot be negative")
	}

	if c.Slug != nil {
		slug := strings.ToLower(*c.Slug)
		c.Slug = &slug
	}
	return nil
}
