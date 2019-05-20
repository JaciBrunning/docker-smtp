FROM golang:alpine

RUN apk add --update postfix opendkim git cyrus-sasl cyrus-sasl-plain

# RUN mkdir -p /tmp/etcd/dl && wget -nc -P /tmp/etcd/dl https://github.com/etcd-io/etcd/releases/download/v3.3.13/etcd-v3.3.13-linux-amd64.tar.gz
# RUN tar xzvf /tmp/etcd/dl/etcd-v3.3.13-linux-amd64.tar.gz && cp etcd-v3.3.13-linux-amd64/etcdctl /usr/bin && cp etcd-v3.3.13-linux-amd64/etcd /usr/bin
# RUN rm -r etcd-v3.3.13-linux-amd64 && rm -r /tmp/etcd

RUN go get go.etcd.io/etcd/clientv3

COPY src /app
RUN cd /app && go build

COPY ./entrypoint.sh /

ENV ETCDCTL_API 3
WORKDIR /app

ENTRYPOINT ["sh", "/entrypoint.sh"]
CMD ./app $DOMAIN $ACME_STORAGE/object $ETCD