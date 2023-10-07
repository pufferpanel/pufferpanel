###
# Builder container
###
FROM node:18-alpine AS node

WORKDIR /build
COPY client .

RUN yarn install && \
    yarn build

FROM golang:1.21-alpine AS builder

ARG tags=nohost
ARG version=devel
ARG sha=devel
ARG swagversion=1.16.2

ENV CGOENABLED=1

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

RUN ~/go/bin/swag init --md . -o web/swagger -g web/loader.go

COPY --from=node /build/frontend/dist /build/pufferpanel/client/frontend/dist

RUN go build -v -buildvcs=false -tags "$tags" -ldflags "-X 'github.com/pufferpanel/pufferpanel/v3.Hash=$sha' -X 'github.com/pufferpanel/pufferpanel/v3.Version=$version'" -o /pufferpanel/pufferpanel github.com/pufferpanel/pufferpanel/v3/cmd && \
    go test -v ./...

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
