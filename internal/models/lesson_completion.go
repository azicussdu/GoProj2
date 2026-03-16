package models

import "time"

type LessonCompletion struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement"`
	UserID      int       `gorm:"column:user_id"`
	LessonID    int       `gorm:"column:lesson_id"`
	CompletedAt time.Time `gorm:"column:completed_at"`
}

func (LessonCompletion) TableName() string {
	return "lesson_completions"
}
