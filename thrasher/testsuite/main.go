package main

import (
	// "fmt"
	"./mgodb"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/thinkhy/logops"
	"log"
	"net/http"
	"time"
)

var Config struct {
	TestsuiteUrl  string
	TestsuitesUrl string
	Addr          string
	InfluxDBAddr  string
	DbAddress     string
	DbName        string
	Expiration    time.Duration
	Secret        string
	TestFlag      bool
}

var opshdlr *logops.Hook // handler for sending operation logs to InfluxDB

// apiError struct containing an error and some other fields that will send to
// client when hit any error
// Refer to: http://blog.golang.org/error-handling-and-go
type apiError struct {
	Error error  `json:"-"`
	Ret   int    `json:"ret"`
	Info  string `json:"info"`
}

type apiHandler func(http.ResponseWriter, *http.Request) *apiError

func loadConfig() {
	// TODO: replace below block with config file processing [ 2016-03-15 thinkhy ]
	Config.TestsuiteUrl = "/testsuite"
	Config.TestsuitesUrl = "/testsuites"
	Config.Addr = ":8001"
	Config.DbName = "testhub"
	Config.DbAddress = "127.0.0.1:27017"
	Config.InfluxDBAddr = "45.55.21.6:8089"
	Config.Expiration = time.Second * 3600 * 24 * 7
	Config.Secret = "testhub2016"
	Config.TestFlag = false
}

func handlers() *mux.Router {
	var router = mux.NewRouter()

	router.HandleFunc(Config.TestsuiteUrl,
		withCheckProcessTime("testsuite.handleInsert", withExtraHeader(withErrorHanlder(handleInsert)))).Methods("POST")

	router.HandleFunc(Config.TestsuiteUrl,
		withCheckProcessTime("testsuite.handleGet", withExtraHeader(withErrorHanlder(handleGet)))).Methods("GET")

	router.HandleFunc(Config.TestsuiteUrl,
		withCheckProcessTime("testsuite.handleUpdate", withExtraHeader(withErrorHanlder(handleUpdate)))).Methods("PUT")

	router.HandleFunc(Config.TestsuiteUrl,
		withCheckProcessTime("testsuite.handleDelete", withExtraHeader(withErrorHanlder(handleDelete)))).Methods("DELETE")

	router.HandleFunc(Config.TestsuitesUrl,
		withCheckProcessTime("testsuite.handleGetTestsuites", withExtraHeader(withErrorHanlder(handleGetTestsuites)))).Methods("GET")

	return router
}

func setupInfluxDB() {
	config := &logops.Config{
		Address: Config.InfluxDBAddr,
		UseUDP:  true,
	}
	var err error
	opshdlr, err = logops.NewHook(config)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func main() {
	loadConfig()
	mgodb.Setup(Config.DbAddress)
	setupInfluxDB()

	mux := handlers()

	log.Println("Starting web server on", Config.Addr)
	err := http.ListenAndServe(Config.Addr, mux)
	if err != nil {
		fmt.Println("Failed to start web server: ", err)
	}

	log.Println("Stopping web server ...")
}
