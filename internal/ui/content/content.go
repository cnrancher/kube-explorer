package content

import (
	"io/fs"
	"net/http"
)

type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

type fsContent interface {
	ToFileServer(basePaths ...string) http.Handler
	Open(name string) (fs.File, error)
}

type Handler interface {
	ServeAssets(middleware func(http.Handler) http.Handler, hext http.Handler) http.Handler
	ServeFaviconDashboard() http.Handler
	GetIndex() ([]byte, error)
	Refresh()
}
