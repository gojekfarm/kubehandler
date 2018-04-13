package kubehandler

import "k8s.io/client-go/tools/cache"

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

func (handler *DefaultHandler) AddFunc(namespace, name string) error {
	return nil
}

func (handler *DefaultHandler) UpdateFunc(namespace, name string) error {
	return nil
}

func (handler *DefaultHandler) DeleteFunc(namespace, name string) error {
	return nil
}
