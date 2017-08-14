package writer

import (
	"github.com/mownier/duyog/logger"
	"net/http"
)

// Response interface
type Response interface {
	Write(w http.ResponseWriter, r *http.Request, statusCode int, data []byte)
}

// WriteResponse method
func WriteResponse(resp Response, w http.ResponseWriter, r *http.Request, statusCode int, data []byte) {
	resp.Write(w, r, statusCode, data)
}

type serviceResponse struct {
	log logger.Response
}

func (s serviceResponse) Write(w http.ResponseWriter, r *http.Request, statusCode int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)

	logger.LogResponse(s.log, r, statusCode, data)
}

// ServiceResponse method
func ServiceResponse(l logger.Response) Response {
	return serviceResponse{
		log: l,
	}
}
