###
# Builder container
###
FROM golang:1.18-alpine AS builder

ARG tags=none
ARG version=devel
ARG sha=devel
ARG goproxy
ARG npmproxy

ENV CGOENABLED=1

ENV npm_config_registry=$npmproxy
ENV GOPROXY=$goproxy

RUN go version && \
    apk add --update --no-cache gcc musl-dev git curl nodejs npm make gcc g++ && \
    mkdir /pufferpanel && \
    wget https://github.com/swaggo/swag/releases/download/v1.6.7/swag_1.6.7_Linux_x86_64.tar.gz && \
    mkdir -p ~/go/bin && \
    tar -zxf swag*.tar.gz -C ~/go/bin && \
    rm -rf swag*.tar.gz

WORKDIR /build/pufferpanel

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN ~/go/bin/swag init -o web/swagger -g web/loader.go && \
    go build -v -buildvcs=false -tags $tags -ldflags "-X 'github.com/pufferpanel/pufferpanel/v2.Hash=$sha' -X 'github.com/pufferpanel/pufferpanel/v2.Version=$version'" -o /pufferpanel/pufferpanel github.com/pufferpanel/pufferpanel/v2/cmd && \
    mv assets/email /pufferpanel/email && \
    cd client && \
    npm install && \
    npm run build && \
    mv dist /pufferpanel/www/

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
    GIN_MODE=release

WORKDIR /pufferpanel

ENTRYPOINT ["/pufferpanel/pufferpanel"]
CMD ["run"]
