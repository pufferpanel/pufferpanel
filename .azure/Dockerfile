FROM alpine

COPY pufferpanel /pufferpanel/
COPY client/dist /pufferpanel/client/dist/
RUN echo '{}' > /pufferpanel/config.json

WORKDIR "/pufferpanel"
EXPOSE 8080
CMD ["/pufferpanel/pufferpanel run"]