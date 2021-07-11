FROM golang:1.13-alpine
ARG MAI_UID=1000
ARG MAI_GID=1000
ENV SERVICE_NAME mnc-interview
ENV GOPRIVATE gitlab.com/pt-mai/*

COPY build/assets/ghz /usr/local/bin
COPY build/assets/repositories /etc/apk/repositories

RUN chmod 644 /etc/apk/repositories && \
    apk update && apk add --no-cache sudo git
RUN apk add --no-cache make tzdata curl && \
    chmod a+x /usr/local/bin/ghz

COPY build/assets/sudoers /etc/sudoers

RUN chmod 0400 /etc/sudoers && \
    addgroup mai -g ${MAI_GID} && \
    addgroup sudo && \
    adduser -h /home/mai -s /bin/sh -u ${MAI_UID} -G users -D mai && \
    adduser mai mai && \
    adduser mai sudo

USER mai
WORKDIR "/go/src/gitlab.com/pt-mai/${SERVICE_NAME}"

COPY --chown=mai:users .netrc /tmp/.netrc

COPY --chown=mai:users build/dev/evans_linux_amd64.tar.gz /tmp

RUN sudo chown -R mai:mai /go && \
    cp /tmp/.netrc ~/.netrc && \
    chmod 600 ~/.netrc && \
    sudo tar -C /usr/bin -xzf /tmp/evans_linux_amd64.tar.gz && \
    sudo rm /tmp/evans_linux_amd64.tar.gz && \
    export GO111MODULE=off && \
    echo "Installing additional tools, please wait..." && \
    go get -u golang.org/x/lint/golint && echo "- Success installing golint" && \
    go get github.com/githubnemo/CompileDaemon && echo "- Success installing CompileDaemon" && \
    go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate && echo "- Success installing golang-migrate" && \
    go get github.com/bayucandra/go-seed-pg && echo "- Success installing go-seed-pg" && \
    unset GO111MODULE


ENTRYPOINT ["sh", "./build/dev/entrypoint.sh"]
