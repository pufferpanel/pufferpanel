###
# Builder container
###
FROM golang:alpine AS builder

ARG tags=none

ENV CGOENABLED=1

RUN go version && \
    apk add --update --no-cache gcc musl-dev git curl nodejs nodejs-npm make gcc g++ python && \
    mkdir /pufferpanel

WORKDIR /build/apufferi
RUN git clone https://github.com/pufferpanel/apufferi /build/apufferi

WORKDIR /build/pufferpanel
COPY . .
RUN echo replace github.com/pufferpanel/apufferi/v4 =\> ../apufferi >> go.mod
RUN go build -v -tags $tags -o /pufferpanel/pufferpanel github.com/pufferpanel/pufferpanel/v2/cmd

WORKDIR /build/pufferpanel/client
RUN npm install && \
    npm run dev-build && \
    mv dist /pufferpanel/www/

###
# Generate final image
###

FROM alpine
COPY --from=builder /pufferpanel /pufferpanel

ENV PATH=$PATH:/pufferpanel

WORKDIR /pufferpanel

EXPOSE 8080
CMD ["/pufferpanel/pufferpanel", "run"]