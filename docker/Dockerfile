FROM golang:1.15.4-alpine3.12 AS dev

ARG USER_NAME='gitter'
ARG USER_UID='1000'
ARG USER_GID='1000'
ARG USER_GROUP='app_users'

ENV GOPATH /home/$USER_NAME
ENV GOFLAGS -mod=vendor
ENV CGO_ENABLED=0

# VSCode Requirements: "openssh"
RUN apk --no-progress --purge --no-cache upgrade \
 && apk --no-progress --purge --no-cache add --upgrade \
    curl \
    git \
    openssh \
    tini \
 && apk --no-progress --purge --no-cache upgrade \
 && rm -vrf /var/cache/apk/* \
 && curl --version \
 && git --version

# Add a non-root group and user, helpful if you dev on Linux.
RUN addgroup --system --gid $USER_GID $USER_GROUP \
 && adduser --system \
    --disabled-password \
    --ingroup $USER_GROUP \
    --uid $USER_UID \
    $USER_NAME

USER gitter

RUN go get -u github.com/git-chglog/git-chglog/cmd/git-chglog \
 && go get -u -v golang.org/x/tools/gopls \
 && go get -u -v github.com/go-delve/delve/cmd/dlv \
 && go get -u -v github.com/ramya-rao-a/go-outline

# VSCode Requirements for pre-installing extensions
RUN mkdir -p /home/$USER_NAME/.vscode-server/extensions /home/$USER_NAME/.vscode-server-insiders/extensions
#  && chown -R $USER_NAME /home/$USER_NAME/.vscode-server /home/$USER_NAME/.vscode-server-insiders

WORKDIR /home/gitter/src/github.com/kohirens/go-gitter

ENTRYPOINT [ "tini", "--" ]

CMD [ "tail", "-f", "/dev/null" ]

FROM dev AS build

ENV GOPATH /home/$USER

COPY ./ /home/$USER/src
RUN whoami \
 && ls -la .

RUN go build

FROM alpine3.12 AS release

COPY  --chown=$USER_NAME:$USER_GROUP --from=build /home/$USER/bin /usr/local/bin

ENTRYPOINT [ "gitter" ]

CMD [ "-h" ]
