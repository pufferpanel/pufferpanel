ARG DOCKER_TAG=base-devel

FROM pufferpanel/pufferpanel:${DOCKER_TAG}

# nodejs
RUN apk add --no-cache nodejs npm && \
    npm --version

# Cleanup
RUN rm -rf /var/cache/apk/*