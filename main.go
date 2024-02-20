package main

import (
	"context"
	dbg "github.com/trymoose/debug"
	"log/slog"
	"os"
	"runtime/debug"
)

//go:generate go run github.com/google/wire/cmd/wire@latest

func main() {
	exit := 1
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic", "recovered", r)
			if dbg.Debug {
				debug.PrintStack()
			}
		}
		os.Exit(exit)
	}()

	srv, cleanup, err := NewServers(context.Background())
	if err != nil {
		slog.Error("failed to init", "error", err)
		return
	}
	defer cleanup()

	if err := srv.Wait(); err != nil {
		slog.Error("server stopped", "error", err)
		return
	}
	exit = 0
}
