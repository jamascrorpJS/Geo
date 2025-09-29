package internal

import (
	"context"
	"net/http"

	"jamascrorpJS/gwatch/pkg/config"
)

type Network interface {
	Start() error
	Shutdown(ctx context.Context) error
}

type network struct {
	s      *http.Server
	config config.Config
}

func New(h http.Handler, config config.Config) Network {
	port := config.GetString("server.port")
	return &network{
		s: &http.Server{
			Addr:    port,
			Handler: h,
		},
	}
}

func (n *network) Start() error {
	return n.s.ListenAndServe()
}

func (n *network) Shutdown(ctx context.Context) error {
	return n.s.Shutdown(ctx)
}
