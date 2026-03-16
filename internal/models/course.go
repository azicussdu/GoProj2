package models

import (
	"errors"
	"strings"
	"time"
)

type Course struct {
	ID          int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title       string     `gorm:"column:title" json:"title"`
	Description *string    `gorm:"column:description" json:"description,omitempty"`
	Slug        string     `gorm:"column:slug" json:"slug"`
	Price       int        `gorm:"column:price" json:"price"`
	Duration    int        `gorm:"column:duration" json:"duration"`
	Level       *string    `gorm:"column:level" json:"level,omitempty"`
	IsActive    bool       `gorm:"column:is_active" json:"is_active"`
	TeacherID   int        `gorm:"column:teacher_id" json:"teacher_id"`
	Teacher     User       `gorm:"foreignKey:TeacherID;references:ID" json:"teacher"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
}

func (Course) TableName() string {
	return "courses"
}

type CreateCourse struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	Slug        string  `json:"slug" binding:"required"`
	Price       int     `json:"price"`
	Duration    int     `json:"duration"`
	Level       *string `json:"level"`
	IsActive    bool    `json:"is_active"`
	TeacherID   int     `json:"teacher_id" binding:"required"`
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
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Slug        *string `json:"slug"`
	Price       *int    `json:"price"`
	Duration    *int    `json:"duration"`
	Level       *string `json:"level"`
	IsActive    *bool   `json:"is_active"`
	TeacherID   *int    `json:"teacher_id"`
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
