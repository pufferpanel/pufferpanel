###
# Builder container
###
FROM golang:alpine AS builder

ARG tags=none
ARG version=devel
ARG sha=devel

ENV CGOENABLED=1

RUN go version && \
    apk add --update --no-cache gcc musl-dev git curl nodejs npm make gcc g++ python2 && \
    mkdir /pufferpanel

WORKDIR /build/pufferpanel
COPY . .
RUN go build -v -tags $tags -ldflags "-X 'github.com/pufferpanel/pufferpanel/v2.Hash=$sha' -X 'github.com/pufferpanel/pufferpanel/v2.Version=$version'" -o /pufferpanel/pufferpanel github.com/pufferpanel/pufferpanel/v2/cmd && \
    mv assets/email /pufferpanel/email && \
    cd client && \
    npm install && \
    npm run dev-build && \
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

WORKDIR /pufferpanel

ENTRYPOINT ["/pufferpanel/pufferpanel"]
CMD ["run"]
