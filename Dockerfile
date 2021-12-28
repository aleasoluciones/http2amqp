FROM golang:1.17.5 AS http2amqp-builder

ENV DEBIAN_FRONTEND noninteractive

WORKDIR /app

COPY . .

# CGO_ENABLED needed for scratch image that serves the binary
# Ref: https://stackoverflow.com/questions/61515186/when-using-cgo-enabled-is-must-and-what-happens
ENV CGO_ENABLED 0

RUN make build

COPY docker-entrypoint.sh /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]


#---


FROM scratch AS http2amqp

COPY --from=http2amqp-builder /app/http2amqp /

EXPOSE 18080

ENTRYPOINT ["/http2amqp"]
