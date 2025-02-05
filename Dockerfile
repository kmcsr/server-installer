# syntax=docker/dockerfile:1.2

ARG GO_VERSION=1.20
ARG REPO=github.com/kmcsr/server-installer
ARG SUB_FOLDER=minecraft_installer

FROM golang:${GO_VERSION}-alpine AS BUILD

ARG REPO
ARG SUB_FOLDER

COPY ./go.mod ./go.sum "/go/src/${REPO}/"
COPY "./*.go" "/go/src/${REPO}/"
COPY "./$SUB_FOLDER" "/go/src/${REPO}/${SUB_FOLDER}"
RUN --mount=type=cache,target=/root/.cache/go-build cd "/go/src/${REPO}" && \
 CGO_ENABLED=0 go build -v -o "/go/bin/minecraft_installer" "./${SUB_FOLDER}"

FROM alpine:latest

COPY --from=BUILD "/go/bin/minecraft_installer" "/usr/local/bin/minecraft_installer"

CMD [ "/usr/local/bin/minecraft_installer" ]
