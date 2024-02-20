package upload

import (
	"github.com/trymoose/sup/pkg/args"
	"github.com/trymoose/sup/pkg/server"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

type Server struct {
	*server.Server
}

func NewUploadServer(args *args.Args, srv *server.Server) (*Server, error) {
	if err := srv.Start(args.UploadAddress, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ct ctw
		defer func() { slog.Info("file upload", "path", r.URL.Path, "wrote", int(ct)) }()
		fn := time.Now().Format("2006-01-02T15:04:05.99") + path.Ext(r.URL.Path)

		f, err := os.Create(filepath.Join(string(args.Data), fn))
		if err != nil {
			slog.Error("failed to open output file", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()

		var body io.Reader = r.Body
		if args.LimitUpload > 0 {
			body = io.LimitReader(body, args.LimitUpload)
		}

		if _, err = io.Copy(io.MultiWriter(f, &ct), r.Body); err != nil {
			slog.Error("failed to write output file", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})); err != nil {
		return nil, err
	}
	return &Server{Server: srv}, nil
}

type ctw int

func (c *ctw) Write(b []byte) (int, error) { *c = (*c) + ctw(len(b)); return len(b), nil }
