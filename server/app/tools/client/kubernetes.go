package client

import (
	"deep-ai-server/app/model"
	"fmt"
	apps "k8s.io/api/apps/v1"
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func int32ptr(i int32) *int32 { return &i }

func (cli Client) CreatePVC(name, storage string, labels map[string]string) error {
	pvc := api.PersistentVolumeClaim{
		ObjectMeta: meta.ObjectMeta{
			Name: name,
			Labels: labels,
		},
		Spec: api.PersistentVolumeClaimSpec{
			AccessModes: []api.PersistentVolumeAccessMode{
				api.ReadWriteMany,
			},
			Resources: api.ResourceRequirements{
				Requests: api.ResourceList{
					api.ResourceStorage: resource.MustParse(storage),
				},
			},
		},
	}

	_, err := k8sClient.CoreV1().PersistentVolumeClaims(cli.namespace).Create(&pvc)
	return err
}

func (cli Client) DeletePVC(name string) error {
	return k8sClient.CoreV1().PersistentVolumeClaims(cli.namespace).Delete(name, &meta.DeleteOptions{})
}

func (cli Client) GetPVPath(PVCName string) (string, error) {
	pvc, err := k8sClient.CoreV1().PersistentVolumeClaims(cli.namespace).Get(PVCName, meta.GetOptions{})
	if err != nil {
		return "", err
	}
	pv, err := k8sClient.CoreV1().PersistentVolumes().Get(pvc.Spec.VolumeName, meta.GetOptions{})
	if err != nil {
		return "", err
	}
	return pv.Spec.NFS.Path, nil
}

func (cli Client) CreateServing(name string, options ServingOption) error {
	image := fmt.Sprintf("%s/tensorflow/serving:2.0.0", cli.registry)
	command := []string{
		"/usr/bin/tensorflow_model_server",
	}
	args := append([]string{ //todo: process args
		"--rest_api_port=9000",
		"--port=8500",
		"--model_base_path=/workspace/model/export",
		"--model_name=model",
	}, options.Args...)
	if options.Framework == model.INFERENCE_FRAMEWORK_PYTORCH {
		image = fmt.Sprintf("%s/pytorch/serving:2.0.0-v1.0", cli.registry)
		command = []string{
			"/entrypoint.sh",
		}
		args = []string{
			options.Args[0],
			"--rest_api_port=9000 --port=8500 --model_base_path=/workspace/model/export --model_name=model",
		}
	}

	service := api.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:        name,
			Namespace:   cli.namespace,
			Labels:      options.Labels,
			Annotations: nil, //todo
		},
		Spec: api.ServiceSpec{
			Ports: []api.ServicePort{
				//{
				//	Name:       "grpc-tf-serving",
				//	Port:       8500,
				//	TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8500},
				//},
				{
					Name:       "http-tf-serving",
					Port:       80,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9000},
					Protocol:   api.ProtocolTCP,
				},
			},
			Selector: options.Labels,
			Type:     api.ServiceTypeClusterIP,
		},
	}

	deployment := apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: cli.namespace,
			Labels:    options.Labels,
		},
		Spec: apps.DeploymentSpec{
			Selector: &meta.LabelSelector{
				MatchLabels: options.Labels,
			},
			Template: api.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: options.Labels,
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						{
							Name:            "tf-serving",
							Image:           image,
							Command:         command,
							Args:            args,
							ImagePullPolicy: api.PullIfNotPresent,
							Ports: []api.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 9000,
								},
								//{
								//	Name:          "rpc",
								//	ContainerPort: 8500,
								//},
							},
							Resources: api.ResourceRequirements{ //todo
								Limits: api.ResourceList{
									api.ResourceCPU:    resource.MustParse("4"),
									api.ResourceMemory: resource.MustParse("4Gi"),
								},
								Requests: api.ResourceList{
									api.ResourceCPU:    resource.MustParse("1"),
									api.ResourceMemory: resource.MustParse("1Gi"),
								},
							},
							VolumeMounts: []api.VolumeMount{
								{
									MountPath: "/workspace/model/",
									Name:      "model-storage",
								},
							},
						},
					},
					Volumes: []api.Volume{
						{
							Name: "model-storage",
							VolumeSource: api.VolumeSource{
								PersistentVolumeClaim: &api.PersistentVolumeClaimVolumeSource{
									ClaimName: options.AIModelPVCName,
									ReadOnly:  false,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := k8sClient.CoreV1().Services(cli.namespace).Create(&service)
	if err != nil {
		return err
	}
	_, err = k8sClient.AppsV1().Deployments(cli.namespace).Create(&deployment)
	if err != nil {
		_ = k8sClient.CoreV1().Services(cli.namespace).Delete(service.Name, &meta.DeleteOptions{})
		return err
	}
	return err
}

func (cli Client) DeleteServing(name string) error {
	if err := k8sClient.CoreV1().Services(cli.namespace).Delete(name, &meta.DeleteOptions{}); err != nil {
		return err
	}
	return k8sClient.AppsV1().Deployments(cli.namespace).Delete(name, &meta.DeleteOptions{})
}

func (cli Client) CreateNotebook(name string, options NotebookOption) error {
	image := "jupyter/tensorflow-1.15.0-notebook:1.0.0"
	if options.Framework == model.CODE_FRAMEWORK_PYTORCH {
		image = "jupyter/pytorch-1.3.1-notebook:1.0.0"
	}

	service := api.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:        name,
			Namespace:   cli.namespace,
			Labels:      options.Labels,
			Annotations: nil, //todo
		},
		Spec: api.ServiceSpec{
			Ports: []api.ServicePort{
				{
					Name:       "jupyter-notebook",
					Port:       80,
					TargetPort: intstr.FromInt(8888),
					Protocol:   api.ProtocolTCP,
				},
			},
			Selector: options.Labels,
			Type:     api.ServiceTypeClusterIP,
		},
	}

	deployment := apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: cli.namespace,
			Labels:    options.Labels,
		},
		Spec: apps.DeploymentSpec{
			Selector: &meta.LabelSelector{
				MatchLabels: options.Labels,
			},
			Template: api.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: options.Labels,
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						{
							Name:            "jupyter-notebook",
							Image:           fmt.Sprintf("%s/%s", cli.registry, image),
							ImagePullPolicy: api.PullIfNotPresent,
							Command: []string{
								"jupyter",
								"notebook",
								"--NotebookApp.token=''",
								fmt.Sprintf("--NotebookApp.base_url='%s'", options.BaseURL),
								"--NotebookApp.notebook_dir='/home/jovyan/work'",
								"--NotebookApp.allow_origin='*'",
								"--NotebookApp.allow_credentials=True",
							},
							Ports: []api.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8888,
								},
							},
							VolumeMounts: []api.VolumeMount{
								{
									MountPath: "/home/jovyan/work",
									Name:      "code-storage",
								},
							},
						},
					},
					Volumes: []api.Volume{
						{
							Name: "code-storage",
							VolumeSource: api.VolumeSource{
								PersistentVolumeClaim: &api.PersistentVolumeClaimVolumeSource{
									ClaimName: options.CodePVCName,
									ReadOnly:  false,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := k8sClient.CoreV1().Services(cli.namespace).Create(&service)
	if err != nil {
		return err
	}
	_, err = k8sClient.AppsV1().Deployments(cli.namespace).Create(&deployment)
	if err != nil {
		_ = k8sClient.CoreV1().Services(cli.namespace).Delete(service.Name, &meta.DeleteOptions{})
		return err
	}

	return nil
}

func (cli Client) DeleteNotebook(name string) error {
	if err := k8sClient.CoreV1().Services(cli.namespace).Delete(name, &meta.DeleteOptions{}); err != nil {
		return err
	}
	return k8sClient.AppsV1().Deployments(cli.namespace).Delete(name, &meta.DeleteOptions{})
}

func (cli Client) CreateTrainingJob(name string, option TrainingOption) error {
	if option.Framework == model.JOB_FRAMEWORK_TENSORFLOW {
		return cli.createTFJob(name, option)
	} else if option.Framework == model.JOB_FRAMEWORK_PYTORCH {
		return cli.createPytorchJob(name, option)
	}
	return nil
}

func (cli Client) DeleteTrainingJob(name, framework string) error {
	if framework == model.JOB_FRAMEWORK_TENSORFLOW {
		return cli.deleteTFJob(name)
	} else if framework == model.JOB_FRAMEWORK_PYTORCH {
		return cli.deletePytorchJob(name)
	}
	return nil
}

//func (cli Client) GetJobStatus(name, framework string) int {
//	if framework == model.JOB_FRAMEWORK_TENSORFLOW {
//		return cli.getTFJobStatus(name)
//	} else if framework == model.JOB_FRAMEWORK_PYTORCH {
//		return cli.getPytorchJobStatus(name)
//	}
//	return model.JOB_STATUS_ERROR
//}

func CreateTFServing() {
	config, err := clientcmd.BuildConfigFromFlags("", "conf/kubeconfig")
	if err != nil {
		panic(err.Error())
	}
	clientset := kubernetes.NewForConfigOrDie(config)

	service := api.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      "mnist-service-local-service",
			Namespace: "kubeflow",
			Labels: map[string]string{
				"app": "mnist",
			},
			Annotations: nil,
		},
		Spec: api.ServiceSpec{
			Ports: []api.ServicePort{
				{
					Name:       "grpc-tf-serving",
					Port:       9000,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9000},
				},
				{
					Name:       "http-tf-serving",
					Port:       8500,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8500},
				},
			},
			Selector: map[string]string{
				"app": "mnist",
			},
			Type: api.ServiceTypeClusterIP,
		},
	}

	deployment := apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name:      "mnist-service-local-deployment",
			Namespace: "kubeflow",
			Labels: map[string]string{
				"app": "mnist",
			},
		},
		Spec: apps.DeploymentSpec{
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{
					"app":     "mnist",
					"version": "v1",
				},
			},
			Template: api.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{
						"app":     "mnist",
						"version": "v1",
					},
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						{
							Name:  "mnist",
							Image: "tensorflow/serving:1.11.1",
							Command: []string{
								"/usr/bin/tensorflow_model_server",
							},
							Args: []string{
								"--port=8500",
								"--rest_api_port=8501",
								"--model_name=mnist",
								"--model_base_path=/mnt/export",
							},
							Env: []api.EnvVar{
								{
									Name:  "modelBasePath",
									Value: "/mnt/export",
								},
							},
							ImagePullPolicy: api.PullIfNotPresent,
							Ports: []api.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8500,
								},
								{
									Name:          "rpc",
									ContainerPort: 9000,
								},
							},
							Resources: api.ResourceRequirements{
								Limits: api.ResourceList{
									api.ResourceCPU:    resource.MustParse("4"),
									api.ResourceMemory: resource.MustParse("4Gi"),
								},
								Requests: api.ResourceList{
									api.ResourceCPU:    resource.MustParse("1"),
									api.ResourceMemory: resource.MustParse("1Gi"),
								},
							},
							VolumeMounts: []api.VolumeMount{
								{
									MountPath: "/mnt",
									Name:      "local-storage",
								},
							},
						},
					},
					Volumes: []api.Volume{
						{
							Name: "config-volume",
							VolumeSource: api.VolumeSource{
								ConfigMap: &api.ConfigMapVolumeSource{
									LocalObjectReference: api.LocalObjectReference{
										Name: "mnist-deploy-config",
									},
								},
							},
						},
						{
							Name: "local-storage",
							VolumeSource: api.VolumeSource{
								PersistentVolumeClaim: &api.PersistentVolumeClaimVolumeSource{
									ClaimName: "mnist-pvc",
									ReadOnly:  false,
								},
							},
						},
					},
				},
			},
		},
	}

	//result1, err := clientset.CoreV1().ConfigMaps("kubeflow").Create(&configMap)
	//fmt.Printf("%+v\n", result1)
	//fmt.Println(err)
	//
	//result2, err := clientset.CoreV1().Services("kubeflow").Create(&service)
	//fmt.Printf("%+v\n", result2)
	//fmt.Println(err)
	//
	//result3, err := clientset.AppsV1().Deployments("kubeflow").Create(&deployment)
	//fmt.Printf("%+v\n", result3)
	//fmt.Println(err)

	err = clientset.CoreV1().Services("kubeflow").Delete(service.Name, &meta.DeleteOptions{})
	fmt.Println(err)

	err = clientset.AppsV1().Deployments("kubeflow").Delete(deployment.Name, &meta.DeleteOptions{})
	fmt.Println(err)

}

func CreatePVC(name, storage string) {
	config, err := clientcmd.BuildConfigFromFlags("", "conf/kubeconfig")
	if err != nil {
		panic(err.Error())
	}
	clientset := kubernetes.NewForConfigOrDie(config)

	pvc := api.PersistentVolumeClaim{
		ObjectMeta: meta.ObjectMeta{
			Name: name,
		},
		Spec: api.PersistentVolumeClaimSpec{
			AccessModes: []api.PersistentVolumeAccessMode{
				api.ReadWriteMany,
			},
			Resources: api.ResourceRequirements{
				Requests: api.ResourceList{
					api.ResourceStorage: resource.MustParse(storage),
				},
			},
		},
	}

	result, err := clientset.CoreV1().PersistentVolumeClaims("kubeflow").Create(&pvc)
	fmt.Printf("%+v\n", result)
	fmt.Println(err)
}
