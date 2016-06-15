Readme for testsuite module
=============================

Data Model
-----------------
   * https://github.com/thinkhy/thrasher/wiki/Data-Model

Golang Packages
-----------------
   * [mgo](src/labix.org/v2/): Rich MongoDB driver for Go
   * [testify](github.com/stretchr/testify/assert)
   * [govalidator](github.com/asaskevich/govalidator)
   * [jwt-go](github.com/dgrijalva/jwt-go)
   * local packages:
     - mkdir -p $GOPATH/src/github.com/thinkhy/thrasher/testsuite
     - ln -s thrasher/testsuite/      $GOPATH/src/github.com/thinkhy/thrasher/testsuite
     - ln -s thrasher/testsuite/mgodb $GOPATH/src/github.com/thinkhy/thrasher/testsuite/mgodb

Testing
----------------
   * long-run: go test -v
   * short-run: go test -v -short
   * cd mgo&&go test -v

Coverage:

```
    go test -coverprofile=coverage.out github.com/thinkhy/thrasher/testsuite
    go tool cover -html=coverage.out
```

Existing issues
-----------------
   * None
