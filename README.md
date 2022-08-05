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
	AddFunc(ctx context.Context, namespace, name string) error
	UpdateFunc(ctx context.Context, namespace, name string) error
	DeleteFunc(ctx context.Context, namespace, name string) error
}
```

Each of the Add/Update/Delete Funcs receive a context and the namespace and the name of the
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
	ctx := context.Background()

	loop.Run(ctx, threadiness)
```

You may use the context passed in to `EventLoop.Run` to stop the EventLoop.

```
	ctx, cancelFunc := context.WithCancel(context.Background())
	loop.Run(ctx, threadiness)

	cancelFunc()
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
AddFunc(ctx context.Context, namespace, name string) error
UpdateFunc(ctx context.Context, namespace, name string) error
DeleteFunc(ctx context.Context, namespace, name string) error
```

Note: The current module supports Kubernetes 1.22 and above.
Kubernetes 1.21 and lower compatible module can be found in branch v0.