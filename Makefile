HUB=quay.io/krubot/terraform-operator

.PHONY: image

image:
	operator-sdk build $(HUB)
