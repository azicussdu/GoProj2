package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Lesson struct {
	ID int `json:"id"`

	CourseID int    `gorm:"not null;index" json:"course_id"`
	Course   Course `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`

	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Content   *string        `json:"content,omitempty"`
	VideoURL  *string        `json:"video_url,omitempty"`
	Duration  int            `gorm:"default:0;not null" json:"duration"`
	Position  int            `gorm:"default:0;not null" json:"position"`
	IsPreview bool           `gorm:"default:false;not null" json:"is_preview"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Lesson) TableName() string {
	return "lessons"
}

type CreateLesson struct {
	CourseID  int     `json:"course_id" binding:"required"`
	Title     string  `json:"title" binding:"required"`
	Content   *string `json:"content"`
	VideoURL  *string `json:"video_url"`
	Duration  int     `json:"duration"`
	Position  int     `json:"position"`
	IsPreview bool    `json:"is_preview"`
}

func (c *CreateLesson) Validate() error {
	if c.CourseID <= 0 {
		return errors.New("invalid course id")
	}

	if strings.TrimSpace(c.Title) == "" {
		return errors.New("lesson title is required")
	}

	if c.Content == nil && c.VideoURL == nil {
		return errors.New("lesson must contain content or video")
	}
	return nil
}

type UpdateLesson struct {
	CourseID  *int    `json:"course_id"`
	Title     *string `json:"title"`
	Content   *string `json:"content"`
	VideoURL  *string `json:"video_url"`
	Duration  *int    `json:"duration"`
	Position  *int    `json:"position"`
	IsPreview *bool   `json:"is_preview"`
}
