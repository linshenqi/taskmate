FROM alpine:3.6

RUN apk update && apk add curl
RUN mkdir -p /etc/taskmate/conf /etc/taskmate/log

COPY ./entrypoint.sh /
COPY ./conf/config.yml /etc/taskmate/config.yml
COPY ./build/taskmate /usr/bin

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD taskmate --config /etc/taskmate/conf/config.yml
