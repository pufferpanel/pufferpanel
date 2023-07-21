###
# Builder container
###
FROM tinygo/tinygo:0.28.1 AS wasm

ARG tinygoversion=0.28.1
ENV GOOS=js \
    GOARCH=wasm

USER root

RUN  apt-get update && \
     apt-get install wget && \
     rm -rf /var/lib/apt/lists/* && \
     wget https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-16/libclang_rt.builtins-wasm32-wasi-16.0.tar.gz && \
     tar -zxf libclang_rt.builtins-wasm32-wasi-16.0.tar.gz lib/wasi/libclang_rt.builtins-wasm32.a && \
     rm libclang_rt.builtins-wasm32-wasi-16.0.tar.gz && \
     mv lib/wasi/libclang_rt.builtins-wasm32.a /usr/local/tinygo/lib/wasi-libc/sysroot/lib/wasm32-wasi/

WORKDIR /build

COPY wasm.json.docker.patch ./
RUN patch -ul /usr/local/tinygo/targets/wasm.json wasm.json.docker.patch

COPY go.mod go.sum ./
RUN go mod download && go mod verify

ENV GOFLAGS="-buildvcs=false"
COPY . .
RUN tinygo build -target=wasm -no-debug -o conditions.wasm github.com/pufferpanel/pufferpanel/v3/conditions/wasm

FROM node:18-alpine AS node
FROM golang:1.20-alpine AS builder

COPY --from=node /usr/lib /usr/lib
COPY --from=node /usr/local/share /usr/local/share
COPY --from=node /usr/local/lib /usr/local/lib
COPY --from=node /usr/local/include /usr/local/include
COPY --from=node /usr/local/bin /usr/local/bin
COPY --from=node /opt /opt

ARG tags=docker
ARG version=devel
ARG sha=devel
ARG goproxy
ARG npmproxy
ARG swagversion=1.8.10

ENV CGOENABLED=1 \
    YARN_REGISTRY=$npmproxy \
    GOPROXY=$goproxy

RUN go version && \
    apk add --update --no-cache gcc musl-dev git curl make gcc g++ && \
    mkdir /pufferpanel && \
    wget https://github.com/swaggo/swag/releases/download/v${swagversion}/swag_${swagversion}_Linux_x86_64.tar.gz && \
    mkdir -p ~/go/bin && \
    tar -zxf swag*.tar.gz -C ~/go/bin && \
    rm -rf swag*.tar.gz

WORKDIR /build/pufferpanel

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

COPY --from=wasm /build/conditions.wasm /build/pufferpanel/client/frontend/public

RUN cd client && \
    yarn install && \
    yarn build

RUN ~/go/bin/swag init -o web/swagger -g web/loader.go && \
    go build -v -buildvcs=false -tags "$tags" -ldflags "-X 'github.com/pufferpanel/pufferpanel/v3.Hash=$sha' -X 'github.com/pufferpanel/pufferpanel/v3.Version=$version'" -o /pufferpanel/pufferpanel github.com/pufferpanel/pufferpanel/v3/cmd

###
# Generate final image
###

FROM alpine
COPY --from=builder /pufferpanel /pufferpanel

EXPOSE 8080 5657
RUN mkdir -p /etc/pufferpanel && \
    mkdir -p /var/lib/pufferpanel

ENV PUFFER_LOGS=/etc/pufferpanel/logs \
    PUFFER_PANEL_TOKEN_PUBLIC=/etc/pufferpanel/public.pem \
    PUFFER_PANEL_TOKEN_PRIVATE=/etc/pufferpanel/private.pem \
    PUFFER_PANEL_DATABASE_DIALECT=sqlite3 \
    PUFFER_PANEL_DATABASE_URL="file:/etc/pufferpanel/pufferpanel.db?cache=shared" \
    PUFFER_DAEMON_SFTP_KEY=/etc/pufferpanel/sftp.key \
    PUFFER_DAEMON_DATA_CACHE=/var/lib/pufferpanel/cache \
    PUFFER_DAEMON_DATA_SERVERS=/var/lib/pufferpanel/servers \
    PUFFER_DAEMON_DATA_MODULES=/var/lib/pufferpanel/modules \
    PUFFER_DAEMON_DATA_BINARIES=/var/lib/pufferpanel/binaries \
    GIN_MODE=release

WORKDIR /pufferpanel

ENTRYPOINT ["/pufferpanel/pufferpanel"]
CMD ["run"]
