FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/terraform-operator \
    USER_UID=0 \
    USER_NAME=terraform-operator

# install operator binary
COPY build/_output/bin/terraform-operator ${OPERATOR}
COPY build/bin /usr/local/bin
COPY modules /var/lib/modules

RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

WORKDIR /srv
USER ${USER_UID}
