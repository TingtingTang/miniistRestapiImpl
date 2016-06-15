package main
import (
        "net/http"
	"log"
)
func main() {
    addr := ":8085"
    http.Handle("/", http.FileServer(http.Dir(".")))
    log.Println("Start file server on ", addr)
    err := http.ListenAndServe(addr, nil)
    if nil != err {
        panic(err)
    }
}
