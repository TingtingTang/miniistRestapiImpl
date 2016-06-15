package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	// "encoding/json"
	"log"
	"net/http"
)

// apiError struct containing an error and some other fields that will send to
// client when hit any error
// Refer to: http://blog.golang.org/error-handling-and-go
type apiError struct {
	Status int    `json:"-"`
	Ret    int    `json:"ret"`
	Info   string `json:"info"`
}

func WriteErrorInfo(w rest.ResponseWriter, e *apiError) {
	w.WriteHeader(e.Status)
	if e != nil { // e is *apiError, not os.Error
		if err := w.WriteJson(*e); err != nil {
			log.Println("[workload.WriteHeader] Failed to encode JSON data")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}



