Readme for workload module
=============================

Data Model
-----------------
   * https://github.com/thinkhy/thrasher/wiki/Data-Model
        
Golang Packages
-----------------
   * [go-json-rest](https://github.com/ant0ine/go-json-rest/rest)
   * [mgo](src/labix.org/v2/): Rich MongoDB driver for Go
   * [testify](github.com/stretchr/testify/assert)
   * [govalidator](github.com/asaskevich/govalidator)
   * local packages:
     - mkdir -p $GOPATH/src/github.com/thinkhy/thrasher
     - ln -s thrasher/workload/      $GOPATH/src/github.com/thinkhy/thrasher/workload/
     - ln -s thrasher/workload/mgodb $GOPATH/src/github.com/thinkhy/thrasher/workload/mgodb
   * [ginkgo](https://github.com/onsi/ginkgo)
     - go get github.com/onsi/ginkgo/ginkgo
     - go get github.com/onsi/gomega 


Command
---------------
   * Check status of RESTAPI server: `curl -i -u admin:@_@admin@_@ localhost:8002/.status`

Existing issues
-----------------
   * None





