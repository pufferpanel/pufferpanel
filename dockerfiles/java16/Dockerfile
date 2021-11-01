ARG DOCKER_TAG=base-devel

## Now our image

FROM pufferpanel/pufferpanel:${DOCKER_TAG}

# enable testing repo

RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories && \
    apk update

ENV JAVA_HOME=/usr/lib/jvm/java-16-openjdk

# java
RUN apk add --no-cache openjdk16 && \
    ln -sfn /usr/lib/jvm/java-16-openjdk/bin/java /usr/bin/java && \
    ln -sfn /usr/lib/jvm/java-16-openjdk/bin/javac /usr/bin/javac && \
    ln -sfn /usr/lib/jvm/java-16-openjdk/bin/java /usr/bin/java16 && \
    ln -sfn /usr/lib/jvm/java-16-openjdk/bin/javac /usr/bin/javac16 && \
    echo "Testing Javac 16 path" && \
    javac16 -version && \
    echo "Testing Java 16 path" && \
    java16 -version && \
    echo "Testing java path" && \
    java -version && \
    echo "Testing javac path" && \
    javac -version

# Cleanup
RUN rm -rf /var/cache/apk/*