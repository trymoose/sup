package ui

import (
	"github.com/trymoose/sup/pkg/args"
	"github.com/trymoose/sup/pkg/server"
	"net/http"
)

type Server struct {
	*server.Server
}

func NewUIServer(args *args.Args, srv *server.Server) (*Server, error) {
	if err := srv.Start(args.UIAddress, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetUI(args, w, r)
		case http.MethodPost:
			SendFiles(args, w, r)
		case http.MethodDelete:
			DeleteFiles(args, w, r)
		default:
			http.NotFound(w, r)
		}
	})); err != nil {
		return nil, err
	}
	return &Server{Server: srv}, nil
}
