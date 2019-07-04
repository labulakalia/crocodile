FROM alpine
ADD ../../conf/config.yaml /config.yaml
ADD job-web /job-web
WORKDIR /
ENTRYPOINT [ "/job-web" ]
