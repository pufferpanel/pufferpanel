###
# Builder container
###
FROM --platform=$BUILDPLATFORM node:20-alpine AS node

WORKDIR /build
COPY client .

RUN rm -rf /build/*/node_modules/ && \
    rm -rf /build/*/dist/

RUN yarn install && \
    yarn build

FROM --platform=$BUILDPLATFORM tonistiigi/xx AS xx

FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS builder

RUN apk add clang lld
COPY --from=xx / /

ARG tags=nohost
ARG version=devel
ARG sha=devel
ARG swagversion=1.16.2
ARG swagarch=x86_64

ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

RUN mkdir /pufferpanel && \
    wget https://github.com/swaggo/swag/releases/download/v${swagversion}/swag_${swagversion}_Linux_$swagarch.tar.gz && \
    mkdir -p ~/go/bin && \
    tar -zxvf swag*.tar.gz -C ~/go/bin && \
    rm -rf swag*.tar.gz

WORKDIR /build/pufferpanel

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN ~/go/bin/swag init --md . -o web/swagger -g web/loader.go

COPY --from=node /build/frontend/dist /build/pufferpanel/client/frontend/dist

ARG TARGETPLATFORM
RUN xx-apk add musl-dev gcc
RUN xx-go build -buildvcs=false -tags "$tags" -ldflags "-X 'github.com/pufferpanel/pufferpanel/v3.Hash=$sha' -X 'github.com/pufferpanel/pufferpanel/v3.Version=$version'" -o /pufferpanel/pufferpanel github.com/pufferpanel/pufferpanel/v3/cmd
RUN xx-verify /pufferpanel/pufferpanel

###
# Generate final image
###

FROM alpine

EXPOSE 8080 5657
RUN mkdir -p /etc/pufferpanel && \
    mkdir -p /var/lib/pufferpanel /var/lib/pufferpanel/servers /var/lib/pufferpanel/binaries /var/lib/pufferpanel/cache && \
    mkdir -p /var/log/pufferpanel
    #addgroup --system -g 1000 pufferpanel && \
    #adduser -D -H --home /var/lib/pufferpanel --ingroup pufferpanel -u 1000 pufferpanel && \
    #chown -R pufferpanel:pufferpanel /etc/pufferpanel /var/lib/pufferpanel /var/log/pufferpanel

ENV GIN_MODE=release \
    PUFFER_DOCKER_ROOT=""

#COPY --from=builder --chown=pufferpanel:pufferpanel --chmod=755 /pufferpanel /pufferpanel/bin
#COPY --from=builder --chown=pufferpanel:pufferpanel --chmod=755 /build/pufferpanel/entrypoint.sh /pufferpanel/bin/entrypoint.sh
#COPY --from=builder --chown=pufferpanel:pufferpanel --chmod=755 /build/pufferpanel/config.docker.json /etc/pufferpanel/config.json
COPY --from=builder --chmod=755 /pufferpanel /pufferpanel/bin
COPY --from=builder --chmod=755 /build/pufferpanel/entrypoint.sh /pufferpanel/bin/entrypoint.sh
COPY --from=builder --chmod=755 /build/pufferpanel/config.docker.json /etc/pufferpanel/config.json

VOLUME /etc/pufferpanel
VOLUME /var/lib/pufferpanel
VOLUME /var/log/pufferpanel

WORKDIR /var/lib/pufferpanel

#USER pufferpanel

RUN /pufferpanel/bin/pufferpanel dbmigrate

ENTRYPOINT ["sh", "/pufferpanel/bin/entrypoint.sh"]
