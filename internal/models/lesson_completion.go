package models

import "time"

type LessonCompletion struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	LessonID    int       `db:"lesson_id"`
	CompletedAt time.Time `db:"completed_at"`
}
