package httpServer

import (
	"net/http"
)

type Controller interface {
	Handle(http.ResponseWriter, *http.Request)
	GetOrderedMiddlewares() []Middleware
	GetUrl() string
}
