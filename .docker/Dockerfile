ARG USER_NAME='app'
ARG USER_UID='1000'
ARG USER_GID='1000'
ARG USER_GROUP='app_users'
ARG REPO='github.com/kohirens/tmpltoapp'
ARG GLIBC_VER='2.35-r0'

FROM golang:1.19-alpine3.16 AS dev

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP
ARG REPO
ARG GLIBC_VER

ENV GOPATH /home/${USER_NAME}
ENV CGO_ENABLED=0
ENV WORK_DIR=/home/${USER_NAME}/src/${REPO}
ENV SHELL=/bin/sh
ENV LD_LIBRARY_PATH="/usr/lib:/lib"

COPY .docker/start.sh /usr/local/bin/start.sh

# VS Code Requirements: openssh, musl, libgcc, libstdc++
RUN apk --no-progress --purge --no-cache upgrade \
 && apk --no-progress --purge --no-cache add --upgrade \
    curl \
    git \
    gnupg \
    libgcc \
    libstdc++ \
    openssh \
 && apk --no-progress --purge --no-cache upgrade \
 && rm -vrf /var/cache/apk/* \
 && rm -rf /tmp/* \
 && curl --version \
 && git --version

# Install vanilla GLibC: https://github.com/sgerrand/alpine-pkg-glibc
RUN curl -o /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
 && curl -LO https://github.com/sgerrand/alpine-pkg-glibc/releases/download/${GLIBC_VER}/glibc-${GLIBC_VER}.apk \
 && apk add --force-overwrite glibc-${GLIBC_VER}.apk \
 && rm glibc-${GLIBC_VER}.apk \
       /etc/apk/keys/sgerrand.rsa.pub

# Add a non-root group and user, helpful if you dev on Linux.
RUN addgroup --system --gid ${USER_GID} ${USER_GROUP} \
 && adduser --system \
    --disabled-password \
    --ingroup ${USER_GROUP} \
    --uid ${USER_UID} \
    ${USER_NAME}

USER ${USER_NAME}

# Make directories with correct permissions
RUN mkdir -p ~/bin ~/src/${REPO} && chmod 0774 -R ~/src/${REPO}

WORKDIR "/home/${USER_NAME}"

# Install Go helpful dev tools.
RUN go install golang.org/x/tools/gopls@latest \
 && go install github.com/go-delve/delve/cmd/dlv@latest \
 && go install github.com/ramya-rao-a/go-outline@latest \
 && go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest \
 && go install honnef.co/go/tools/cmd/staticcheck@latest \
 && go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest \
 && go install github.com/cweill/gotests/...@latest

COPY --chown=${USER_NAME}:${USER_GROUP} .devcontainer/download-vs-code-server.sh .

# Install VS Code Server and Requirements for pre-installing extensions
RUN chmod +x ./download-vs-code-server.sh && ./download-vs-code-server.sh \
 && mkdir -p /home/${USER_NAME}/.vscode-server/extensions \
 && mkdir -p /home/${USER_NAME}/.vscode-server-insiders/extensions \
 && rm ./download-vs-code-server.sh

ENV PATH="/home/${USER_NAME}/bin:${PATH}"

WORKDIR "/home/${USER_NAME}/src/${REPO}"

ENTRYPOINT [ "sh" ]

CMD [ "/usr/local/bin/start.sh" ]

FROM dev AS build

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP
ARG REPO

COPY --chown=${USER_NAME}:${USER_GROUP} . .

RUN go mod tidy && go build && go install

FROM alpine:3.16 AS release

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

COPY --from=build --chown=${USER_NAME}:${USER_GROUP} /home/${USER_NAME}/bin/tmpltoapp /home/${USER_NAME}/bin

ENV PATH="/home/${USER_NAME}/bin:${PATH}"

ENTRYPOINT [ "tmpltoapp" ]

CMD [ "-h" ]

# docker build -t tmpltoapp --rm -f .docker/Dockerfile --progress plain .
# docker run --rm tmpltoapp
