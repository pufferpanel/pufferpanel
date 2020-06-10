###
# Builder container
###
FROM golang:alpine AS builder

ARG tags=none

ENV CGOENABLED=1

RUN go version && \
    apk add --update --no-cache gcc musl-dev git curl nodejs nodejs-npm make gcc g++ python2 && \
    mkdir /pufferpanel

WORKDIR /build/pufferpanel
COPY . .
RUN go build -v -tags $tags -o /pufferpanel/pufferpanel github.com/pufferpanel/pufferpanel/v2/cmd && \
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
VOLUME /etc/pufferpanel

ENV PUFFER_LOGS=/etc/pufferpanel/logs \
    PUFFER_WEB_HOST=0.0.0.0:8080 \
    PUFFER_PANEL_DATABASE_SESSION=60 \
    PUFFER_PANEL_DATABASE_DIALECT=sqlite3 \
    PUFFER_PANEL_DATABASE_URL="file:/etc/pufferpanel/pufferpanel.db?cache=shared" \
    PUFFER_PANEL_DATABASE_LOG=false \
    PUFFER_PANEL_TOKEN_PRIVATE=/etc/pufferpanel/private.pem \
    PUFFER_PANEL_WEB_FILES=/pufferpanel/www \
    PUFFER_PANEL_EMAIL_TEMPLATES=/pufferpanel/email/emails.json \
    PUFFER_PANEL_EMAIL_PROVIDER=debug \
    PUFFER_PANEL_SETTINGS_COMPANYNAME=PufferPanel \
    PUFFER_PANEL_SETTINGS_MASTERURL=http://localhost:8080 \
    PUFFER_DAEMON_CONSOLE_BUFFER=50 \
    PUFFER_DAEMON_CONSOLE_FORWARD=false \
    PUFFER_DAEMON_SFTP_HOST=0.0.0.0:5657 \
    PUFFER_DAEMON_SFTP_KEY=/etc/pufferpanel/sftp.key \
    PUFFER_DAEMON_AUTH_URL=http://localhost:8080 \
    PUFFER_DAEMON_AUTH_CLIENTID=none \
    PUFFER_DAEMON_DATA_CACHE=/etc/pufferpanel/cache \
    PUFFER_DAEMON_DATA_SERVERS=/etc/pufferpanel/servers \
    PUFFER_DAEMON_DATA_MODULES=/etc/pufferpanel/modules \
    PUFFER_DAEMON_DATA_CRASHLIMIT=3

WORKDIR /pufferpanel

ENTRYPOINT ["/pufferpanel/pufferpanel"]
CMD ["run"]