package kubehandler

import (
	"context"
	"reflect"

	"k8s.io/client-go/tools/cache"
)

//EventLoop represents a central EventHandler registry which runs in a loop
type EventLoop interface {
	Run(ctx context.Context, threadiness int) error
	Register(handler EventHandler)
}

type eventLoop struct {
	workqueue WorkQueue
}

func (loop *eventLoop) Run(ctx context.Context, threadiness int) error {
	return loop.workqueue.Run(ctx, threadiness)
}

func (loop *eventLoop) Register(handler EventHandler) {
	loop.workqueue.AddSynced(handler.GetSynced())

	loop.workqueue.RegisterAddHandler(handler.GetName(), handler.AddFunc)
	loop.workqueue.RegisterUpdateHandler(handler.GetName(), handler.UpdateFunc)
	loop.workqueue.RegisterDeleteHandler(handler.GetName(), handler.DeleteFunc)

	handler.GetInformer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(event interface{}) {
			loop.workqueue.EnqueueAdd(handler.GetName(), event)
		},
		UpdateFunc: func(oldEvent, newEvent interface{}) {

			oldVersion, oldOk := resourceVersion(oldEvent)
			newVersion, newOk := resourceVersion(newEvent)

			if oldOk && newOk && oldVersion == newVersion {
				// Periodic resync will send update events for all known Object.
				// Two different versions of the same Object will always have different RVs.
				return
			}
			loop.workqueue.EnqueueUpdate(handler.GetName(), newEvent)
		},
		DeleteFunc: func(deletedEvent interface{}) {
			loop.workqueue.EnqueueDelete(handler.GetName(), deletedEvent)
		},
	})
}

//NewEventLoop instantiates a workqueue backed EventLoop
func NewEventLoop(name string) EventLoop {
	return &eventLoop{workqueue: NewWorkQueue(name)}
}

func resourceVersion(event interface{}) (string, bool) {

	result, ok := getStringValueByFieldName(event, "ObjectMeta")
	if !ok {
		return "", ok
	}
	result, ok = getStringValueByFieldName(result, "ResourceVersion")
	if !ok {
		return "", ok
	}
	return result.(string), true
}

func getStringValueByFieldName(n interface{}, fieldName string) (interface{}, bool) {
	s := reflect.ValueOf(n)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	if s.Kind() != reflect.Struct {
		return "", false
	}
	f := s.FieldByName(fieldName)
	if !f.IsValid() {
		return "", false
	}

	return f.Interface(), true
}
