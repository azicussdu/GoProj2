package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Course struct {
	ID          int     `gorm:"primaryKey;autoIncrement" json:"id"` // primaryKey;autoIncrement - kerek emes a tak id bolsa
	Title       string  `gorm:"type:varchar(255);not null" json:"title"`
	Description *string `json:"description,omitempty"`
	Slug        string  `gorm:"type:varchar(255);unique;not null" json:"slug"`
	Price       int     `gorm:"default:0;not null" json:"price"`
	Duration    int     `gorm:"default:0;not null" json:"duration"`
	Level       *string `gorm:"type:varchar(50)" json:"level,omitempty"`
	IsActive    bool    `gorm:"default:false;not null" json:"is_active"`

	TeacherID int  `gorm:"not null" json:"teacher_id"`
	Teacher   User `gorm:"foreignKey:TeacherID;constraint:OnDelete:RESTRICT" json:"teacher"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
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
