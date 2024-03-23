ARG DOCKER_TAG=base-devel

FROM pufferpanel/pufferpanel:${DOCKER_TAG} AS builder

FROM ubuntu:20.04

COPY --from=builder /pufferpanel /pufferpanel
EXPOSE 8080 5657
RUN mkdir -p /etc/pufferpanel && \
    mkdir -p /var/lib/pufferpanel

ARG DEBIAN_FRONTEND=noninteractive
ARG APTPROXY

# Because we use Alpine, we need musl to use the binary we generated
RUN if [ -n "$APTPROXY" ] ; then \
      echo "deb $APTPROXY/ubuntu-focal/ focal main restricted universe multiverse" > /etc/apt/sources.list ; \
      echo "deb $APTPROXY/ubuntu-focal-updates/ focal-updates main restricted universe multiverse" >> /etc/apt/sources.list ; \
      echo "deb $APTPROXY//ubuntu-focal-backports/ focal-backports main restricted universe multiverse" >> /etc/apt/sources.list ;\
    fi

RUN echo 'path-exclude /usr/share/doc/*' >/etc/dpkg/dpkg.cfg.d/docker-minimal && \
    echo 'path-exclude /usr/share/man/*' >>/etc/dpkg/dpkg.cfg.d/docker-minimal && \
    echo 'path-exclude /usr/share/groff/*' >>/etc/dpkg/dpkg.cfg.d/docker-minimal && \
    echo 'path-exclude /usr/share/info/*' >>/etc/dpkg/dpkg.cfg.d/docker-minimal && \
    echo 'path-exclude /usr/share/lintian/*' >>/etc/dpkg/dpkg.cfg.d/docker-minimal && \
    echo 'path-exclude /usr/share/linda/*' >>/etc/dpkg/dpkg.cfg.d/docker-minimal && \
    echo 'path-exclude /usr/share/locale/*' >>/etc/dpkg/dpkg.cfg.d/docker-minimal && \
    echo 'path-include /usr/share/locale/en*' >>/etc/dpkg/dpkg.cfg.d/docker-minimal

RUN apt-get update && \
    mkdir -p /usr/share/binfmts && \
    apt-get install -y git wget curl zip unzip musl

# java
RUN apt-get install --no-install-recommends -y openjdk-8-jdk-headless && \
    apt-get install --no-install-recommends -y openjdk-17-jdk-headless && \
    ln -sfn /usr/lib/jvm/java-8-openjdk-amd64/bin/java /usr/bin/java8 && \
    ln -sfn /usr/lib/jvm/java-8-openjdk-amd64/bin/javac /usr/bin/javac8 && \
    ln -sfn /usr/lib/jvm/java-17-openjdk-amd64/bin/java /usr/bin/java17 && \
    ln -sfn /usr/lib/jvm/java-17-openjdk-amd64/bin/javac /usr/bin/javac17 && \
    java -version && javac -version && \
    java8 -version && javac8 -version && \
    java17 -version && javac17 -version

# nodejs
RUN apt-get install -y nodejs

# srcds
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN echo steam steam/question select "I AGREE" | debconf-set-selections && \
    echo steam steam/license note '' | debconf-set-selections

ENV LANG 'en_US.UTF-8'
ENV LANGUAGE 'en_US:en'

RUN dpkg --add-architecture i386 && \
     apt-get update -y && \
     apt-get install -y --no-install-recommends ca-certificates locales steamcmd && \
     locale-gen en_US.UTF-8 && \
     ln -s /usr/games/steamcmd /usr/bin/steamcmd

# Cleanup
RUN apt-get clean all && rm -rf /var/lib/apt/lists/*

ENV PUFFER_LOGS=/etc/pufferpanel/logs \
    PUFFER_PANEL_TOKEN_PUBLIC=/etc/pufferpanel/public.pem \
    PUFFER_PANEL_TOKEN_PRIVATE=/etc/pufferpanel/private.pem \
    PUFFER_PANEL_DATABASE_DIALECT=sqlite3 \
    PUFFER_PANEL_DATABASE_URL="file:/etc/pufferpanel/pufferpanel.db?cache=shared" \
    PUFFER_DAEMON_SFTP_KEY=/etc/pufferpanel/sftp.key \
    PUFFER_DAEMON_DATA_CACHE=/var/lib/pufferpanel/cache \
    PUFFER_DAEMON_DATA_SERVERS=/var/lib/pufferpanel/servers \
    PUFFER_DAEMON_DATA_MODULES=/var/lib/pufferpanel/modules

WORKDIR /pufferpanel

ENTRYPOINT ["/pufferpanel/pufferpanel"]
CMD ["run"]
