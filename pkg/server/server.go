package server

import (
	"context"
	"errors"
	"github.com/trymoose/sup/pkg/args"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
)

type Server struct {
	srv     *http.Server
	started atomic.Bool
	done    chan struct{}
	err     error
}

func NewServer(ctx context.Context, args *args.Args) (*Server, func()) {
	srv := &Server{
		srv: &http.Server{
			BaseContext: func(net.Listener) context.Context { return ctx },
			ConnContext: func(ctx context.Context, _ net.Conn) context.Context {
				ctx, cancel := context.WithTimeout(ctx, args.UploadTimeout)
				context.AfterFunc(ctx, cancel)
				return ctx
			},
		},
		done: make(chan struct{}),
	}

	return srv, func() {
		if srv.started.Load() {
			defer srv.srv.Close()
			ctx, cancel := context.WithTimeout(context.Background(), args.ShutdownTimeout)
			defer cancel()
			slog.Error("failed to shutdown", "error", srv.srv.Shutdown(ctx))
		} else {
			slog.Error("not started")
		}
	}
}

func (srv *Server) Start(addr string, mux http.Handler) error {
	if srv.started.CompareAndSwap(false, true) {
		srv.srv.Addr = addr
		srv.srv.Handler = mux
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		go func() {
			defer ln.Close()
			defer close(srv.done)
			slog.Info("starting http", "address", ln.Addr().String())
			srv.err = srv.srv.Serve(ln)
		}()
		return nil
	}
	return errors.New("already started")
}

func (srv *Server) Wait() error {
	<-srv.done
	return srv.err
}
