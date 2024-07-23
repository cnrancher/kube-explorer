package content

import (
	"embed"
	"io/fs"
	"net/http"
	"path/filepath"
)

func NewEmbedded(staticContent embed.FS, prefix string) Handler {
	return newFS(&embedFS{
		pathPrefix:    prefix,
		staticContent: staticContent,
	})
}

var _ fsContent = &embedFS{}

type embedFS struct {
	pathPrefix    string
	staticContent embed.FS
}

// Open implements fsContent.
func (e *embedFS) Open(name string) (fs.File, error) {
	return e.staticContent.Open(joinEmbedFilepath(e.pathPrefix, name))
}

// ToFileServer implements fsContent.
func (e *embedFS) ToFileServer(basePaths ...string) http.Handler {
	handler := fsFunc(func(name string) (fs.File, error) {
		assetPath := joinEmbedFilepath(joinEmbedFilepath(basePaths...), name)
		return e.Open(assetPath)
	})

	return http.FileServer(http.FS(handler))
}

func (e *embedFS) Refresh() error { return nil }

func joinEmbedFilepath(paths ...string) string {
	return filepath.ToSlash(filepath.Join(paths...))
}
