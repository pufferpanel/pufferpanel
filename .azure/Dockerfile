FROM alpine

COPY pufferpanel /pufferpanel/
COPY client/dist /pufferpanel/client/dist/
RUN echo '{}' > /pufferpanel/config.json

ENV PATH "/pufferpanel:${PATH}"
EXPOSE 8080

CMD ["pufferpanel"]