HUB=quay.io/krubot/terraform-operator

all: generate manifests image

generate:
	operator-sdk generate k8s

manifests: controller-gen
	$(CONTROLLER_GEN) crd:trivialVersions=true paths=./pkg/apis/... output:artifacts:config=deploy/00-crds

image:
	operator-sdk build $(HUB)
	#CGO_ENABLED=0 GOOS=linux go build \
	#  -o "${PWD}/terraform-operator/build/_output/bin/terraform-operator" \
	#	${PWD}/cmd/manager
	#docker build -f build/Dockerfile -t "$(HUB)" .

packages:
	go get -u ./...
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
