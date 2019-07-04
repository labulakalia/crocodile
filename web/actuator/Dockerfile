FROM alpine
ADD config.yaml /config.yaml
ADD actuator-web /actuator-web
WORKDIR /
ENTRYPOINT [ "/actuator-web" ]
