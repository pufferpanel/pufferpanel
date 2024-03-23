ARG DOCKER_TAG=base-devel

FROM pufferpanel/pufferpanel:${DOCKER_TAG}

# enable testing repo

RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories && \
    apk update

ENV JAVA_HOME=/usr/lib/jvm/java-17-openjdk

# java
RUN apk add --no-cache openjdk8 openjdk17 && \
    ln -sfn /usr/lib/jvm/java-17-openjdk/bin/java /usr/bin/java && \
    ln -sfn /usr/lib/jvm/java-17-openjdk/bin/javac /usr/bin/javac && \
    ln -sfn /usr/lib/jvm/java-17-openjdk/bin/java /usr/bin/java17 && \
    ln -sfn /usr/lib/jvm/java-17-openjdk/bin/javac /usr/bin/javac17 && \
    ln -sfn /usr/lib/jvm/java-1.8-openjdk/bin/java /usr/bin/java8 && \
    ln -sfn /usr/lib/jvm/java-1.8-openjdk/bin/javac /usr/bin/javac8 && \
    echo "Testing Javac 8 path" && \
    javac8 -version && \
    echo "Testing Java 8 path" && \
    java8 -version && \
    echo "Testing Javac 17 path" && \
    javac17 -version && \
    echo "Testing Java 17 path" && \
    java17 -version && \
    echo "Testing java path" && \
    java -version && \
    echo "Testing javac path" && \
    javac -version

# Cleanup
RUN rm -rf /var/cache/apk/*