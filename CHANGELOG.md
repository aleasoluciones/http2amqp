## [1.0.0] - 2021-12-17
* Migrate to Go modules, work with Go 1.17
* Add Jenkins pipeline
* Update dependencies
* Add docker-compose for running a RabbitMQ in dev
* Add environment variables
* Clean Makefile
* Improve README

## [0.4.8] - 2019-02-22
* Port is parameterized

## [0.4.7] - 2018-10-29
* Internal: Update go version 1.10 and 1.11
* Internal: Update godep version
* Internal: Update Makefile and dependencies
* Internal: Update go version for docker image

## [0.4.6] - 2018-10-25
* Systems: Use network host mode

## [0.4.5] - 2018-01-30
* Internal: Added godep for dependency management
            - Updated Makefile
            - Excluded vendor directory from git
            - Updated .travis
            - Updated dependencies

## [0.4.4] - 2017-12-15
* Internal: Removed go version 1.6 and 1.7. Added 1.8.x, 1.9.x and master versions

## [0.4.3] - 2016-10-24

* Use manual prefix "docker/" for docker_log_tag.It's necessary for docker v1.12
* Update go version to 1.6.3
* Travis configuration uses 1.6 1.7 and tip
* golint is not supported in go version 1.5
* Remove update_deps doesnt work correctly :-(

## [0.4.2] - 2016-04-20

* Update deprecated docker compose configuration syslog-tag

## [0.4.2] - 2016-04-19

* Internal: Unify Makefile and travis.yml

## [0.4.2] - 2016-01-11

* Publish queries with the timeout value as amqp message ttl so messages are expired after timeout

## [0.4.1] - 2015-12-30

* Bugfix: ignore messages that can not be deserialized.

## [0.4.0] - 2015-12-02

* return 408 error code when timeout happens

## [0.4.0] - 2015-08-11

* Support a query string argument to especify a request timeout in milliseconds (e.g. curl http://localhost/foo?timeout=1500)

## [0.3.0] - 2015-06-17

* Support multiple methods
* simplified code
* API changes

## [0.2.0] - 2015-02-24

* Bugfix. Memory leak in the map that stores pending responses.
* Improved logging to show query id, topic and criteria.

## [0.1.1] - 2015-02-23

* Nothing. Just released to test deployment

## [0.1.0] - 2015-02-23

* Initial release
