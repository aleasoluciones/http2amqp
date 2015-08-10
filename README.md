# http2amqp

[![Build Status](https://travis-ci.org/aleasoluciones/http2amqp.svg)](https://travis-ci.org/aleasoluciones/http2amqp)


## Generating new version

Update code and commit changes.
Generate a new tag and push the tag. The version will be automatically upload to [github releases](https://github.com/aleasoluciones/http2amqp/releases)

Example:
```
git tag v0.3.0
git push
git push --tags
```

Will generate the 0.3.0 version at https://github.com/aleasoluciones/http2amqp/releases/download/v0.3.0/http2amqp

## TODO
 - test timeout parameter for each request
 - implement delay parameter for echo server to allow tests timeouts
 
