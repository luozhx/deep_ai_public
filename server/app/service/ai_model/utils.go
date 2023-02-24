package ai_model

import (
	"deep-ai-server/app/model"
	"deep-ai-server/app/tools/client"
	"fmt"
	core "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"os"
	"path/filepath"
	"strconv"
)

func createPVC(aiModel *model.AIModel, c chan string, stopListener chan struct{}) {
	cli := client.NewClient(client.NAMESPACE_DEEP_AI)
	informer := cli.GetInformerFactory(
		fmt.Sprintf(
			"app=model,id=%d,user=%d",
			aiModel.ID,
			aiModel.UserID,
		)).Core().V1().PersistentVolumeClaims().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: nil,
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPVC := oldObj.(*core.PersistentVolumeClaim)
			newPVC := newObj.(*core.PersistentVolumeClaim)

			if oldPVC.Status.Phase == newPVC.Status.Phase {
				return
			}

			_ = aiModel.UpdatePVCStatus(newPVC.Status.Phase)
			switch newPVC.Status.Phase {
			case core.ClaimPending:
				break
			case core.ClaimBound:
				close(stopListener)
				pvcPath, _ := cli.GetPVPath(newPVC.Name)
				_ = aiModel.UpdatePVC(newPVC.Name, pvcPath)
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

	pvcName := fmt.Sprintf("user%d-model%d", aiModel.UserID, aiModel.ID)
	if err := cli.CreatePVC(pvcName, "1Gi", map[string]string{
		"app":  "model",
		"id":   strconv.Itoa(int(aiModel.ID)),
		"user": strconv.Itoa(int(aiModel.UserID)),
	}); err != nil {
		close(stopListener)
		_ = aiModel.UpdatePVCStatus(core.ClaimLost)
	}
}

func createNNModelFile(aiModel *model.AIModel, c chan string) {
	status := <-c
	defer close(c)
	if status == "fail" {
		_ = aiModel.UpdateAIModelStatus(model.AIMODEL_STATUS_UNAVAILABLE)
		return
	}

	if aiModel.Framework == model.AIMODEL_FRAMEWORK_PYTORCH {
		f, err := os.Create(filepath.Join(aiModel.PersistentVolumePath, "nn_model.py"))
		if err != nil {
			_ = aiModel.UpdateAIModelStatus(model.AIMODEL_STATUS_UNAVAILABLE)
			return
		}
		_, err = f.Write([]byte(aiModel.Network))
		if err != nil {
			_ = aiModel.UpdateAIModelStatus(model.AIMODEL_STATUS_UNAVAILABLE)
			return
		}
		_ = f.Close()
	}
	_ = aiModel.UpdateAIModelStatus(model.AIMODEL_STATUS_IDLE)
}
