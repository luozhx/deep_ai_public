package dataset

import (
	"deep-ai-server/app/model"
	"deep-ai-server/app/tools/client"
	"deep-ai-server/app/tools/file"
	"fmt"
	core "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"mime/multipart"
	"path/filepath"
	"strconv"
)

func createPVC(dataset *model.Dataset, c chan string, stopListener chan struct{}) {
	cli := client.NewClient(client.NAMESPACE_DEEP_AI)
	informer := cli.GetInformerFactory(
		fmt.Sprintf(
			"app=dataset,id=%d,user=%d",
			dataset.ID,
			dataset.UserID,
		)).Core().V1().PersistentVolumeClaims().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    nil,
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPVC := oldObj.(*core.PersistentVolumeClaim)
			newPVC := newObj.(*core.PersistentVolumeClaim)

			if oldPVC.Status.Phase == newPVC.Status.Phase {
				return
			}

			_ = dataset.UpdatePVCStatus(newPVC.Status.Phase)
			switch newPVC.Status.Phase {
			case core.ClaimPending:
				break
			case core.ClaimBound:
				close(stopListener)
				pvcPath, _ := cli.GetPVPath(newPVC.Name)
				_ = dataset.UpdatePVC(newPVC.Name, pvcPath)
				c <- "success"
				break
			case core.ClaimLost:
			default:
				c <- "fail"
			}
		},
		DeleteFunc: nil,
	})
	go informer.Run(stopListener)

	pvcName := fmt.Sprintf("user%d-dataset%d", dataset.UserID, dataset.ID)
	if err := cli.CreatePVC(pvcName, "10Gi", map[string]string{
		"app":  "dataset",
		"id":   strconv.Itoa(int(dataset.ID)),
		"user": strconv.Itoa(int(dataset.UserID)),
	}); err != nil {
		close(stopListener)
		_ = dataset.UpdatePVCStatus(core.ClaimLost)
	}
}

func saveDataset(dataset *model.Dataset, fileData *multipart.FileHeader, c chan string) {
	status := <- c
	defer close(c)
	if status == "fail" {
		_ = dataset.UpdateDatasetStatus(model.DATASET_STATUS_UNAVAILABLE)
		return
	}
	path := filepath.Join(dataset.PersistentVolumePath, fmt.Sprintf("%s.zip", dataset.Name))
	if err := file.Save(fileData, path); err != nil {
		_ = dataset.UpdateDatasetStatus(model.DATASET_STATUS_UNAVAILABLE)
		return
	}
	if err := file.Unzip(path, dataset.PersistentVolumePath); err != nil {
		_ = dataset.UpdateDatasetStatus(model.DATASET_STATUS_UNAVAILABLE)
		return
	}
	_ = dataset.UpdateDatasetStatus(model.DATASET_STATUS_IDLE)
}
