package logger

import (
	"log"
	"net/http"
)

// Request interface
type Request interface {
	Log(req *http.Request)
}

// Response interface
type Response interface {
	Log(req *http.Request, statusCode int, data []byte)
}

// LogRequest method
func LogRequest(r Request, req *http.Request) {
	r.Log(req)
}

// LogResponse method
func LogResponse(r Response, req *http.Request, statusCode int, data []byte) {
	r.Log(req, statusCode, data)
}

type requestLog struct{}
type responseLog struct{}

func (requestLog) Log(r *http.Request) {
	log.Println("["+r.Method+"]", r.URL.Path)
}

func (responseLog) Log(r *http.Request, statusCode int, data []byte) {
	log.Println("["+r.Method+"]", r.URL.Path, "\n", string(data))
}

// RequestLog method
func RequestLog() Request {
	return requestLog{}
}

// ResponseLog method
func ResponseLog() Response {
	return responseLog{}
}
