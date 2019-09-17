HUB=quay.io/krubot/terraform-operator

all: generate image

generate:
	operator-sdk generate k8s
	operator-sdk generate openapi

image:
	operator-sdk build $(HUB)

packages:
	go get -u
	go mod tidy
	go mod vendor

fmt:
	go fmt ./...

.PHONY: all image generate
