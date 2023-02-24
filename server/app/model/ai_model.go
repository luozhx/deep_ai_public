package model

import (
	"deep-ai-server/app/db"
	v1 "k8s.io/api/core/v1"
	"time"
)

type AIModel struct {
	// gorm.Job
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Name             string
	Description      string
	Framework        string
	FrameworkVersion string
	Network          string `sql:"type:text;"`

	Status           int
	PVCStatus        int

	PersistentVolumeClaimName string `json:"-"`
	PersistentVolumePath      string `json:"-"`

	// foreign key
	UserID uint `json:"-"`
}

const (
	AIMODEL_STATUS_CREATING = iota
	AIMODEL_STATUS_UNAVAILABLE
	AIMODEL_STATUS_IDLE
	AIMODEL_STATUS_BUSY
)

const (
	AIMODEL_FRAMEWORK_TENSORFLOW = "tensorflow"
	AIMODEL_FRAMEWORK_PYTORCH    = "pytorch"
)

func (aiModel *AIModel) UpdatePVC(name string, path string) error {
	aiModel.PersistentVolumePath = path
	aiModel.PersistentVolumeClaimName = name
	return db.DB().Model(aiModel).Update(&aiModel).Error
}

func (aiModel AIModel) UpdatePVCStatus(phase v1.PersistentVolumeClaimPhase) error {
	return db.DB().Model(&aiModel).Update(&AIModel{
		PVCStatus: pvcStatusMap[phase],
	}).Error
}

func (aiModel AIModel) UpdateAIModelStatus(status int) error {
	return db.DB().Model(&aiModel).Update(&AIModel{
		Status: status,
	}).Error
}

