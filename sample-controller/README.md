## Sample Controller

This is a sample controller that uses kubehandler. It can be used as a starting
point for your very own kubernetes controller.

### Setup

Run `make` to build the sample controller.

You will need a service account and a clusterrolebinding which you can create
by running:
```
$ kubectl create serviceaccount logging-controller
$ kubectl create clusterrolebinding logging-controller-admin-binding --clusterrole=view --serviceaccount=default:logging-controller
```

To run it locally and set it up with a local minikube cluster, you can run:
```
make minikube-dev
```

You can then view the logs of the logging-controller pod to view pods.

```
$ kubectl logs -f <logging controller pod name here>
```
