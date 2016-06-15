package main

import (
	// "fmt"
	"./mgodb"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/thinkhy/logops"
	"log"
	"net/http"
	"time"
)

const (
	TESTAPIKEY = "PLSHAVEATRY"
)

var opshdlr *logops.Hook // handler for sending operation logs to InfluxDB

var Config struct {
	LoginUrl          string
	LogoutUrl         string
	LoinUrl           string
	JoinUrl           string
	UserUrl           string
	ChangePasswordUrl string
	Addr              string
	InfluxDBAddr      string
	DbAddress         string
	DbName            string
	Expiration        time.Duration
	Secret            string
}

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
	Config.LoginUrl = "/login"
	Config.LogoutUrl = "/logout"
	Config.ChangePasswordUrl = "/user/password"
	Config.JoinUrl = "/join"
	Config.UserUrl = "/user"
	Config.Addr = ":8000"
	Config.DbName = "testhub"
	Config.DbAddress = "127.0.0.1:27017"
	Config.InfluxDBAddr = "45.55.21.6:8089"
	Config.Expiration = time.Second * 3600 * 24 * 7
	Config.Secret = "testhub2016"
}

func handlers() *mux.Router {
	var router = mux.NewRouter()

	router.HandleFunc(Config.LoginUrl,
		withCheckProcessTime(Config.LoginUrl, withExtraHeader(withErrorHanlder(handleLogin)))).Methods("POST")
	router.HandleFunc(Config.JoinUrl,
		withCheckProcessTime(Config.JoinUrl, withExtraHeader(withErrorHanlder(handleJoin)))).Methods("POST")
	router.HandleFunc(Config.LogoutUrl,
		withCheckProcessTime(Config.LogoutUrl, withExtraHeader(withErrorHanlder(handleLogout)))).Methods("POST") // logout does authentication in handleLogout
	router.HandleFunc(Config.UserUrl,
		withCheckProcessTime("delete user", withExtraHeader(withErrorHanlder(handleDeleteUser)))).Methods("DELETE")
	// withCheckProcessTime("get userinfo", withErrorHanlder(handleGetUser))).Methods("GET")
	router.HandleFunc(Config.UserUrl,
		withCheckProcessTime("get userinfo", withExtraHeader(withErrorHanlder(handleGetUser)))).Methods("GET")

	router.HandleFunc(Config.UserUrl,
		withCheckProcessTime("update userinfo", withExtraHeader(withErrorHanlder(handleUpdateUser)))).Methods("PUT")
	router.HandleFunc(Config.ChangePasswordUrl,
		withCheckProcessTime("change password", withExtraHeader(withErrorHanlder(handleChangePassword)))).Methods("POST") // ChangePassword does authentication becuase password already be provided

	return router
}

func withErrorHanlder(fn apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if e := fn(w, r); e != nil { // e is *apiError, not os.Error
			if err := json.NewEncoder(w).Encode(e); err != nil {
				log.Println("[user withErrorHanlder] Failed to encode JSON data")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isValidAPIKey(r.URL.Query().Get("key")) {
			respondErr(w, r, http.StatusUnauthorized, "invalid API key")
			return
		}
		fn(w, r)
	}
}

func isValidAPIKey(key string) bool {
	return key == TESTAPIKEY
}

func checkProcessTime(tag string, start time.Time) {
	duration := float64(time.Since(start).Nanoseconds()) / float64(time.Millisecond)
	log.Printf("[%s] request-processing time is %vms\n", tag, duration)
}

func withCheckProcessTime(tag string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer checkProcessTime(tag, time.Now())
		fn(w, r)
	}
}

func withExtraHeader(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// with Content Type
		w.Header().Set("Content-Type", "application/json")

		// with CORS
		// Refer to Blog: http://www.2cto.com/Article/201509/441863.html
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		fn(w, r)
	}
}

func withAuthentication(w http.ResponseWriter, r *http.Request, user string) (*apiError, string) {
	// Authenticate user
	jwtString := r.Header.Get("Authorization")
	var uuid string
	key := []byte(Config.Secret)
	log.Println("[user.withAuthentication] Verifying JWT token ", jwtString)
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what we expect
		// Comparing agaist *jwt.SigningMethodHSA if using an HS256 token
		// Refer to: https://github.com/dgrijalva/jwt-go/issues/123
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return key, errors.New("Unexpected signing method")
		}

		if token.Claims["user"].(string) != user {
			return key, errors.New("User doesn't match")
		}
		uuid = token.Claims["id"].(string)
		return key, nil
	})

	if err != nil || token.Valid == false {
		str := ""
		if err != nil {
			str = fmt.Sprintf("Failed to authenticate user: %s", err)
		} else {
			str = fmt.Sprintf("Failed to authenticate user: %s", token.Claims["user"])
		}
		log.Println("[user.withAuthentication] ", str)
		w.WriteHeader(http.StatusUnauthorized)
		return &apiError{Error: err, Ret: -1, Info: str}, ""
	}

	return nil, uuid
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
	setupInfluxDB()
	mgodb.Setup(Config.DbAddress)

	mux := handlers()
	// http.Handle("/", r)

	log.Println("Starting web server on", Config.Addr)
	err := http.ListenAndServe(Config.Addr, mux)
	if err != nil {
		fmt.Println("Failed to start web server: ", err)
	}
	log.Println("Stopping web server ...")
}
