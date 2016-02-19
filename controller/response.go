package controller

import (
	"encoding/json"
	"fmt"
	"log"
)

// Response represents common json format of controller's output
type Response struct {
	Result  bool        `json:"result"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Errf is short-hand function to reply error code to user
func (r *Response) Errf(code int, format string, args ...interface{}) *Response {
	return r.Err(code, fmt.Sprintf(format, args...))
}

// Err is short-hand function to reply error code to user
func (r *Response) Err(code int, msg string) *Response {
	r.Result = false
	r.Message = msg
	data := make(map[string]int)
	data["code"] = code
	r.Data = data
	return r
}

// Failf is short-hand function to reply error code to user
func (r *Response) Failf(format string, args ...interface{}) *Response {
	return r.Fail(fmt.Sprintf(format, args...))
}

// Fail is short-hand function to reply error code to user
func (r *Response) Fail(msg string) *Response {
	r.Result = false
	r.Message = msg
	return r
}

// Ok is short-hand function to reply success message to user
func (r *Response) Ok(data interface{}) *Response {
	r.Result = true
	r.Data = data
	return r
}

// Do will send this response to user
func (r *Response) Do(w *json.Encoder) {
	if err := w.Encode(r); err != nil {
		log.Printf("Response: While encoding json data: %s", err)
		return
	}
}
