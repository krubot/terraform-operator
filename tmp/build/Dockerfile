FROM alpine:3.6

RUN adduser -D terraform-operator
USER terraform-operator

ADD tmp/_output/bin/terraform-operator /usr/local/bin/terraform-operator
