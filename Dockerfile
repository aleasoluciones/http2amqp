FROM golang:1.17.2 AS http2amqp-builder
ENV DEBIAN_FRONTEND noninteractive

RUN mkdir -pv $GOPATH/github.com/aleasoluciones/http2amqp
WORKDIR $GOPATH/github.com/aleasoluciones/http2amqp
COPY . $GOPATH/github.com/aleasoluciones/http2amqp

# CGO_ENABLED needed for scratch image that serves the binary
# Ref: https://stackoverflow.com/questions/61515186/when-using-cgo-enabled-is-must-and-what-happens
ENV CGO_ENABLED 0
RUN make jenkins
RUN cp http2amqp /

FROM scratch AS http2amqp-compiled

COPY --from=build-stage http2amqp /

EXPOSE 18080
ENTRYPOINT ["/http2amqp", "-verbose"]
