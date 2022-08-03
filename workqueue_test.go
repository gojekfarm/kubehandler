package kubehandler_test

import (
	"context"
	"testing"
	"time"

	kubehandlerv2 "github.com/gojektech/kubehandler/v2"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
)

func TestShouldEnqueueIntoTheUnderlyingWorkQueue(t *testing.T) {
	workQueue := kubehandlerv2.NewWorkQueue("WorkqueueTest")
	workQueue.EnqueueAdd("someKind", &appsv1.Deployment{})
	timeCompleted := make(chan string, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeCompleted <- "done"
	}()
	select {
	case <-timeCompleted:
		assert.Equal(t, 1, workQueue.Length())
	case <-time.After(2 * time.Second):
		assert.Fail(t, "Nothing in the work queue after timeout")
	}
}

func TestShouldCallRegisteredAddFuncWhenAddEventIsReceived(t *testing.T) {
	workQueue := kubehandlerv2.NewWorkQueue("WorkqueueTest2")
	kind := "Foo"
	addHandlerCalled := make(chan bool, 1)
	stopChan := make(chan struct{}, 1)
	workQueue.RegisterAddHandler(kind, func(ctx context.Context, namespace, name string) error {
		addHandlerCalled <- true
		return nil
	})
	workQueue.EnqueueAdd(kind, &appsv1.Deployment{})
	go workQueue.Run(context.TODO(), 1)
	assert.True(t, <-addHandlerCalled)
	close(stopChan)
}

func TestShouldCallRegisteredUpdateFuncWhenUpdateEventIsReceived(t *testing.T) {
	workQueue := kubehandlerv2.NewWorkQueue("WorkqueueTest3")
	kind := "Foo"
	updateHandlerCalled := make(chan bool, 1)
	stopChan := make(chan struct{}, 1)
	workQueue.RegisterUpdateHandler(kind, func(ctx context.Context, namespace, name string) error {
		updateHandlerCalled <- true
		return nil
	})
	workQueue.EnqueueUpdate(kind, &appsv1.Deployment{})
	go workQueue.Run(context.TODO(), 1)
	assert.True(t, <-updateHandlerCalled)
	close(stopChan)
}

func TestShouldCallRegisteredDeleteFuncWhenDeleteEventIsReceived(t *testing.T) {
	workQueue := kubehandlerv2.NewWorkQueue("WorkqueueTest4")
	kind := "Foo"
	deleteHandlerCalled := make(chan bool, 1)
	stopChan := make(chan struct{}, 1)
	workQueue.RegisterDeleteHandler(kind, func(ctx context.Context, namespace, name string) error {
		deleteHandlerCalled <- true
		return nil
	})
	workQueue.EnqueueDelete(kind, &appsv1.Deployment{})
	go workQueue.Run(context.TODO(), 1)
	assert.True(t, <-deleteHandlerCalled)
	close(stopChan)
}
