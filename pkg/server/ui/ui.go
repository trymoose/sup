package ui

import (
	"bytes"
	_ "embed"
	"github.com/trymoose/sup/pkg/args"
	"github.com/trymoose/sup/pkg/server/ui/files"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func GetUI(args *args.Args, w http.ResponseWriter, r *http.Request) {
	switch {
	case serveFile(args, w, r):
	case serveStaticFile(w, r):
	default:
		serveIndex(args, w)
	}
}

func serveIndex(args *args.Args, w http.ResponseWriter) {
	if err := execTemplate(args, w); err != nil {
		slog.Error("failed to serve index", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func execTemplate(args *args.Args, w http.ResponseWriter) error {
	dot, err := getDot(args)
	if err != nil {
		return err
	}

	tmpl, err := getIndexTemplate()
	if err != nil {
		return err
	}

	return tmpl.Execute(w, dot)
}

func getIndexTemplate() (*template.Template, error) {
	idx := files.GetIndex()
	defer idx.Close()
	b, err := io.ReadAll(idx)
	if err != nil {
		return nil, err
	}
	return template.New("").Parse(string(b))
}

func getDot(args *args.Args) (dot files.IndexDot, err error) {
	var fns []os.DirEntry
	fns, err = os.ReadDir(string(args.Data))

	for _, fi := range fns {
		if !fi.IsDir() {
			dot.Files = append(dot.Files, fi.Name())
		}
	}

	for name := range ReadDestinations(args) {
		dot.Destinations = append(dot.Destinations, name)
	}

	dot.Sort()
	return
}

func serveStaticFile(w http.ResponseWriter, r *http.Request) bool {
	_, fn := path.Split(r.URL.Path)
	if d, ok := files.StaticFiles[fn]; ok {
		slog.Info("serving file", "filename", fn)
		r = r.Clone(r.Context())
		r.URL.RawQuery = ""
		http.ServeFileFS(w, r, &singleFileFS{name: fn, r: bytes.NewReader(d.Data)}, d.Name)
		return true
	}
	return false
}

func serveFile(args *args.Args, w http.ResponseWriter, r *http.Request) bool {
	dir, fn := path.Split(r.URL.Path)
	if strings.Contains(dir, "files") {
		slog.Info("serving file", "filename", fn)
		http.ServeFileFS(w, r, &singleFileFS{name: fn, path: string(args.Data)}, fn)
		return true
	}
	return false
}

type singleFileFS struct {
	name string
	r    *bytes.Reader
	path string
}

func (s *singleFileFS) Open(name string) (fs.File, error) {
	if name == s.name {
		if s.r != nil {
			return s, nil
		}
		return os.Open(filepath.Join(s.path, name))
	}
	return nil, fs.ErrNotExist
}

func (s *singleFileFS) Stat() (fs.FileInfo, error)                { return s, nil }
func (s *singleFileFS) Read(b []byte) (int, error)                { return s.r.Read(b) }
func (s *singleFileFS) Seek(off int64, whence int) (int64, error) { return s.r.Seek(off, whence) }
func (s *singleFileFS) Close() error                              { return nil }
func (s *singleFileFS) Name() string                              { return s.name }
func (s *singleFileFS) Size() int64                               { return s.r.Size() }
func (s *singleFileFS) Mode() fs.FileMode                         { return os.ModePerm }

var modTime = time.Now()

func (s *singleFileFS) ModTime() time.Time { return modTime }
func (s *singleFileFS) IsDir() bool        { return false }
func (s *singleFileFS) Sys() any           { return nil }
