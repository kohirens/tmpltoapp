ARG REPO='{{REPO}}'

FROM alpine:latest as dev

ENV WORK_DIR=/home/${USER_NAME}/src/${REPO}

ENTRYPOINT [ "tail" ]

CMD [ "-f", "/dev/null" ]
