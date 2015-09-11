package common

import (
	"fmt"
	"log"
	"net/http"
)

func Fatalf(w http.ResponseWriter, err error, msg string, args ...interface{}) {
	Fatal(w, err, fmt.Sprintf(msg, args...))
}

func Fatal(w http.ResponseWriter, err error, msg string) {
	w.WriteHeader(500)
	w.Write([]byte(msg))
	Error(err, msg)
}

func Errorf(err error, msg string, args ...interface{}) {
	Error(err, fmt.Sprintf(msg, args...))
}

func Error(err error, msg string) {
	log.Printf("%s: %s", msg, err)
}
