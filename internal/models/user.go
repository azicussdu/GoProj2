package models

import "time"

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleTeacher UserRole = "teacher"
	RoleStudent UserRole = "student"
)

func IsValidRole(role string) bool {
	switch UserRole(role) {
	case RoleAdmin, RoleTeacher, RoleStudent:
		return true
	}
	return false
}

type User struct {
	ID           int       `gorm:"type:serial" json:"type:id"`
	FullName     string    `gorm:"type:varchar(255);not null" json:"full_name"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         UserRole  `gorm:"type:user_role;default:'student';not null"`
	IsActive     bool      `gorm:"default:true;not null" json:"is_active"`
	CreatedAt    time.Time `gorm:"default:now();not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:now();not null" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type CreateUser struct {
	FullName     string `json:"full_name" binding:"required"`
	Email        string `json:"email" binding:"required"`
	PasswordHash string `json:"-"`
}

type RegisterUser struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type ChangeUserRoleInput struct {
	Role string `json:"role" binding:"required"`
}
