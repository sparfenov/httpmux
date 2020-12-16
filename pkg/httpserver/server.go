package httpserver

import (
	"net"
	"net/http"
)

type Server struct {
	*http.Server

	// Maximum simultaneous requests.
	// set above zero to limit simultaneous requests
	LimitReq uint
}

// Used for middlewares
func (srv *Server) Use(handlerFunc func(http.Handler) http.Handler) {
	srv.Handler = handlerFunc(srv.Handler)
}

// Overrides standard ListenAndServe with option to limit simultaneous requests
func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	if srv.LimitReq > 0 {
		ln = NewLimitedListener(ln, srv.LimitReq)
	}

	return srv.Serve(ln)
}
