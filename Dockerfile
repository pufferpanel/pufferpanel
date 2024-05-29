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
RUN xx-go build -buildvcs=false -tags "$tags" -ldflags "-X 'github.com/pufferpanel/pufferpanel/v3.Hash=$sha' -X 'github.com/pufferpanel/pufferpanel/v3.Version=$version'" -o /pufferpanel/pufferpanel github.com/pufferpanel/pufferpanel/v3/cmd && \
    xx-verify /pufferpanel/pufferpanel

###
# Generate final image
###

FROM alpine

EXPOSE 8080 5657
RUN mkdir -p /etc/pufferpanel && \
    mkdir -p /var/lib/pufferpanel && \
    mkdir -p /var/log/pufferpanel && \
    addgroup --system -g 1000 pufferpanel && \
    adduser -D -H --home /var/lib/pufferpanel --ingroup pufferpanel -u 1000 pufferpanel && \
    chown -R pufferpanel:pufferpanel /etc/pufferpanel /var/lib/pufferpanel /var/log/pufferpanel

ENV GIN_MODE=release \
    PUFFER_CONFIG=/etc/pufferpanel/config.json

WORKDIR /var/lib/pufferpanel

COPY --from=builder --chown=pufferpanel:pufferpanel /pufferpanel /pufferpanel/bin
COPY --from=builder --chown=pufferpanel:pufferpanel /build/pufferpanel/config.linux.json /etc/pufferpanel/config.json

USER pufferpanel

ENTRYPOINT ["/pufferpanel/bin/pufferpanel"]
CMD ["run"]
