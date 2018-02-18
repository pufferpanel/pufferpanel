# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM gobuffalo/buffalo as builder

RUN mkdir -p $GOPATH/src/github.com/pufferpanel/pufferpanel
WORKDIR $GOPATH/src/github.com/pufferpanel/pufferpanel

# this will cache the npm install step, unless package.json changes
ADD package.json .
ADD yarn.lock .
RUN yarn install --no-progress
ADD . .
RUN go get $(go list ./... | grep -v /vendor/)
RUN buffalo build --static -o /bin/app

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

# Comment out to run the binary in "production" mode:
# ENV GO_ENV=production
ENV ADDR=0.0.0.0

WORKDIR /bin/

COPY --from=builder /bin/app .

EXPOSE 3000

# Comment out to run the migrations before running the binary:
# CMD /bin/app migrate; /bin/app
ENTRYPOINT exec /bin/app
