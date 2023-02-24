package model

import v1 "k8s.io/api/core/v1"

const (
	PVC_STATUS_PENDING = iota
	PVC_STATUS_BOUND
	PVC_STATUS_LOST
)

var (
	pvcStatusMap = map[v1.PersistentVolumeClaimPhase]int{
		v1.ClaimPending: PVC_STATUS_PENDING,
		v1.ClaimBound: PVC_STATUS_BOUND,
		v1.ClaimLost: PVC_STATUS_LOST,
	}
)

