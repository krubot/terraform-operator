# Build the manager binary
FROM golang:1.13 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/manager/ cmd/manager/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager cmd/manager/main.go

# Use a minimal redhat container so that we are also able to debug on the container
FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV USER_WORKDIR=/srv \
    USER_UID=0 \
    GROUP_UID=0

RUN mkdir -p ${USER_WORKDIR} \
    && chown ${USER_UID}:${GROUP_UID} ${USER_WORKDIR} \
    && chmod 777 ${USER_WORKDIR}

WORKDIR ${USER_WORKDIR}

COPY --from=builder /workspace/manager .
COPY modules /opt/modules

USER ${USER_UID}:${USER_UID}

ENTRYPOINT ["/srv/manager"]
