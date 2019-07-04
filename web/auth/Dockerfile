FROM alpine
ADD ../../conf/config.yaml /config.yaml
ADD auth-web /auth-web
WORKDIR /
ENTRYPOINT [ "/auth-web" ]
