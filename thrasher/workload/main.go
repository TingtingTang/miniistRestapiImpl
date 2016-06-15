package main

import (
	"./mgodb"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/thinkhy/logops"
	"log"
	"net/http"
	"time"
)

var opshdlr *logops.Hook // handler for sending operation logs to InfluxDB

var Config struct {
	WorkloadUrl  string
	WorkloadsUrl string
	Addr         string
	DbAddress    string
	InfluxDBAddr string
	DbName       string
	Expiration   time.Duration
	Secret       string
	TestFlag     bool
}

func loadConfig() {
	// TODO: replace below block with config file processing [ 2016-03-15 thinkhy ]
	Config.WorkloadUrl = "/workload"
	Config.WorkloadsUrl = "/workloads"
	Config.Addr = ":8002"
	Config.DbName = "testhub"
	Config.InfluxDBAddr = "45.55.21.6:8089"
	Config.DbAddress = "127.0.0.1:27017"
	Config.Expiration = time.Hour * 24 * 5 // [2016-05-09 sync expiration with front-end(JS) ]
	Config.Secret = "testhub2016"
	Config.TestFlag = false
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

func MakeApi() (api *rest.Api) {
	api = rest.NewApi()

	// TODO: Add JWT middleware

	// Add CORS support
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return true
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{
			"Accept",
			"Content-Type",
			"X-Custom-Header",
			"Origin",
			"Authorization",
		},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})

	// Add status auth
	statusMw := &rest.StatusMiddleware{}
	api.Use(statusMw)
	api.Use(rest.DefaultDevStack...)

	auth := &rest.AuthBasicMiddleware{
		Realm: "admin zone",
		Authenticator: func(userId string, password string) bool {
			// TODO: [2016-05-05] remove the smell hard-coding
			if userId == "admin" && password == "@_@admin@_@" {
				return true
			}
			return false
		},
	}
	router, err := rest.MakeRouter(
		// Setup a /.status endpoint protected with basic authentication
		rest.Get("/.status", auth.MiddlewareFunc(
			func(w rest.ResponseWriter, r *rest.Request) {
				w.WriteJson(statusMw.GetStatus())
			},
		)),

		rest.Get(Config.WorkloadUrl, GetWorkload),
		rest.Get(Config.WorkloadsUrl, GetWorkloads),
		rest.Post(Config.WorkloadUrl, InsertWorkload),
		rest.Put(Config.WorkloadUrl, UpdateWorkload),
		rest.Delete(Config.WorkloadUrl, DeleteWorkload),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	return
}

func main() {
	loadConfig()
	setupInfluxDB()
	mgodb.Setup(Config.DbAddress)

	api := MakeApi()

	log.Println("Starting web server on", Config.Addr)
	err := http.ListenAndServe(Config.Addr, api.MakeHandler())
	if err != nil {
		fmt.Println("Failed to start web server: ", err)
	}

	return
}
