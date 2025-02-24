package httpServer

import (
	"net/http"
)

type Controller interface {
	GetOrderedMiddlewares() []Middleware
	GetUrl() string
	Handle(http.ResponseWriter, *http.Request)
}
