package core

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type Server interface {
	// Config returns server config
	Config() Config

	// Router returns server router
	Router() *mux.Router

	// ListenAndServe runs http server
	ListenAndServe() error
}

/*
New returns fresh gopypi server instance bound to config.
*/
func New(cfg Config) (result Server, err error) {

	var (
		chain  alice.Chain
	)

	if chain, err = InitRouter(cfg); err != nil {
		return
	}

	result = &server{
		config: cfg,
		chain:  chain,
		router: cfg.Router(),
	}

	return
}

/*
server implements Server interface
*/
type server struct {
	config Config
	chain  alice.Chain
	router *mux.Router
}

/*
Return server config
*/
func (s *server) Config() Config {
	return s.config
}

/*
Router returns instantiated mux router
*/
func (s *server) Router() *mux.Router {
	return s.router
}

/*
ListenAndServe starts http server and listens to requests
*/
func (s *server) ListenAndServe() (err error) {
	s.Config().Logger().Info("Attempting to listen on: 0.0.0.0:9900")

	final := s.chain.Then(s.Router())

	err = http.ListenAndServe("0.0.0.0:9900", final)
	return
}
