ARG DOCKER_TAG=base-devel

## Now our image

FROM pufferpanel/pufferpanel:${DOCKER_TAG}

ENV JAVA_HOME=/usr/lib/jvm/java-1.8-openjdk

# java
RUN apk add --no-cache openjdk8 && \
    ln -sfn /usr/lib/jvm/java-1.8-openjdk/bin/java /usr/bin/java8 && \
    ln -sfn /usr/lib/jvm/java-1.8-openjdk/bin/javac /usr/bin/javac8 && \
    ln -sfn /usr/lib/jvm/java-1.8-openjdk/bin/java /usr/bin/java && \
    ln -sfn /usr/lib/jvm/java-1.8-openjdk/bin/javac /usr/bin/javac && \
    echo "Testing Javac 8 path" && \
    javac8 -version && \
    echo "Testing Java 8 path" && \
    java8 -version && \
    echo "Testing java path" && \
    java -version && \
    echo "Testing javac path" && \
    javac -version

# Cleanup
RUN rm -rf /var/cache/apk/*