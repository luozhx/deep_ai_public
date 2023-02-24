package client

import (
	conf "deep-ai-server/app/config"
	pytorch "github.com/kubeflow/pytorch-operator/pkg/client/clientset/versioned"
	tensorflow "github.com/kubeflow/tf-operator/pkg/client/clientset/versioned"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type TrainingOption struct {
	Command        []string
	Args           []string
	Framework      string
	CodePVCName    string
	DatasetPVCName string
	AIModelPVCName string
	WorkerReplicas int32
	Labels         map[string]string
}

type ServingOption struct {
	Args           []string
	AIModelPVCName string
	Framework      string
	Labels         map[string]string
}

type NotebookOption struct {
	CodePVCName string
	Framework   string
	BaseURL     string
	Labels      map[string]string
}

var (
	k8sClient        *k8s.Clientset
	pytorchClient    *pytorch.Clientset
	tensorflowClient *tensorflow.Clientset
)

//var (
//	statusMap = map[string]int{
//		"Created":    model.JOB_STATUS_QUEUING,
//		"Running":    model.JOB_STATUS_TRAINING,
//		"Restarting": model.JOB_STATUS_ERROR,
//		"Succeeded":  model.JOB_STATUS_SUCCESS,
//		"Failed":     model.JOB_STATUS_ERROR,
//	}
//)

const (
	RETRY_TIMES = 30
)

const (
	//NAMESPACE_DEEP_AI = "dev"
	NAMESPACE_DEEP_AI = "deep-ai"
)

func Init(kubeconfigPath string) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		panic(err.Error())
	}
	//TestTFClient = tensorflow.NewForConfigOrDie(config)
	//TestClient = k8s.NewForConfigOrDie(config)

	k8sClient = k8s.NewForConfigOrDie(config)
	pytorchClient = pytorch.NewForConfigOrDie(config)
	tensorflowClient = tensorflow.NewForConfigOrDie(config)
}

type Client struct {
	namespace string
	registry  string
}

func NewClient(namespace string) Client {
	cli := Client{
		namespace: namespace,
		registry:  conf.Docker.Registry,
	}
	return cli
}
