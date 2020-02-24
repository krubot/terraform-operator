HUB=quay.io/krubot/terraform-operator

.PHONY: all image generate

all: generate manifests image

generate:
	$(CONTROLLER_GEN) object paths=./pkg/apis/terraform/v1alpha1 output:dir=./pkg/apis/terraform/v1alpha1

manifests: controller-gen
	$(CONTROLLER_GEN) crd paths=./pkg/apis/... output:artifacts:config=deploy/00-crds

image:
	CGO_ENABLED=0 GOOS=linux go build \
	  -o "${PWD}/build/_output/bin/terraform-operator" \
		${PWD}/cmd/manager
	docker build -f build/Dockerfile -t "$(HUB)" .

fmt:
	go fmt ./...

controller-gen:
  ifeq (, $(shell which controller-gen))
  go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5
  CONTROLLER_GEN=$(shell go env GOPATH)/bin/controller-gen
  else
  CONTROLLER_GEN=$(shell which controller-gen)
  endif
