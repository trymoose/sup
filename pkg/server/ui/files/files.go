package files

import (
	"bytes"
	_ "embed"
	"github.com/trymoose/debug"
	"io"
	"os"
	"path/filepath"
	"slices"
)

//go:generate wget -N https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/js/bootstrap.bundle.min.js
//go:generate wget -N https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css
//go:generate wget -N https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css
//go:generate wget -N https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/fonts/bootstrap-icons.woff
//go:generate wget -N https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/fonts/bootstrap-icons.woff2
//go:generate go run github.com/trymoose/sup/pkg/server/ui/files/internal/favicon https://icons.getbootstrap.com/assets/icons/archive-fill.svg BDC3C7

var (
	//go:embed favicon.svg
	Favicon []byte
	//go:embed index.gohtml
	index []byte
	//go:embed index.js
	IndexJS []byte
	//go:embed index.css
	IndexCSS []byte
	//go:embed bootstrap.bundle.min.js
	BootstrapJS []byte
	//go:embed bootstrap.min.css
	BootstrapCSS []byte
	//go:embed bootstrap-icons.min.css
	BootstrapIcons []byte
	//go:embed bootstrap-icons.woff
	BootstrapIconsWoff []byte
	//go:embed bootstrap-icons.woff2
	BootstrapIconsWoff2 []byte
	//go:embed manifest.json
	Manifest []byte
	//go:embed google-touch-icon.png
	GoogleTouchIcon []byte
	//go:embed apple-touch-icon.png
	AppleTouchIcon []byte
)

var StaticFiles = map[string]struct {
	Data []byte
	Name string
}{
	"favicon.svg":             {Name: "favicon.svg", Data: Favicon},
	"index.js":                {Name: "index.js", Data: IndexJS},
	"bootstrap.bundle.min.js": {Name: "bootstrap.bundle.min.js", Data: BootstrapJS},
	"bootstrap.min.css":       {Name: "bootstrap.min.css", Data: BootstrapCSS},
	"bootstrap-icons.min.css": {Name: "bootstrap-icons.min.css", Data: BootstrapIcons},
	"bootstrap-icons.woff":    {Name: "bootstrap-icons.woff", Data: BootstrapIconsWoff},
	"bootstrap-icons.woff2":   {Name: "bootstrap-icons.woff2", Data: BootstrapIconsWoff2},
	"index.css":               {Name: "index.css", Data: IndexCSS},
	"manifest.json":           {Name: "manifest.json", Data: Manifest},
	"google-touch-icon.png":   {Name: "google-touch-icon.png", Data: GoogleTouchIcon},
	"apple-touch-icon.png":    {Name: "apple-touch-icon.png", Data: AppleTouchIcon},
}

type IndexDot struct {
	Files        []string
	Destinations []string
}

func (d *IndexDot) Sort() {
	slices.Sort(d.Files)
	slices.Sort(d.Destinations)
}

func GetIndex() io.ReadCloser {
	return FileOrBytes("index.gohtml", index)
}

func FileOrBytes(name string, b []byte) io.ReadCloser {
	if debug.Debug {
		return GetFile(name)
	}
	return io.NopCloser(bytes.NewReader(b))
}

func StaticFilePath() []string {
	return []string{"pkg", "server", "ui", "files"}
}

func GetFile(name string) io.ReadCloser {
	f, err := os.Open(filepath.Join(append(StaticFilePath(), name)...))
	if err != nil {
		return newErrReader(err)
	}
	return f
}

type errReader struct{ err error }

func newErrReader(err error) io.ReadCloser   { return io.NopCloser(errReader{err}) }
func (e errReader) Read([]byte) (int, error) { return 0, e.err }
