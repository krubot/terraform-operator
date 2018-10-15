FROM golang:alpine AS terraform

ENV TERRAFORM_VERSION=0.11.7

RUN apk add --update git bash openssh

ENV TF_DEV=true
ENV TF_RELEASE=true

WORKDIR $GOPATH/src/github.com/hashicorp/terraform
RUN git clone https://github.com/hashicorp/terraform.git ./ && \
    git checkout v${TERRAFORM_VERSION} && \
    /bin/bash scripts/build.sh

RUN mkdir -p /out && \
    cp $GOPATH/bin/terraform /out/terraform

# alpine base image
FROM golang:alpine

RUN apk add --update py-pip && \
    pip install --upgrade awscli && \
    apk --no-cache add ca-certificates && \
    rm -rf /var/cache/*

# create user with data work directory
RUN addgroup operator && \
    adduser -D -h /tmp terraform && \
    chown -R terraform:operator /tmp 

# change to user
USER terraform
WORKDIR /tmp

# copy files from terraform build and directory
COPY --from=terraform /out/terraform /usr/bin/terraform
ADD tmp/build/run-terraform-validate.sh /tmp/scripts/run-terraform-validate.sh
ADD tmp/_output/bin/terraform-operator /usr/local/bin/terraform-operator
