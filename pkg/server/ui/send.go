package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/trymoose/sup/pkg/args"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type SendRequest struct {
	Files        []string `json:"files"`
	Destinations []string `json:"destinations"`
}

func SendFiles(args *args.Args, w http.ResponseWriter, r *http.Request) {
	var req SendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Files) == 0 {
		http.Error(w, "no files to send", http.StatusBadRequest)
		return
	} else if len(req.Destinations) == 0 {
		http.Error(w, "no destinations given", http.StatusBadRequest)
		return
	}

	dest := ReadDestinations(args)
	if len(dest) == 0 {
		http.Error(w, "no destinations exist", http.StatusInternalServerError)
		return
	}

	dests := make([]Dest, len(req.Destinations))
	for i, d := range req.Destinations {
		v, ok := dest[d]
		if !ok {
			http.Error(w, fmt.Sprintf("destination %q does not exist", d), http.StatusBadRequest)
		}
		dests[i] = v
	}

	for _, fn := range req.Files {
		if err := sendFile(args, fn, dests); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func sendFile(args *args.Args, fn string, dests []Dest) (err error) {
	out := make([]io.Writer, 0, len(dests))
	cl := make([]io.Closer, 0, len(dests))
	defer func() {
		if err != nil {
			for _, cl := range cl {
				err = errors.Join(err, cl.Close())
			}
			for i := range out {
				err = errors.Join(err, dests[i].Remove(fn))
			}
		} else {
			err = os.Remove(filepath.Join(filepath.Join(string(args.Data), fn)))
		}
	}()

	f, err := os.Open(filepath.Join(string(args.Data), fn))
	if err != nil {
		return err
	}
	cl = append(cl, f)

	for _, d := range dests {
		o, err := d.Create(fn)
		if err != nil {
			return err
		}
		out = append(out, o)
		cl = append(cl, o)
	}

	if _, err = io.Copy(io.MultiWriter(out...), f); err != nil {
		return err
	}
	return nil
}
