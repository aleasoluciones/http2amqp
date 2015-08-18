FROM buildpack-deps:jessie-scm


# gcc for cgo
RUN apt-get update && apt-get install -y gcc libc6-dev make --no-install-recommends && rm -rf /var/lib/apt/lists/*

ENV GOLANG_GOOS linux
ENV GOLANG_GOARCH amd64
ENV GOLANG_VERSION 1.4

RUN curl -sSL https://golang.org/dl/go$GOLANG_VERSION.$GOLANG_GOOS-$GOLANG_GOARCH.tar.gz | tar -v -C /usr/local -xz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -pv $GOPATH/src/github.com/aleasoluciones

RUN mkdir $GOPATH/src/github.com/aleasoluciones/http2amqp
COPY . $GOPATH/src/github.com/aleasoluciones/http2amqp
RUN cd $GOPATH/src/github.com/aleasoluciones/http2amqp; make deps; make build
RUN mkdir -pv /go/bin/
RUN cp $GOPATH/src/github.com/aleasoluciones/http2amqp/http2amqp /http2amqp

COPY docker-entrypoint.sh /

EXPOSE 18080
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["http2amqp"]
