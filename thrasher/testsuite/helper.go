package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	//"reflect"
	"time"
)

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

func withErrorHanlder(fn apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if e := fn(w, r); e != nil { // e is *apiError, not os.Error
			if err := json.NewEncoder(w).Encode(e); err != nil {
				log.Println("[testsuite withErrorHanlder] Failed to encode JSON data")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

func withAuthentication(w http.ResponseWriter, r *http.Request, user string) *apiError {
	// Authenticate user
	jwtString := r.Header.Get("Authorization")
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
		return key, nil
	})

	if err != nil || token.Valid == false {
		str := ""
		if err != nil {
			str = fmt.Sprintf("Failed to authenticate user %s", err)
		} else {
			str = fmt.Sprintf("Failed to authenticate user %s", user)
		}
		log.Println("[testsuite.withAuthentication] ", str)
		w.WriteHeader(http.StatusUnauthorized)
		return &apiError{Error: err, Ret: -1, Info: str}
	}

	return nil
}
