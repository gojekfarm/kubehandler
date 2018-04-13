package kubehandler

import (
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	WorkqueueAddEvent    string = "add"
	WorkqueueUpdateEvent string = "update"
	WorkqueueDeleteEvent string = "delete"
)

//WorkQueueHandler defines the contract of a handler function
type WorkQueueHandler func(namespace, name string) error

//WorkQueue manages the rate limiting interface
type WorkQueue interface {
	Run(threadiness int, stopCh <-chan struct{}) error
	AddSynced(cache.InformerSynced)
	EnqueueAdd(kind string, obj interface{})
	EnqueueUpdate(kind string, obj interface{})
	EnqueueDelete(kind string, obj interface{})
	RegisterAddHandler(kind string, handler WorkQueueHandler)
	RegisterUpdateHandler(kind string, handler WorkQueueHandler)
	RegisterDeleteHandler(kind string, handler WorkQueueHandler)
	Length() int
}

type workQueue struct {
	name             string
	workqueue        workqueue.RateLimitingInterface
	informerSynced   []cache.InformerSynced
	addHandlerMap    map[string]WorkQueueHandler
	updateHandlerMap map[string]WorkQueueHandler
	deleteHandlerMap map[string]WorkQueueHandler
}

func (q *workQueue) shutDown() {
	q.workqueue.ShutDown()
}

func (q *workQueue) Length() int {
	return q.workqueue.Len()
}

func (q *workQueue) AddSynced(informer cache.InformerSynced) {
	q.informerSynced = append(q.informerSynced, informer)
}

//Run is the WorkQueue entry point
func (q *workQueue) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer q.shutDown()

	if ok := cache.WaitForCacheSync(stopCh, q.informerSynced...); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(q.runWorker, time.Second, stopCh)
	}

	<-stopCh

	return nil
}

func (q *workQueue) runWorker() {
	for q.processNextWorkItem() {
	}
}

func (q *workQueue) handleAdd(kind, namespace, name string) error {
	return q.addHandlerMap[kind](namespace, name)
}

func (q *workQueue) handleUpdate(kind, namespace, name string) error {
	return q.updateHandlerMap[kind](namespace, name)
}

func (q *workQueue) handleDelete(kind, namespace, name string) error {
	return q.deleteHandlerMap[kind](namespace, name)
}

func (q *workQueue) syncHandler(key string) error {
	splitKey := strings.Split(key, ":")
	eventType, kind, nsKey := splitKey[0], splitKey[1], splitKey[2]
	namespace, name, err := cache.SplitMetaNamespaceKey(nsKey)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	switch eventType {
	case WorkqueueAddEvent:
		return q.handleAdd(kind, namespace, name)
	case WorkqueueUpdateEvent:
		return q.handleUpdate(kind, namespace, name)
	case WorkqueueDeleteEvent:
		return q.handleDelete(kind, namespace, name)
	default:
		runtime.HandleError(fmt.Errorf("invalid event type: %s", eventType))
		return nil
	}
}

func (q *workQueue) processNextWorkItem() bool {
	obj, shutdown := q.workqueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer q.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			q.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		if err := q.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		q.workqueue.Forget(obj)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

func (q *workQueue) EnqueueAdd(kind string, obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}

	q.workqueue.AddRateLimited(fmt.Sprintf("%s:%s:%s", WorkqueueAddEvent, kind, key))
}

func (q *workQueue) EnqueueUpdate(kind string, obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}

	q.workqueue.AddRateLimited(fmt.Sprintf("%s:%s:%s", WorkqueueUpdateEvent, kind, key))
}

func (q *workQueue) EnqueueDelete(kind string, obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}

	q.workqueue.AddRateLimited(fmt.Sprintf("%s:%s:%s", WorkqueueDeleteEvent, kind, key))
}

func (q *workQueue) RegisterAddHandler(kind string, handler WorkQueueHandler) {
	q.addHandlerMap[kind] = handler
}

func (q *workQueue) RegisterUpdateHandler(kind string, handler WorkQueueHandler) {
	q.updateHandlerMap[kind] = handler
}

func (q *workQueue) RegisterDeleteHandler(kind string, handler WorkQueueHandler) {
	q.deleteHandlerMap[kind] = handler
}

//NewWorkQueue creates a WorkQueue with a name
func NewWorkQueue(name string) WorkQueue {
	return &workQueue{
		name:             name,
		workqueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), name),
		addHandlerMap:    make(map[string]WorkQueueHandler),
		updateHandlerMap: make(map[string]WorkQueueHandler),
		deleteHandlerMap: make(map[string]WorkQueueHandler),
	}
}
