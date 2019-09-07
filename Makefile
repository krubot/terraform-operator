HUB=quay.io/krubot/terraform-operator

all: generate image

generate:
	operator-sdk generate k8s
	operator-sdk generate openapi

image:
	operator-sdk build $(HUB)

.PHONY: all image generate
