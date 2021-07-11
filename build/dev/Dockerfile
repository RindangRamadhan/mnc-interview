FROM golang:1.13-alpine
ARG RAMA_UID=1000
ARG RAMA_GID=1000
ENV SERVICE_NAME mnc-interview

COPY build/assets/ghz /usr/local/bin
COPY build/assets/repositories /etc/apk/repositories

RUN chmod 644 /etc/apk/repositories && \
    apk update && apk add --no-cache sudo git
RUN apk add --no-cache make tzdata curl && \
    chmod a+x /usr/local/bin/ghz

COPY build/assets/sudoers /etc/sudoers

RUN chmod 0400 /etc/sudoers && \
    addgroup rama -g ${RAMA_GID} && \
    addgroup sudo && \
    adduser -h /home/rama -s /bin/sh -u ${RAMA_UID} -G users -D rama && \
    adduser rama rama && \
    adduser rama sudo

USER rama
WORKDIR "/go/src/github.com/rindangramadhan/${SERVICE_NAME}"

COPY --chown=rama:users .netrc /tmp/.netrc

COPY --chown=rama:users build/dev/evans_linux_amd64.tar.gz /tmp

RUN sudo chown -R rama:rama /go && \
    cp /tmp/.netrc ~/.netrc && \
    chmod 600 ~/.netrc && \
    export GO111MODULE=off && \
    echo "Installing additional tools, please wait..." && \
    go get -u golang.org/x/lint/golint && echo "- Success installing golint" && \
    go get github.com/githubnemo/CompileDaemon && echo "- Success installing CompileDaemon" && \
    unset GO111MODULE


ENTRYPOINT ["sh", "./build/dev/entrypoint.sh"]
