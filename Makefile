HUB=quay.io/krubot/terraform-operator

.PHONY: all image generate

all: generate manifests image

generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths=./pkg/apis/...

manifests: controller-gen
	$(CONTROLLER_GEN) crd paths=./pkg/apis/... output:artifacts:config=deploy/00-crds

image:
	CGO_ENABLED=0 GOOS=linux go build \
	  -o "${PWD}/build/_output/bin/terraform-operator" \
		${PWD}/cmd/manager
	docker build -f build/Dockerfile -t "$(HUB)" .

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
