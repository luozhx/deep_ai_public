package model

import (
	"deep-ai-server/app/db"
	"time"
)

type Inference struct {
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
	Dimension        string
	ModelFile        string

	//Num    int `json:"-"`
	//CPU    int `json:"-"`
	//Memory int `json:"-"`
	//GPU    int `json:"-"`

	ServiceName string

	UserID    uint `json:"-"`
	AIModelID uint `json:"-"`
}

const (
	INFERENCE_STATUS_PROCESSING = iota
	INFERENCE_STATUS_AVAILABLE
	INFERENCE_STATUS_UNAVAILABLE
)

const (
	INFERENCE_FRAMEWORK_TENSORFLOW = "tensorflow"
	INFERENCE_FRAMEWORK_PYTORCH    = "pytorch"
)

func (inference Inference) UpdateInferenceStatus(status int) error {
	return db.DB().Model(&inference).Update(&Inference{
		Status: status,
	}).Error
}
