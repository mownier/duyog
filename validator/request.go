package validator

import (
	"github.com/mownier/duyog/progerr"
	"net/http"
	"strings"
)

// Request interface
type Request interface {
	Validate(*http.Request, ...string) error
}

// ValidateRequest method
func ValidateRequest(r Request, req *http.Request, methods ...string) error {
	return r.Validate(req, methods...)
}

type serviceRequest struct{}

func (serviceRequest) Validate(r *http.Request, methods ...string) error {
	if len(methods) == 0 {
		return nil
	}

	ok := false
	method := strings.ToLower(r.Method)

	for _, m := range methods {
		if strings.ToLower(m) == method {
			ok = true
			break
		}
	}

	if ok == false {
		return progerr.MethodNotAllowed
	}

	return nil
}

type jsonRequest struct {
	serviceRequest
}

func (j jsonRequest) Validate(r *http.Request, methods ...string) error {
	if err := ValidateRequest(j.serviceRequest, r, methods...); err != nil {
		return err
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return progerr.ContentTypeNotJSON
	}

	return nil
}

// ServiceRequest method
func ServiceRequest() Request {
	return serviceRequest{}
}

// JSONRequest method
func JSONRequest() Request {
	return jsonRequest{
		serviceRequest: ServiceRequest().(serviceRequest),
	}
}
