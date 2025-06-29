package middleware

import (
	"github.com/Hesam-Eskandari/gollum/library/httpServer"
	"net/http"
)

type HttpMethodValidator interface {
	httpServer.Middleware
	SetAllowedMethods(httpMethods []string)
}

func NewHttpMethodValidator() HttpMethodValidator {
	return &httpMethodValidatorImpl{
		methodSet: make(map[string]struct{}),
	}
}

type httpMethodValidatorImpl struct {
	methodSet map[string]struct{}
}

func (hmv *httpMethodValidatorImpl) SetAllowedMethods(httpMethods []string) {
	for _, method := range httpMethods {
		hmv.methodSet[method] = struct{}{}
	}
}

func (hmv *httpMethodValidatorImpl) ProcessBeforeHandler(w http.ResponseWriter, r *http.Request) bool {
	if _, ok := hmv.methodSet[r.Method]; !ok {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func (hmv *httpMethodValidatorImpl) ProcessAfterHandler(w http.ResponseWriter, r *http.Request) {
}
