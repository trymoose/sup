package ui

import (
	"bytes"
	"errors"
	"github.com/trymoose/sup/pkg/args"
	"gopkg.in/yaml.v3"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type (
	Destinations map[string]Dest
	Dest         interface {
		Remove(string) error
		Create(string) (io.WriteCloser, error)
	}
	FSDest struct {
		Path string `yaml:"path"`
	}
)

func (d *FSDest) Remove(s string) error {
	return os.Remove(filepath.Join(d.Path, s))
}

func (d *FSDest) Create(s string) (io.WriteCloser, error) {
	return os.Create(filepath.Join(d.Path, s))
}

var destinations = map[string]func() Dest{
	"path": func() Dest { return new(FSDest) },
}

func ReadDestinations(args *args.Args) Destinations {
	d := Destinations{}
	if f, err := os.Open(string(args.DestinationsFile)); err == nil {
		defer f.Close()
		if yaml.NewDecoder(f).Decode(&d); err != nil {
			if e := new(UnmarshalError); errors.As(err, &e) {
				e.Log(slog.Default())
			} else {
				slog.Error("failed to decode destinations", "error", err)
			}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		slog.Error("failed to open destinations", "error", err)
	}
	return d
}

type UnmarshalError struct {
	Message string
	Name    string
	Type    *string
	Base    error
}

func (u *UnmarshalError) Error() string {
	var buf bytes.Buffer
	u.Log(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelError,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})))
	return strings.TrimSpace(buf.String())
}

func (u *UnmarshalError) Log(l *slog.Logger) {
	a := []any{"name", u.Name}
	if u.Type != nil {
		a = append(a, "type", *u.Type)
	}
	if u.Base != nil {
		a = append(a, "error", u.Base)
	}
	l.Error(u.Message, a...)
}

func (d *Destinations) UnmarshalYAML(value *yaml.Node) error {
	*d = Destinations{}
	var m map[string]yaml.Node
	if err := value.Decode(&m); err != nil {
		return err
	}

	for name, config := range m {
		var t struct {
			Type string `yaml:"type"`
		}
		if err := config.Decode(&t); err != nil {
			return &UnmarshalError{Message: "failed to get destination type", Name: name, Base: err}
		}

		cons, ok := destinations[t.Type]
		if !ok {
			return &UnmarshalError{Message: "type does not exist", Name: name, Type: &t.Type}
		}

		v := cons()
		if err := config.Decode(v); err != nil {
			return &UnmarshalError{Message: "failed to decode destination config", Name: name, Type: &t.Type, Base: err}
		}
		(*d)[name] = v
	}
	return nil
}
