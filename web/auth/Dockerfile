FROM alpine
ADD config.yaml /config.yaml
ADD auth-web /auth-web
WORKDIR /
ENTRYPOINT [ "/auth-web" ]
