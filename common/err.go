// This file is part of Darius. See License.txt for license information.

package common

import (
	"fmt"
	"log"
	"net/http"
)

// Fatalf is short-hand function to reply error message to user
func Fatalf(w http.ResponseWriter, err error, msg string, args ...interface{}) {
	Fatal(w, err, fmt.Sprintf(msg, args...))
}

// Fatal is short-hand function to reply error message to user
func Fatal(w http.ResponseWriter, err error, msg string) {
	w.WriteHeader(500)
	w.Write([]byte(msg))
	Error(err, msg)
}

// Errorf is short-hand function to reply error message to user
func Errorf(err error, msg string, args ...interface{}) {
	Error(err, fmt.Sprintf(msg, args...))
}

// Error is short-hand function to reply error message to user
func Error(err error, msg string) {
	log.Printf("%s: %s", msg, err)
}
