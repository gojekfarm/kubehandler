.PHONY: all
all: build fmt
APP=logging-controller

linux-build:
	env GOOS=linux go build -o $(APP)-linux

compile:
	go build -o $(APP)

build: compile fmt

fmt:
	go fmt $(ALL_PACKAGES)

docker-build: linux-build
	docker build --no-cache -t $(APP) . ;\

docker-dev: docker-build
	kubectl config use-context docker-for-desktop
	kubectl apply -f $(APP)-deployment.yaml

minikube-dev: linux-build
	@eval $$(SHELL=bash minikube docker-env) ;\
	docker build --no-cache -t $(APP) . ;\
	kubectl config use-context minikube
	kubectl apply -f $(APP)-deployment.yaml
