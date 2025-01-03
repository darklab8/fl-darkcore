package web

/*
Entrypoint for front and for dev web server?
*/

import (
	"fmt"
	"log"
	"net/http"

	"github.com/darklab8/fl-darkcore/darkcore/builder"
	"github.com/darklab8/fl-darkcore/darkcore/web/registry"
)

type Web struct {
	filesystem *builder.Filesystem
	registry   *registry.Registion
	mux        *http.ServeMux

	site_root           string
	site_root_acceptors []string
}

type WebOpt func(w *Web)

func WithSiteRootAcceptors(acceptors []string, site_root string) WebOpt {
	return func(w *Web) {
		w.site_root_acceptors = acceptors
		w.site_root = site_root
	}
}

func NewWeb(filesystem *builder.Filesystem, opts ...WebOpt) *Web {
	w := &Web{
		filesystem: filesystem,
		registry:   registry.NewRegister(),
		mux:        http.NewServeMux(),
	}

	for _, opt := range opts {
		opt(w)
	}

	w.registry.Register(w.NewEndpointStatic())

	w.registry.Register(w.NewEndpointPing())

	return w
}

type WebServeOpts struct {
	Port *int
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		next.ServeHTTP(w, r)
	})
}

func (w *Web) Serve(opts WebServeOpts) {
	w.registry.Foreach(func(e *registry.Endpoint) {
		w.mux.HandleFunc(string(e.Url), e.Handler)
	})

	ip := "0.0.0.0"
	var port int = 8000
	if opts.Port != nil {
		port = *opts.Port
	}

	fmt.Printf("launching web server, visit http://localhost:%d to check it!\n", port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), CorsMiddleware(w.mux)); err != nil {
		log.Fatal(err)
	}
}
