ARG DOCKER_TAG=base-devel

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

# srcds
RUN apk add --no-cache bash
COPY --from=srcdsbuilder /lib/i386-linux-gnu /lib/
COPY --from=srcdsbuilder /root/installer/linux32/libstdc++.so.6 /lib/
COPY --from=srcdsbuilder /root/installer/steamcmd.sh /usr/lib/games/steam/
COPY --from=srcdsbuilder /root/installer/linux32/steamcmd /usr/lib/games/steam/
COPY --from=srcdsbuilder /usr/games/steamcmd /usr/bin/steamcmd

# Cleanup
RUN rm -rf /var/cache/apk/*