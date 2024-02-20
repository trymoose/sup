//go:build wireinject

package main

import (
	"context"
	"errors"
	"github.com/google/wire"
	"github.com/trymoose/sup/pkg/args"
	"github.com/trymoose/sup/pkg/server"
	"github.com/trymoose/sup/pkg/server/ui"
	"github.com/trymoose/sup/pkg/server/upload"
)

func NewServers(ctx context.Context) (*Servers, func(), error) {
	panic(wire.Build(args.NewArgs, ProvideUploadServer, ProvideUIServer, ProvideServers))
}

func ProvideUploadServer(ctx context.Context, args *args.Args) (*upload.Server, func(), error) {
	panic(wire.Build(server.NewServer, upload.NewUploadServer))
}

func ProvideUIServer(ctx context.Context, args *args.Args) (*ui.Server, func(), error) {
	panic(wire.Build(server.NewServer, ui.NewUIServer))
}

type Servers struct {
	UploadServer *upload.Server
	UIServer     *ui.Server
}

func (srv *Servers) Wait() error {
	return errors.Join(srv.UploadServer.Wait(), srv.UIServer.Wait())
}

func ProvideServers(us *upload.Server, cs *ui.Server) *Servers {
	panic(wire.Build(wire.Struct(new(Servers), "*")))
}
