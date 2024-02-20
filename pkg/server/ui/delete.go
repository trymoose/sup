package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/trymoose/sup/pkg/args"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func DeleteFiles(args *args.Args, w http.ResponseWriter, r *http.Request) {
	var fns []string
	if err := json.NewDecoder(r.Body).Decode(&fns); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if len(fns) == 0 {
		http.Error(w, "no files to delete", http.StatusBadRequest)
		return
	}

	var err error
	var deleted []string
	for _, fn := range fns {
		if fn == "" {
			err = errors.Join(err, errors.New("filename cannot be empty"))
		} else if e := os.Remove(filepath.Join(string(args.Data), fn)); err != nil {
			err = errors.Join(e, fmt.Errorf("delete %q: %w", fn, err))
		} else {
			deleted = append(deleted, fn)
		}
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		slog.Info("deleted files", "files", deleted)
	}
}
