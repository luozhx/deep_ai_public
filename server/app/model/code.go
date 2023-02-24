package model

import (
	"deep-ai-server/app/db"
	v1 "k8s.io/api/core/v1"
	"time"
)

type Code struct {
	// gorm.Job
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Name             string
	Description      string
	Framework        string
	FrameworkVersion string

	Status         int
	PVCStatus      int
	NotebookStatus int

	ServiceName               string
	PersistentVolumeClaimName string `json:"-"`
	PersistentVolumePath      string `json:"-"`

	UserID uint `json:"-"`
}

const (
	CODE_STATUS_UNAVAILABLE = iota
	CODE_STATUS_IDLE
	CODE_STATUS_BUSY
)

const (
	CODE_NOTEBOOK_STATUS_PROCESSING = iota
	CODE_NOTEBOOK_STATUS_AVAILABLE
	CODE_NOTEBOOK_STATUS_UNAVAILABLE
)

const (
	CODE_FRAMEWORK_TENSORFLOW = "tensorflow"
	CODE_FRAMEWORK_PYTORCH    = "pytorch"
)

func (code Code) UpdatePVC(name string, path string) error {
	return db.DB().Model(&code).Update(&Code{
		PersistentVolumePath:      path,
		PersistentVolumeClaimName: name,
	}).Error
}

func (code Code) UpdatePVCStatus(phase v1.PersistentVolumeClaimPhase) error {
	return db.DB().Model(&code).Update(&Code{
		PVCStatus: pvcStatusMap[phase],
	}).Error
}

func (code Code) UpdateNotebookStatus(status int) error {
	return db.DB().Model(&code).Update(&Code{
		NotebookStatus: status,
	}).Error
}

func (code Code) UpdateCodeStatus(status int) error {
	return db.DB().Model(&code).Update(&Code{
		Status: status,
	}).Error
}
