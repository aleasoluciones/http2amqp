#!/bin/bash

env | grep -v '^PATH' >> /etc/environment

case "$1" in
    "integration-tests")
        source dev/env_develop
        exec make test
    ;;

esac

exec "$@"