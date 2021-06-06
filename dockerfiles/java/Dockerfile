ARG DOCKER_TAG=devel

FROM pufferpanel/pufferpanel:${DOCKER_TAG}

# enable testing repo

RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    apk update

# java
RUN apk add --no-cache openjdk8 openjdk16 && \
    export JAVA8_BIN=/usr/lib/jvm/java-1.8-openjdk/bin/java && \
    export JAVAC8_BIN=/usr/lib/jvm/java-1.8-openjdk/bin/javac && \
    export JAVA16_BIN=/usr/lib/jvm/java-16-openjdk/bin/java && \
    export JAVAC16_BIN=/usr/lib/jvm/java-16-openjdk/bin/javac && \
    ln -sfn $JAVA8_BIN /usr/bin/java && \
    ln -sfn $JAVAC8_BIN /usr/bin/javac && \
    echo "Testing Javac 8 ENV variable" && \
    $JAVAC8_BIN -version && \
    echo "Testing Java 8 ENV variable" && \
    $JAVA8_BIN -version && \
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