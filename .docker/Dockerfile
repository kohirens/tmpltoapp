ARG USER_NAME='app'
ARG USER_UID='1000'
ARG USER_GID='1000'
ARG USER_GROUP='app_users'
ARG REPO='github.com/kohirens/bootup'

FROM golang:1.16-alpine3.14 AS dev

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP
ARG REPO

ENV GOPATH /home/${USER_NAME}
ENV CGO_ENABLED=0
ENV WORK_DIR=/home/${USER_NAME}/src/${REPO}
ENV SHELL=/bin/sh

COPY .devcontainer/download-vs-code-server.sh .

# VS Code Requirements: openssh, musl, libgcc, libstdc++
RUN apk --no-progress --purge --no-cache upgrade \
 && apk --no-progress --purge --no-cache add --upgrade \
    curl \
    git \
    gnupg \
    libgcc \
    libstdc++ \
    openssh \
    tini \
 && apk --no-progress --purge --no-cache upgrade \
 && rm -vrf /var/cache/apk/* \
 && curl --version \
 && git --version

# Install vanilla GLibC: https://github.com/sgerrand/alpine-pkg-glibc
RUN curl -o /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
 && curl -LO https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.32-r0/glibc-2.32-r0.apk \
 && apk add glibc-2.32-r0.apk

# Add a non-root group and user, helpful if you dev on Linux.
RUN addgroup --system --gid ${USER_GID} ${USER_GROUP} \
 && adduser --system \
    --disabled-password \
    --ingroup ${USER_GROUP} \
    --uid ${USER_UID} \
    ${USER_NAME}

USER ${USER_NAME}

# Make directories with correct permissions
RUN mkdir -p ~/bin

# Install Go helpful dev tools.
RUN go get -u -v golang.org/x/tools/gopls \
    github.com/go-delve/delve/cmd/dlv \
    github.com/ramya-rao-a/go-outline \
    github.com/cweill/gotests/... \
    github.com/uudashr/gopkgs/v2/cmd/gopkgs \
    honnef.co/go/tools/cmd/staticcheck \
    github.com/git-chglog/git-chglog/cmd/git-chglog

# Install VS Code Server and Requirements for pre-installing extensions
RUN ./download-vs-code-server.sh \
 && mkdir -p /home/${USER_NAME}/.vscode-server/extensions \
 && mkdir -p /home/${USER_NAME}/.vscode-server-insiders/extensions

ENV PATH="/home/${USER_NAME}/bin:${PATH}"

WORKDIR "/home/${USER_NAME}/src/${REPO}"

ENTRYPOINT [ "tini", "--" ]

CMD [ "sh" ]

FROM dev AS build

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP
ARG REPO

ENV GOPATH /home/${USER_NAME}
ENV CGO_ENABLED=0

COPY . .

RUN go build && go install

FROM alpine:3.13 AS release

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP

# Add a non-root group and user, helpful if you dev on Linux.
RUN addgroup --system --gid ${USER_GID} ${USER_GROUP} \
 && adduser --system \
    --disabled-password \
    --ingroup ${USER_GROUP} \
    --uid ${USER_UID} \
    ${USER_NAME}

USER ${USER_NAME}

# Make directories with correct permissions
RUN mkdir -p ~/bin

COPY --from=build --chown=${USER_NAME}:${USER_GROUP} /home/${USER_NAME}/bin/bootup /home/${USER_NAME}/bin

ENV PATH="/home/${USER_NAME}/bin:${PATH}"

ENTRYPOINT [ "bootup" ]

CMD [ "-h" ]

# docker build -t bootup --rm -f .docker/Dockerfile --progress plain .
# docker run --rm bootup
