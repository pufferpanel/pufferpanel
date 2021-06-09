ARG DOCKER_TAG=base-devel

## Now our image

FROM pufferpanel/pufferpanel:${DOCKER_TAG}

ENV JAVA_HOME=/usr/lib/jvm/java-1.8-openjdk

# java
RUN apk add --no-cache openjdk8 && \
    export JAVA8_BIN=/usr/lib/jvm/java-1.8-openjdk/bin/java && \
    export JAVAC8_BIN=/usr/lib/jvm/java-1.8-openjdk/bin/javac && \
    ln -sfn $JAVAC8_BIN /usr/bin/javac && \
    echo "Testing Javac 8 ENV variable" && \
    $JAVAC8_BIN -version && \
    echo "Testing Java 8 ENV variable" && \
    $JAVA8_BIN -version && \
    echo "Testing java in path" && \
    java -version && \
    echo "Testing javac in path" && \
    javac -version

# Cleanup
RUN rm -rf /var/cache/apk/*