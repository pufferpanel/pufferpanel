###
# Builder container
###
FROM golang:alpine AS builder

ARG tags=none

ENV CGOENABLED=1

RUN go version && \
    apk add --update --no-cache gcc musl-dev git curl nodejs nodejs-npm make gcc g++ python && \
    mkdir /pufferpanel

WORKDIR /build/pufferpanel
COPY . .
RUN go build -v -tags $tags -o /pufferpanel/pufferpanel github.com/pufferpanel/pufferpanel/v2/cmd && \
    mv assets/email /pufferpanel/email

WORKDIR /build/pufferpanel/client
RUN npm install && \
    npm run dev-build && \
    mv dist /pufferpanel/www/

###
# Generate final image
###

FROM alpine
COPY --from=builder /pufferpanel /pufferpanel

EXPOSE 8080
VOLUME /etc/pufferpanel

ENV PUFFERPANEL_DATABASE_SESSION=60 \
    PUFFERPANEL_DATABASE_DIALECT=sqlite3 \
    PUFFERPANEL_DATABASE_URL="file:/etc/pufferpanel/pufferpanel.db?cache=shared" \
    PUFFERPANEL_TOKEN_PRIVATE=/etc/pufferpanel/private.pem \
    PUFFERPANEL_TOKEN_PUBLIC=/etc/pufferpanel/public.pem \
    PUFFERPANEL_WEB_HOST=0.0.0.0:8080 \
    PUFFERPANEL_WEB_FILES=/pufferpanel/web \
    PUFFERPANEL_EMAIL_TEMPLATES=/pufferpanel/email/emails.json \
    PUFFERPANEL_EMAIL_PROVIDER=debug \
    PUFFERPANEL_SETTINGS_COMPANYNAME=PufferPanel \
    PUFFERPANEL_SETTINGS_MASTERURL=http://localhost:8080 \
    PUFFERPANEL_LOGS=/etc/pufferpanel/logs

WORKDIR /pufferpanel

ENTRYPOINT ["/pufferpanel/pufferpanel"]
CMD ["run"]