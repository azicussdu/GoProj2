package models

import "time"

type Lesson struct {
	ID        int     `db:"id" json:"id"`
	CourseID  int     `db:"course_id" json:"course_id"`
	Title     string  `db:"title" json:"title"`
	Content   *string `db:"content" json:"content,omitempty"`
	VideoURL  *string `db:"video_url" json:"video_url,omitempty"`
	Duration  int     `db:"duration" json:"duration"`
	Position  int     `db:"position" json:"position"`
	IsPreview bool    `db:"is_preview" json:"is_preview"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateLesson struct {
	CourseID  int     `db:"course_id" json:"course_id" binding:"required"`
	Title     string  `db:"title" json:"title" binding:"required"`
	Content   *string `db:"content" json:"content"`
	VideoURL  *string `db:"video_url" json:"video_url"`
	Duration  int     `db:"duration" json:"duration"`
	Position  int     `db:"position" json:"position"`
	IsPreview bool    `db:"is_preview" json:"is_preview"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UpdateLesson struct {
	CourseID  *int    `db:"course_id" json:"course_id"`
	Title     *string `db:"title" json:"title"`
	Content   *string `db:"content" json:"content"`
	VideoURL  *string `db:"video_url" json:"video_url"`
	Duration  *int    `db:"duration" json:"duration"`
	Position  *int    `db:"position" json:"position"`
	IsPreview *bool   `db:"is_preview" json:"is_preview"`
}
