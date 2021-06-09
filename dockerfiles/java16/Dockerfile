ARG DOCKER_TAG=base-devel

## Now our image

FROM pufferpanel/pufferpanel:${DOCKER_TAG}

# enable testing repo

RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    apk update

ENV JAVA_HOME=/usr/lib/jvm/java-16-openjdk

# java
RUN apk add --no-cache openjdk16 && \
    export JAVA16_BIN=/usr/lib/jvm/java-16-openjdk/bin/java && \
    export JAVAC16_BIN=/usr/lib/jvm/java-16-openjdk/bin/javac && \
    ln -sfn $JAVAC16_BIN /usr/bin/javac && \
    echo "Testing Javac 16 ENV variable" && \
    $JAVAC16_BIN -version && \
    echo "Testing Java 16 ENV variable" && \
    $JAVA16_BIN -version && \
    echo "Testing java in path" && \
    java -version && \
    echo "Testing javac in path" && \
    javac -version

# Cleanup
RUN rm -rf /var/cache/apk/*