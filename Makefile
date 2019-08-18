HUB=github.com/krubot/terraform-operator

build:
	export GO111MODULE=on
	go build $(HUB)/cmd/terraform-operator
