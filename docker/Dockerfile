FROM golang:1.15.4-alpine3.12 AS dev

ARG USER_NAME='gitter'
ARG USER_UID='1000'
ARG USER_GID='1000'
ARG USER_GROUP='app_users'

COPY ./docker/rsyslog.conf /etc/rsyslog.conf

ENV GOPATH /home/$USER_NAME

RUN apk --no-progress --purge --no-cache upgrade \
 && apk --no-progress --purge --no-cache add --upgrade \
    curl \
    git \
    rsyslog \
 && apk --no-progress --purge --no-cache upgrade \
 && rm -vrf /var/cache/apk/* \
 && curl --version \
 && git --version

RUN apk --no-progress --purge --no-cache upgrade \
 && apk --no-progress --purge --no-cache add --upgrade \
 && go get -u github.com/git-chglog/git-chglog/cmd/git-chglog \
 && go get -u -v golang.org/x/tools/gopls

# Add a non-root group and user, helpful if you dev on Linux.
RUN addgroup --system --gid $USER_GID $USER_GROUP \
 && adduser --system \
    --disabled-password \
    --ingroup $USER_GROUP \
    --uid $USER_UID \
    $USER_NAME \
 && adduser $USER_NAME daemon

USER gitter

WORKDIR /home/gitter/src/github.com/kohirens/go-gitter

# ENTRYPOINT [ "syslogd" ]
ENTRYPOINT [ "rsyslogd" ]

CMD [ "-n" ]

FROM dev AS build

ENV GOPATH /home/$USER

COPY ./ /home/$USER/src
RUN whoami \
 && ls -la .

RUN go build

FROM alpine3.12 AS release

COPY --from=build /home/$USER/bin /usr/local/bin

ENTRYPOINT [ "gitter" ]

CMD [ "" ]
