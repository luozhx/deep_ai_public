package inference

import (
	"deep-ai-server/app/db"
	"deep-ai-server/app/model"
	"deep-ai-server/app/tools/client"
	"fmt"
	core "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"path/filepath"
	"strconv"
)

func createServing(inference model.Inference, AIModel model.AIModel, stopListener chan struct{}) {
	cli := client.NewClient(client.NAMESPACE_DEEP_AI)
	servingName := fmt.Sprintf("user%d-serving%d", inference.UserID, inference.ID)
	count := 0
	informer := cli.GetInformerFactory(fmt.Sprintf("app=serving,id=%d,user=%d", inference.ID, inference.UserID)).Core().V1().Pods().Informer()
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
				_ = inference.UpdateInferenceStatus(model.INFERENCE_STATUS_UNAVAILABLE)
				_ = cli.DeleteServing(servingName)
				return
			}

			switch newPod.Status.Phase {
			case core.PodPending:
				_ = inference.UpdateInferenceStatus(model.INFERENCE_STATUS_PROCESSING)
				break
			case core.PodRunning:
				close(stopListener)
				_ = inference.UpdateInferenceStatus(model.INFERENCE_STATUS_AVAILABLE)
				break
			case core.PodFailed:
			case core.PodUnknown:
			default:
				close(stopListener)
				_ = inference.UpdateInferenceStatus(model.INFERENCE_STATUS_UNAVAILABLE)
				_ = cli.DeleteServing(servingName)
				return
			}
		},
		DeleteFunc: nil,
	})
	go informer.Run(stopListener)

	args := make([]string, 0)
	if inference.Framework == model.INFERENCE_FRAMEWORK_PYTORCH {
		args = []string{fmt.Sprintf("--input-size %s --input-name x --output-name y "+ // TODO: INPUT-NAME OUTPUT-NAME
			"--model-file %s --export-dir /workspace/model/export/114514", inference.Dimension, filepath.Join("/workspace/model/", inference.ModelFile))}
	}
	options := client.ServingOption{
		Framework:      inference.Framework,
		Args:           args,
		AIModelPVCName: AIModel.PersistentVolumeClaimName,
		Labels: map[string]string{
			"app":  "serving",
			"id":   strconv.Itoa(int(inference.ID)),
			"user": strconv.Itoa(int(inference.UserID)),
		},
	}

	if err := cli.CreateServing(servingName, options); err != nil {
		close(stopListener)
		_ = inference.UpdateInferenceStatus(model.INFERENCE_STATUS_UNAVAILABLE)
		return
	}
	db.DB().Model(inference).Update(&model.Inference{ServiceName: servingName})
}
