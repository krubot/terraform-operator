FROM golang:alpine AS terraform

ENV TERRAFORM_VERSION=0.10.0

RUN apk add --update git bash openssh

ENV TF_DEV=true
ENV TF_RELEASE=true

WORKDIR $GOPATH/src/github.com/hashicorp/terraform
RUN git clone https://github.com/hashicorp/terraform.git ./ && \
    git checkout v${TERRAFORM_VERSION} && \
    /bin/bash scripts/build.sh

RUN mkdir -p /out && \
    cp $GOPATH/bin/terraform /out/terraform

FROM alpine:3.6

RUN adduser -D terraform-operator
USER terraform-operator
WORKDIR /etc/infra

COPY --from=terraform /out/terraform /usr/bin/terraform
ADD tmp/build/run-terraform-validate.sh /usr/local/bin/run-terraform-validate.sh
ADD tmp/_output/bin/terraform-operator /usr/local/bin/terraform-operator
