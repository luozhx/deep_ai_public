package code

import (
	"deep-ai-server/app/db"
	"deep-ai-server/app/model"
	"deep-ai-server/app/tools/client"
	"fmt"
	core "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"strconv"
)

func createPVC(code model.Code, c chan string, stopListener chan struct{}) {
	cli := client.NewClient(client.NAMESPACE_DEEP_AI)
	informer := cli.GetInformerFactory(
		fmt.Sprintf(
			"app=code,id=%d,user=%d",
			code.ID,
			code.UserID,
		)).Core().V1().PersistentVolumeClaims().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: nil,
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPVC := oldObj.(*core.PersistentVolumeClaim)
			newPVC := newObj.(*core.PersistentVolumeClaim)

			if oldPVC.Status.Phase == newPVC.Status.Phase {
				return
			}

			_ = code.UpdatePVCStatus(newPVC.Status.Phase)
			switch newPVC.Status.Phase {
			case core.ClaimPending:
				break
			case core.ClaimBound:
				close(stopListener)
				pvcPath, _ := cli.GetPVPath(newPVC.Name)
				_ = code.UpdatePVC(newPVC.Name, pvcPath)
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

	pvcName := fmt.Sprintf("user%d-code%d", code.UserID, code.ID)
	if err := cli.CreatePVC(pvcName, "1Gi", map[string]string{
		"app":  "code",
		"id":   strconv.Itoa(int(code.ID)),
		"user": strconv.Itoa(int(code.UserID)),
	}); err != nil {
		close(stopListener)
		_ = code.UpdatePVCStatus(core.ClaimLost)
		c <- "fail"
	}
}

func createNotebook(code model.Code, c chan string, stopListener chan struct{}) {
	status := <-c
	defer close(c)
	if status == "fail" {
		return
	}
	cli := client.NewClient(client.NAMESPACE_DEEP_AI)
	pvcName := fmt.Sprintf("user%d-code%d", code.UserID, code.ID)
	serviceName := pvcName
	count := 0
	informer := cli.GetInformerFactory(fmt.Sprintf("app=notebook,id=%d,user=%d", code.ID, code.UserID)).Core().V1().Pods().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: nil,
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod := oldObj.(*core.Pod)
			newPod := newObj.(*core.Pod)
			if oldPod.Status.Phase == newPod.Status.Phase && newPod.Status.Phase == core.PodPending {
				if count < 5 {
					count++
					return
				}
				close(stopListener)
				_ = code.UpdateNotebookStatus(model.CODE_NOTEBOOK_STATUS_UNAVAILABLE)
				_ = cli.DeleteNotebook(serviceName)
				return
			}

			switch newPod.Status.Phase {
			case core.PodPending:
				_ = code.UpdateNotebookStatus(model.CODE_NOTEBOOK_STATUS_PROCESSING)
				break
			case core.PodRunning:
				close(stopListener)
				_ = code.UpdateNotebookStatus(model.CODE_NOTEBOOK_STATUS_AVAILABLE)
				_ = code.UpdateCodeStatus(model.CODE_STATUS_IDLE)
				break
			case core.PodFailed:
			case core.PodUnknown:
			default:
				close(stopListener)
				_ = code.UpdateNotebookStatus(model.CODE_NOTEBOOK_STATUS_UNAVAILABLE)
				_ = cli.DeleteNotebook(serviceName)
				return
			}
		},
		DeleteFunc: nil,
	})
	go informer.Run(stopListener)

	err := cli.CreateNotebook(serviceName, client.NotebookOption{
		CodePVCName: pvcName,
		Framework:   code.Framework,
		BaseURL:     fmt.Sprintf("/api/v1/code/%d/notebook/%s", code.ID, serviceName),
		Labels: map[string]string{
			"app":  "notebook",
			"id":   strconv.Itoa(int(code.ID)),
			"user": strconv.Itoa(int(code.UserID)),
		},
	})
	if err != nil {
		close(stopListener)
		return
	}
	db.DB().Model(&code).Update(model.Code{ServiceName: serviceName})
}
