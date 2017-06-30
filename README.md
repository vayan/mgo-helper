# mgo-helper

This lib regroups all the golang mongo driver mgo utils that we developped to integrate into our projects. 

[![CircleCI](https://circleci.com/gh/transcovo/mgo-helper.svg?style=shield&circle-token=3664fc4c8d1f7578b306d5aac5cc6bda59ac0eca)](https://circleci.com/gh/transcovo/mgo-helper)
[![codecov](https://codecov.io/gh/transcovo/mgo-helper/branch/master/graph/badge.svg)](https://codecov.io/gh/transcovo/mgo-helper)
[![GoDoc](https://godoc.org/github.com/transcovo/mgo-helper?status.svg)](https://godoc.org/github.com/transcovo/mgo-helper)

## Requirements

Minimum Go version: 1.7

## Installation

- If you are using govendor
    ```shell
    govendor fetch github.com/transcovo/mgo-helper
    ```

- standard way (not recommended)
    ```shell
    got get -u github.com/transcovo/mgo-helper
    ```

## Usage

The main function to be used to connect is `InitMongoFromConfig`
by providing a mongo.Configuration object

```go
mongoConfig := Configuration{
    PingFrequency: 100,
    SSLCert:       []byte{},
    UseSSL:        false,
    URL:           "mongodb://localhost:27017/some-test-db",
}
db, teardown := InitMongoFromConfig(mongoConfig)
defer teardown()
```

## Tests

Tests with go are a little weird: you can't run all the tests at once for your project.

Still, to do that, run `./tools/test.sh`. It will recursively run the tests for all packages. It's the fastest 
test script if you're following TDD. To add a coverage report generation to the tests (much slower), run `./tools/coverage.sh`.

You can run `./tools/coverage.sh --html` to run the tests with the coverage report and open the coverage
result in the browser.
