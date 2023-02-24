package model

import (
	"deep-ai-server/app/db"
	"time"
)

type Job struct {
	// gorm.Job
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Name             string
	Description      string
	Status           int
	Framework        string
	FrameworkVersion string
	EntryPoint       string `json:"-"`
	Args             string `json:"-"`

	Num    int    `json:"-"`
	CPU    int    `json:"-"`
	Memory string `json:"-"`
	GPU    int    `json:"-"`

	// foreign key
	UserID    uint `json:"-"`
	CodeID    uint `json:"-"`
	DatasetID uint `json:"-"`
	AIModelID uint `json:"-"`
}

const (
	JOB_STATUS_CREATING = iota
	JOB_STATUS_TRAINING
	JOB_STATUS_FINISHED
	JOB_STATUS_STOPPED
	JOB_STATUS_FAILED
)

const (
	JOB_FRAMEWORK_TENSORFLOW = "tensorflow"
	JOB_FRAMEWORK_PYTORCH    = "pytorch"
)

func (job Job) UpdateJobStatus(status int) error {
	return db.DB().Model(&job).Update(&Job{
		Status: status,
	}).Error
}
