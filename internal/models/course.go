package models

import (
	"time"
)

type Course struct {
	ID           int     `db:"id" json:"id"`
	Title        string  `db:"title" json:"title"`
	Description  *string `db:"description" json:"description,omitempty"`
	Slug         string  `db:"slug" json:"slug"`
	Price        int     `db:"price" json:"price"`
	Duration     int     `db:"duration" json:"duration"`
	Level        *string `db:"level" json:"level,omitempty"`
	IsActive     bool    `db:"is_active" json:"is_active"`
	InstructorID int     `db:"instructor_id" json:"instructor_id"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateCourse struct {
	Title        string  `db:"title" json:"title" binding:"required"`
	Description  *string `db:"description" json:"description"`
	Slug         string  `db:"slug" json:"slug" binding:"required"`
	Price        int     `db:"price" json:"price"`
	Duration     int     `db:"duration" json:"duration"`
	Level        *string `db:"level" json:"level"`
	IsActive     bool    `db:"is_active" json:"is_active"`
	InstructorID int     `db:"instructor_id" json:"instructor_id" binding:"required"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UpdateCourse struct {
	Title        *string `db:"title" json:"title"`
	Description  *string `db:"description" json:"description"`
	Slug         *string `db:"slug" json:"slug"`
	Price        *int    `db:"price" json:"price"`
	Duration     *int    `db:"duration" json:"duration"`
	Level        *string `db:"level" json:"level"`
	IsActive     *bool   `db:"is_active" json:"is_active"`
	InstructorID *int    `db:"instructor_id" json:"instructor_id"`
}
