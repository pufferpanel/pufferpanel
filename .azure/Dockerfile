FROM alpine

COPY . /pufferpanel/
RUN ls -l /pufferpanel/
RUN echo '{}' > /pufferpanel/config.json

ENV PATH "/pufferpanel:${PATH}"
EXPOSE 8080

CMD ["pufferpanel"]