package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"image"
	"image/png"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"
)

func main() {
	slog.Info("downloading svg")
	fail := true
	defer catch(&fail)
	slog.Info("downloading favicon", "url", os.Args[1])
	icon := downloadIcon()
	slog.Info("saving favicon")
	data := saveFavicon(icon)
	slog.Info("generating device pngs")
	svg := tryValue(oksvg.ReadReplacingCurrentColor(data, "#ffffff"))
	writePNGs(svg)
	fail = false
}

func writePNGs(svg *oksvg.SvgIcon) {
	try(os.WriteFile("apple-touch-icon.png", toPNG(svg, 180, 180), os.ModePerm))
	try(os.WriteFile("google-touch-icon.png", toPNG(svg, 512, 512), os.ModePerm))
}

func downloadIcon() io.ReadCloser {
	resp := tryValue(http.Get(os.Args[1]))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		try(httpError(resp.StatusCode))
	}
	return resp.Body
}

func saveFavicon(icon io.ReadCloser) io.Reader {
	defer icon.Close()
	var buf bytes.Buffer
	tryValue(io.Copy(&buf, icon))
	try(os.WriteFile("favicon.svg", bytes.Replace(buf.Bytes(), []byte(`fill="currentColor"`), []byte(`fill="#`+os.Args[2]+`"`), 1), os.ModePerm))
	return &buf
}

func toPNG(icon *oksvg.SvgIcon, w int, h int) []byte {
	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)
	var buf bytes.Buffer
	try(png.Encode(&buf, rgba))
	return buf.Bytes()
}

func catch(fail *bool) {
	if r := recover(); r != nil {
		if e := new(throw); errors.As(r.(error), &e) {
			slog.Error("throw", "error", e.err, "location", fmt.Sprintf("{ %s:%d }", e.fn, e.ln))
		} else {
			slog.Error("panic", "recover", r)
		}
		os.Exit(1)
	} else if *fail {
		os.Exit(1)
	}
}

func try(err error) {
	_try(err)
}

func tryValue[T any](t T, err error) T {
	_try(err)
	return t
}

func _try(err error) {
	if err != nil {
		t := &throw{err: err}
		_, t.fn, t.ln, _ = runtime.Caller(2)
		panic(t)
	}
}

type throw struct {
	ln  int
	fn  string
	err error
}

func (tr *throw) Unwrap() error { return tr.err }

func (tr *throw) Error() string {
	return tr.err.Error()
}

type httpError int

func (he httpError) Error() string {
	return fmt.Sprintf("got http %03d: %q", he, http.StatusText(int(he)))
}
