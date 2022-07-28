package kubehandler

import (
	"context"

	"k8s.io/client-go/tools/cache"
)

//EventHandler represents Event Handling for a Resource
type EventHandler interface {
	GetName() string
	GetSynced() cache.InformerSynced
	GetInformer() cache.SharedInformer
	AddFunc(ctx context.Context, namespace, name string) error
	UpdateFunc(ctx context.Context, namespace, name string) error
	DeleteFunc(ctx context.Context, namespace, name string) error
}
