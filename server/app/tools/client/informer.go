package client

import (
	pytorch "github.com/kubeflow/pytorch-operator/pkg/client/informers/externalversions/pytorch/v1"
	tensorflow "github.com/kubeflow/tf-operator/pkg/client/informers/externalversions/tensorflow/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

// example: labelSelectorSrc="app=notebook,user=1"
func (cli Client) GetInformerFactory(labelSelectorSrc string) informers.SharedInformerFactory {
	return informers.NewSharedInformerFactoryWithOptions(
		k8sClient,
		0,
		informers.WithNamespace(cli.namespace),
		informers.WithTweakListOptions(func(options *v1.ListOptions) {
			options.LabelSelector = labelSelectorSrc
		}),
	)
}

func (cli Client) GetTFJobInformer(labelSelectorSrc string) cache.SharedIndexInformer {
	return tensorflow.NewFilteredTFJobInformer(
		tensorflowClient,
		cli.namespace,
		0,
		cache.Indexers{},
		func(options *v1.ListOptions) {
			options.LabelSelector = labelSelectorSrc
		},
	)
}

func (cli Client) GetPyTorchJobInformer(labelSelectorSrc string) cache.SharedIndexInformer {
	return pytorch.NewFilteredPyTorchJobInformer(
		pytorchClient,
		cli.namespace,
		0,
		cache.Indexers{},
		func(options *v1.ListOptions) {
			options.LabelSelector = labelSelectorSrc
		},
	)
}
