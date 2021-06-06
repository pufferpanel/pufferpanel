ARG DOCKER_TAG=devel

## SRCDS is dumb, so we have to do crazy stuff
# copied from https://github.com/steamcmd/docker/blob/master/dockerfiles/alpine-3/Dockerfile
FROM steamcmd/steamcmd:ubuntu-18 as srcdsbuilder
ENV USER root
ENV HOME /root/installer
WORKDIR $HOME
RUN apt-get update \
 && apt-get install -y --no-install-recommends curl tar
RUN curl http://media.steampowered.com/installer/steamcmd_linux.tar.gz \
    --output steamcmd.tar.gz --silent
RUN tar -xvzf steamcmd.tar.gz && rm steamcmd.tar.gz

## Now our image

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

# nodejs
RUN apk add --no-cache nodejs npm && \
    npm --version

# srcds
RUN apk add --no-cache bash
COPY --from=srcdsbuilder /lib/i386-linux-gnu /lib/
COPY --from=srcdsbuilder /root/installer/linux32/libstdc++.so.6 /lib/
COPY --from=srcdsbuilder /root/installer/steamcmd.sh /usr/lib/games/steam/
COPY --from=srcdsbuilder /root/installer/linux32/steamcmd /usr/lib/games/steam/
COPY --from=srcdsbuilder /usr/games/steamcmd /usr/bin/steamcmd

# Cleanup
RUN rm -rf /var/cache/apk/*