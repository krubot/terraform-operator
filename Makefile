HUB=quay.io/krubot/terraform-operator

all: generate manifests image

generate:
	operator-sdk generate k8s

manifests: controller-gen
	$(CONTROLLER_GEN) crd:trivialVersions=true paths="./..." output:crd:artifacts:config=deploy/crds

image:
	operator-sdk build $(HUB)

packages:
	go get -u
	go mod tidy
	go mod vendor

fmt:
	go fmt ./...

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.2
CONTROLLER_GEN=$(shell go env GOPATH)/bin/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

.PHONY: all image generate
