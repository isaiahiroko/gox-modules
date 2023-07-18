package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var (
	ErrMethodNotAllowed = errors.New("method is not allowed for this router")
	ErrBadPath          = errors.New("every path definition must conform to [Method]:[Url]")
	ErrNilHandler       = errors.New("nill handler provided")
)

type Endpoint struct {
	Method  string
	Path    string
	Handler httprouter.Handle
}

type Server struct{}

func (s *Server) Start(addr string, routes []Endpoint) error {
	router := httprouter.New()

	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.Handler)
	}

	fmt.Printf("Server running on: %s", addr)

	err := http.ListenAndServe(addr, router)

	return err
}

func New() *Server {
	return &Server{}
}
