FROM alpine
ADD config.yaml /config.yaml
ADD job-web /job-web
WORKDIR /
ENTRYPOINT [ "/job-web" ]
