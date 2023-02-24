package model

import (
	"time"
)

type User struct {
	// gorm.Job
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Username      string    `gorm:"not null;unique_index"`
	PasswordHash  string    `gorm:"not null" json:"-"`
	Email         string    `gorm:"not null;unique_index"`
	Role          int       `gorm:"not null"`
	LastLoginTime time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	Codes      []Code      `json:"-"`
	Datasets   []Dataset   `json:"-"`
	Jobs       []Job       `json:"-"`
	Inferences []Inference `json:"-"`
	AIModels   []AIModel   `json:"-"`
}

const (
	USER_ROLE_ADMIN = iota
	USER_ROLE_USER
)
