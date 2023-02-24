package model

import (
	"deep-ai-server/app/db"
	v1 "k8s.io/api/core/v1"
	"time"
)

type Dataset struct {
	// gorm.Job
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Name        string
	Description string
	Size        int64

	Status    int
	PVCStatus int

	PersistentVolumeClaimName string `json:"-"`
	PersistentVolumePath      string `json:"-"`

	UserID uint `json:"-"`
}

const (
	DATASET_STATUS_CREATING = iota
	DATASET_STATUS_UNAVAILABLE
	DATASET_STATUS_IDLE
	DATASET_STATUS_BUSY
)

func (dataset *Dataset) UpdatePVC(name string, path string) error {
	dataset.PersistentVolumePath = path
	dataset.PersistentVolumeClaimName = name
	return db.DB().Model(dataset).Update(&dataset).Error
}

func (dataset Dataset) UpdatePVCStatus(phase v1.PersistentVolumeClaimPhase) error {
	return db.DB().Model(&dataset).Update(&Dataset{
		PVCStatus: pvcStatusMap[phase],
	}).Error
}

func (dataset Dataset) UpdateDatasetStatus(status int) error {
	return db.DB().Model(&dataset).Update(&Dataset{
		Status: status,
	}).Error
}
