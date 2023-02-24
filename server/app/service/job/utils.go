package job

import (
	"deep-ai-server/app/model"
	"deep-ai-server/app/tools/client"
	"fmt"
	pytorch "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1"
	common "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
	tensorflow "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1"
	"k8s.io/client-go/tools/cache"
	"strconv"
)

func createJob(job model.Job, code model.Code, dataset model.Dataset, AIModel model.AIModel, args []string, stopListener chan struct{}) {
	cli := client.NewClient(client.NAMESPACE_DEEP_AI)
	jobName := fmt.Sprintf("user%d-job%d", job.UserID, job.ID)
	var informer cache.SharedIndexInformer
	if job.Framework == model.JOB_FRAMEWORK_TENSORFLOW {
		informer = cli.GetTFJobInformer(fmt.Sprintf(
			"app=training,id=%d,user=%d",
			job.ID,
			job.UserID,
		))
	} else if job.Framework == model.JOB_FRAMEWORK_PYTORCH {
		informer = cli.GetPyTorchJobInformer(fmt.Sprintf(
			"app=training,id=%d,user=%d",
			job.ID,
			job.UserID,
		))
	}

	if informer == nil {
		close(stopListener)
		return
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: nil,
		UpdateFunc: func(oldObj, newObj interface{}) {
			var statuses map[common.ReplicaType]*common.ReplicaStatus
			if job.Framework == model.JOB_FRAMEWORK_TENSORFLOW {
				statuses = newObj.(*tensorflow.TFJob).Status.ReplicaStatuses
			} else if job.Framework == model.AIMODEL_FRAMEWORK_PYTORCH {
				statuses = newObj.(*pytorch.PyTorchJob).Status.ReplicaStatuses
			}

			active, succeeded, failed := 0, 0, 0
			for _, v := range statuses {
				active += int(v.Active)
				succeeded += int(v.Succeeded)
				failed += int(v.Failed)
			}

			if failed > 0 {
				close(stopListener)
				_ = job.UpdateJobStatus(model.JOB_STATUS_FAILED)
			} else if active+succeeded < len(statuses) {
				_ = job.UpdateJobStatus(model.JOB_STATUS_CREATING)
			} else if active > 0 && active+succeeded == len(statuses) {
				_ = job.UpdateJobStatus(model.JOB_STATUS_TRAINING)
			} else if succeeded == len(statuses) {
				close(stopListener)
				_ = job.UpdateJobStatus(model.JOB_STATUS_FINISHED)
			}
		},
		DeleteFunc: nil,
	})
	go informer.Run(stopListener)

	options := client.TrainingOption{
		Command: []string{
			"python3",
			job.EntryPoint,
		},
		Args:           args,
		Framework:      job.Framework,
		CodePVCName:    code.PersistentVolumeClaimName,
		DatasetPVCName: dataset.PersistentVolumeClaimName,
		AIModelPVCName: AIModel.PersistentVolumeClaimName,
		WorkerReplicas: int32(job.Num),
		Labels: map[string]string{
			"app":  "training",
			"id":   strconv.Itoa(int(job.ID)),
			"user": strconv.Itoa(int(job.UserID)),
		},
	}

	if err := cli.CreateTrainingJob(jobName, options); err != nil {
		close(stopListener)
		_ = job.UpdateJobStatus(model.JOB_STATUS_FAILED)
		return
	}
}
