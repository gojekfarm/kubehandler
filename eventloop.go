package kubehandler

import "k8s.io/client-go/tools/cache"

//EventLoop represents a central EventHandler registry which runs in a loop
type EventLoop interface {
	Run(threadiness int, stopCh <-chan struct{}) error
	Register(handler EventHandler)
}

type eventLoop struct {
	workqueue WorkQueue
}

func (loop *eventLoop) Run(threadiness int, stopCh <-chan struct{}) error {
	return loop.workqueue.Run(threadiness, stopCh)
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
