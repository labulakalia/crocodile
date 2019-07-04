FROM alpine
ADD ../../conf/config.yaml /config.yaml
ADD actuator-web /actuator-web
WORKDIR /
ENTRYPOINT [ "/actuator-web" ]
