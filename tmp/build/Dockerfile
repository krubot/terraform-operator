FROM alpine:3.6

RUN adduser -D app-operator
USER app-operator

ADD tmp/_output/bin/app-operator /usr/local/bin/app-operator
