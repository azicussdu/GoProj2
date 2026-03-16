package models

import (
	"errors"
	"strings"
	"time"
)

type Lesson struct {
	ID        int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CourseID  int        `gorm:"column:course_id" json:"course_id"`
	Title     string     `gorm:"column:title" json:"title"`
	Content   *string    `gorm:"column:content" json:"content,omitempty"`
	VideoURL  *string    `gorm:"column:video_url" json:"video_url,omitempty"`
	Duration  int        `gorm:"column:duration" json:"duration"`
	Position  int        `gorm:"column:position" json:"position"`
	IsPreview bool       `gorm:"column:is_preview" json:"is_preview"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
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
