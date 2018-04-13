package kubehandler

import "k8s.io/client-go/tools/cache"

//EventHandler represents Event Handling for a Resource
type EventHandler interface {
	GetName() string
	GetSynced() cache.InformerSynced
	GetInformer() cache.SharedInformer
	AddFunc(namespace, name string) error
	UpdateFunc(namespace, name string) error
	DeleteFunc(namespace, name string) error
}
