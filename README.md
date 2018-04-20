# Kubehandler

An event dispatcher for Kubernetes controllers.

__Note: This is alpha software. Please use with caution.__

### Sample Controller

There is a sample controller available in `sample-controller/`. You can use
that as a starting point for building a new controller.

### EventHandler

Kubehandler defines a Go interface, EventHandler. Any type that implements this
interface can be used to handle events.

```
type EventHandler interface {
	GetName() string
	GetSynced() cache.InformerSynced
	GetInformer() cache.SharedInformer
	AddFunc(namespace, name string) error
	UpdateFunc(namespace, name string) error
	DeleteFunc(namespace, name string) error
}
```

Each of the Add/Update/Delete Funcs receive the namespace and the name of the
resource that has been modified (or created or deleted).
`kubehandler.DefaultHandler` implements a DefaultHandler that accepts all
events and does nothing. In order to make use of this behaviour, you can use
the bundled DefaultHandler by inheriting from it:

```
type DeploymentsHandler struct {
	kubehandler.DefaultHandler
}
```

You will need an EventLoop to consume events. You can create one like so:

```
	loop := kubehandler.NewEventLoop("workqueueName")
```

You can then register the EventHandler with the EventLoop.

```
	loop.Register(&DeploymentsHandler{
		DefaultHandler: kubehandler.DefaultHandler{
			Synced:   deploymentsInformer.Informer().HasSynced,
			Informer: deploymentsInformer.Informer(),
		},
	})

```

Finally, run the EventLoop.


```
	threadiness := 2
	stopCh := make(chan struct{})

	loop.Run(threadiness, stopCh)
```

You may use the channel passed in to `EventLoop.Run` to stop the EventLoop.

```
	close(stopCh)
```

## DefaultHandler
DefaultHandler provides implementations of EventHandler.

	GetName() string
Returns defaultHandler.Name

	GetSynced() cache.InformerSynced
Returns defaultHandler.Synced

	GetInformer() cache.SharedInformer
Returns defaultHandler.Informer

The following functions are no-ops.
```
AddFunc(namespace, name string) error
UpdateFunc(namespace, name string) error
DeleteFunc(namespace, name string) error
```
