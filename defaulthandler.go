package kubehandler

import (
	"context"

	"k8s.io/client-go/tools/cache"
)

type DefaultHandler struct {
	Synced   cache.InformerSynced
	Informer cache.SharedInformer
	Name     string
}

func (handler *DefaultHandler) GetName() string {
	return handler.Name
}
func (handler *DefaultHandler) GetSynced() cache.InformerSynced {
	return handler.Synced
}
func (handler *DefaultHandler) GetInformer() cache.SharedInformer {
	return handler.Informer
}

func (handler *DefaultHandler) AddFunc(ctx context.Context, namespace, name string) error {
	return nil
}

func (handler *DefaultHandler) UpdateFunc(ctx context.Context, namespace, name string) error {
	return nil
}

func (handler *DefaultHandler) DeleteFunc(ctx context.Context, namespace, name string) error {
	return nil
}
