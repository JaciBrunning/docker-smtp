FROM alpine

RUN apk add --update postfix opendkim

COPY ./entrypoint.sh /

ENTRYPOINT ["sh", "/entrypoint.sh"]
CMD postfix start; opendkim -f