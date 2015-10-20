FROM scratch


COPY http2amqp /

EXPOSE 18080
ENTRYPOINT ["/http2amqp"]
