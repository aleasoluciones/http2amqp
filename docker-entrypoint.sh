#!/bin/bash
set -e

# Inspired on https://github.com/docker-library/postgres script

case "$1" in
    "http2amqp")
        shift
        OPTIONS="$@"
        exec /go/bin/http2amqp ${OPTIONS}
    ;;
esac
exec "$@"
