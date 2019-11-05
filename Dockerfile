###
# Builder container
###
FROM golang:alpine AS builder

ARG tags=none

ENV CGOENABLED=1

RUN go version && \
    apk add --update --no-cache gcc musl-dev git curl nodejs nodejs-npm make gcc g++ python && \
    mkdir /pufferpanel && \
    wget https://github.com/swaggo/swag/releases/download/v1.6.3/swag_1.6.3_Linux_x86_64.tar.gz && \
    tar -zxf swag*.tar.gz && \
    mv swag /go/bin/

WORKDIR /build/apufferi
RUN git clone https://github.com/pufferpanel/apufferi /build/apufferi

#swag init --parseDependency -g web/api/loader.go && \
WORKDIR /build/pufferpanel
COPY . .
RUN echo replace github.com/pufferpanel/apufferi/v4 =\> ../apufferi >> go.mod && \
    go get -u github.com/pufferpanel/pufferpanel/v2/cmd && \
    go build -v -tags $tags -o /pufferpanel/pufferpanel github.com/pufferpanel/pufferpanel/v2/cmd && \
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

ENV PATH=$PATH:/pufferpanel

WORKDIR /pufferpanel

EXPOSE 8080
CMD ["/pufferpanel/pufferpanel", "run"]