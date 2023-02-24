package client

import (
	"fmt"
	common "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
	tf "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1"
	operator "github.com/kubeflow/tf-operator/pkg/client/clientset/versioned"
	k8s "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func (cli Client) createTFJob(name string, options TrainingOption) error {
	template := k8s.PodTemplateSpec{
		Spec: k8s.PodSpec{
			Containers: []k8s.Container{
				{
					Command: options.Command,
					Args:    options.Args,
					Image:   fmt.Sprintf("%s/tensorflow/tensorflow:1.15.0-py3", cli.registry),
					Name:    "tensorflow", // must be tensorflow!!!
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

	job := &tf.TFJob{
		TypeMeta: meta.TypeMeta{
			Kind:       tf.Kind,
			APIVersion: tf.GroupVersion,
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: cli.namespace,
			Labels: options.Labels,
		},
		Spec: tf.TFJobSpec{
			TFReplicaSpecs: map[tf.TFReplicaType]*common.ReplicaSpec{
				tf.TFReplicaTypePS: {
					Replicas:      int32ptr(1),
					Template:      template,
					RestartPolicy: common.RestartPolicyNever,
				},
				tf.TFReplicaTypeChief: {
					Replicas:      int32ptr(1),
					Template:      template,
					RestartPolicy: common.RestartPolicyNever,
				},
				tf.TFReplicaTypeWorker: {
					Replicas:      int32ptr(options.WorkerReplicas),
					Template:      template,
					RestartPolicy: common.RestartPolicyNever,
				},
			},
		},
	}

	_, err := tensorflowClient.KubeflowV1().TFJobs(cli.namespace).Create(job)
	return err
}

//func (cli Client) getTFJobStatus(jobName string) int {
//	job, err := tensorflowClient.KubeflowV1().TFJobs(cli.namespace).Get(jobName, meta.GetOptions{})
//	if err != nil {
//		return model.JOB_STATUS_ERROR
//	}
//	if len(job.Status.Conditions) == 0 {
//		return model.JOB_STATUS_QUEUING
//	}
//	return statusMap[string(job.Status.Conditions[len(job.Status.Conditions)-1].Type)]
//}

func (cli Client) deleteTFJob(name string) error {
	return tensorflowClient.KubeflowV1().TFJobs(cli.namespace).Delete(name, &meta.DeleteOptions{})
}

func CreateTFJob() {
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
						"/opt/model.py",
					},
					Args: []string{
						fmt.Sprintf("--tf-model-dir=%s", "/mnt"),
						fmt.Sprintf("--tf-export-dir=%s", "/mnt/export"),
						fmt.Sprintf("--tf-train-steps=%s", "200"),
						fmt.Sprintf("--tf-batch-size=%s", "100"),
						fmt.Sprintf("--tf-learning-rate=%s", "0.01"),
					},
					Image: "k8s-1:5000/kubeflow/mytfmodel:v1",
					Name:  "tensorflow", // must be tensorflow!!!
					VolumeMounts: []k8s.VolumeMount{
						{
							Name:      "local-storage",
							MountPath: "/mnt",
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
							ClaimName: "tf-mnist",
							ReadOnly:  false,
						},
					},
				},
			},
			RestartPolicy: k8s.RestartPolicyOnFailure,
		},
	}

	job := &tf.TFJob{
		TypeMeta: meta.TypeMeta{
			Kind:       tf.Kind,
			APIVersion: tf.GroupVersion,
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "tf-mnist",
			Namespace: "kubeflow",
		},
		Spec: tf.TFJobSpec{
			TFReplicaSpecs: map[tf.TFReplicaType]*common.ReplicaSpec{
				tf.TFReplicaTypeChief: {
					Replicas: int32ptr(1),
					Template: template,
				},
				tf.TFReplicaTypePS: {
					Replicas: int32ptr(1),
					Template: template,
				},
				tf.TFReplicaTypeWorker: {
					Replicas: int32ptr(2),
					Template: template,
				},
			},
		},
	}

	result, err := client.KubeflowV1().TFJobs("kubeflow").Create(job)
	fmt.Printf("%+v\n", result)
	fmt.Println(err)

	//err = client.KubeflowV1().TFJobs("kubeflow").Delete(job.Name, &meta.DeleteOptions{})
	//fmt.Println(err)
}
