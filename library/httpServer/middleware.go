package httpServer

import (
	"net/http"
)

// Middleware is usually singleton
type Middleware interface {
	// ProcessBeforeHandler runs before request is processed by the handler function
	ProcessBeforeHandler(w http.ResponseWriter, r *http.Request) (isSuccessful bool)
	// ProcessAfterHandler runs after request is processed by the handler function
	ProcessAfterHandler(w http.ResponseWriter, r *http.Request)
}
