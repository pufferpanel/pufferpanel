FROM golang:alpine

RUN apk add --update --no-cache gcc musl-dev git curl nodejs nodejs-npm make gcc g++ python && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/pufferpanel/pufferpanel/

COPY . .

ENV CGOENABLED=1

RUN mkdir -p /etc/pufferpanel && \
    echo "{}" > /etc/pufferpanel/config.json

RUN dep ensure -v && \
    go install -v github.com/pufferpanel/pufferpanel

WORKDIR /go/src/github.com/pufferpanel/pufferpanel/client
RUN npm install && \
    npm run dev-build && \
    mkdir /go/bin/client && \
    mkdir -p /pufferpanel/client/ && \
    mv dist /pufferpanel/client/

WORKDIR /pufferpanel
RUN rm -rf /go/src/github.com/pufferpanel && \
    apk del gcc musl-dev git curl nodejs nodejs-npm make gcc g++ python

EXPOSE 8080
CMD ["pufferpanel"]