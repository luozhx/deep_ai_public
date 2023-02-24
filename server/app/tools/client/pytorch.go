package client

import (
	"fmt"
	pytorch "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1"
	operator "github.com/kubeflow/pytorch-operator/pkg/client/clientset/versioned"
	common "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
	k8s "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func (cli Client) createPytorchJob(name string, options TrainingOption) error {
	template := k8s.PodTemplateSpec{
		Spec: k8s.PodSpec{
			Containers: []k8s.Container{
				{
					Command: options.Command,
					Args:    options.Args,
					Image:   fmt.Sprintf("%s/pytorch/pytorch:1.3-cuda10.1-cudnn7-runtime", cli.registry),
					Name:    "pytorch",
					VolumeMounts: []k8s.VolumeMount{
						{
							Name:      "code-storage",
							MountPath: "/workspace/code/",
						},
						{
							Name:      "dataset-storage",
							MountPath: "/workspace/dataset/",
						},
						{
							Name:      "model-storage",
							MountPath: "/workspace/model/",
						},
					},
					WorkingDir: "/workspace/code/",
				},
			},
			Volumes: []k8s.Volume{
				{
					Name: "code-storage",
					VolumeSource: k8s.VolumeSource{
						PersistentVolumeClaim: &k8s.PersistentVolumeClaimVolumeSource{
							ClaimName: options.CodePVCName,
							ReadOnly:  false,
						},
					},
				},
				{
					Name: "dataset-storage",
					VolumeSource: k8s.VolumeSource{
						PersistentVolumeClaim: &k8s.PersistentVolumeClaimVolumeSource{
							ClaimName: options.DatasetPVCName,
							ReadOnly:  true,
						},
					},
				},
				{
					Name: "model-storage",
					VolumeSource: k8s.VolumeSource{
						PersistentVolumeClaim: &k8s.PersistentVolumeClaimVolumeSource{
							ClaimName: options.AIModelPVCName,
							ReadOnly:  false,
						},
					},
				},
			},
		},
	}

	job := &pytorch.PyTorchJob{
		TypeMeta: meta.TypeMeta{
			Kind:       pytorch.Kind,
			APIVersion: pytorch.GroupVersion,
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: cli.namespace,
			Labels: options.Labels,
		},
		Spec: pytorch.PyTorchJobSpec{
			PyTorchReplicaSpecs: map[pytorch.PyTorchReplicaType]*common.ReplicaSpec{
				pytorch.PyTorchReplicaTypeMaster: {
					Replicas:      int32ptr(1),
					RestartPolicy: common.RestartPolicyNever,
					Template:      template,
				},
				pytorch.PyTorchReplicaTypeWorker: {
					Replicas:      int32ptr(options.WorkerReplicas),
					RestartPolicy: common.RestartPolicyNever,
					Template:      template,
				},
			},
		},
	}
	_, err := pytorchClient.KubeflowV1().PyTorchJobs(cli.namespace).Create(job)
	return err
}

//func (cli Client) getPytorchJobStatus(jobName string) int {
//	job, err := pytorchClient.KubeflowV1().PyTorchJobs(cli.namespace).Get(jobName, meta.GetOptions{})
//	if err!= nil {
//		return model.JOB_STATUS_ERROR
//	}
//	if len(job.Status.Conditions) == 0 {
//		return model.JOB_STATUS_QUEUING
//	}
//	return statusMap[string(job.Status.Conditions[len(job.Status.Conditions)-1].Type)]
//}

func (cli Client) deletePytorchJob(name string) error {
	return pytorchClient.KubeflowV1().PyTorchJobs(cli.namespace).Delete(name, &meta.DeleteOptions{})
}

func CreatePytorchJob() {
	config, err := clientcmd.BuildConfigFromFlags("", "conf/kubeconfig")
	if err != nil {
		panic(err.Error())
	}
	client := operator.NewForConfigOrDie(config)

	template := k8s.PodTemplateSpec{
		Spec: k8s.PodSpec{
			Containers: []k8s.Container{
				{
					Command: []string{
						"python",
						"-u",
						"/opt/pytorch_dist_mnist/mnist_DDP.py",
					},
					Args: []string{
						fmt.Sprintf("--modelpath=%s", "/mnt/kubeflow-gcfs/pytorch/model"),
					},
					Image: "k8s-1:5000/kubeflow/pytorch-mnist-cpu:v1",
					Name:  "pytorch",
					VolumeMounts: []k8s.VolumeMount{
						{
							MountPath: "/mnt",
							Name:      "local-storage",
						},
					},
					WorkingDir: "/opt",
				},
			},
			Volumes: []k8s.Volume{
				{
					Name: "local-storage",
					VolumeSource: k8s.VolumeSource{
						PersistentVolumeClaim: &k8s.PersistentVolumeClaimVolumeSource{
							ClaimName: "pytorch-mnist",
							ReadOnly:  false,
						},
					},
				},
			},
			RestartPolicy: k8s.RestartPolicyOnFailure,
		},
	}

	job := &pytorch.PyTorchJob{
		TypeMeta: meta.TypeMeta{
			Kind:       pytorch.Kind,
			APIVersion: pytorch.GroupVersion,
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "pytorch-mnist",
			Namespace: "kubeflow",
		},
		Spec: pytorch.PyTorchJobSpec{
			PyTorchReplicaSpecs: map[pytorch.PyTorchReplicaType]*common.ReplicaSpec{
				pytorch.PyTorchReplicaTypeMaster: {
					Replicas:      int32ptr(1),
					RestartPolicy: common.RestartPolicyOnFailure,
					Template:      template,
				},
				pytorch.PyTorchReplicaTypeWorker: {
					Replicas:      int32ptr(2),
					RestartPolicy: common.RestartPolicyOnFailure,
					Template:      template,
				},
			},
		},
	}
	result, err := client.KubeflowV1().PyTorchJobs("kubeflow").Create(job)
	fmt.Printf("%+v\n", result)
	fmt.Println(err)

	//err := client.KubeflowV1().PyTorchJobs("kubeflow").Delete(job.Name, &meta.DeleteOptions{})
	//fmt.Println(err)
}
