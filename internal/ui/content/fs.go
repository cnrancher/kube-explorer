package content

import (
	"io"
	"net/http"
	"path/filepath"
	"sync"
)

var _ Handler = &handler{}

func newFS(content fsContent) Handler {
	return &handler{
		content: content,
		cacheFS: &sync.Map{},
	}
}

type handler struct {
	content fsContent
	cacheFS *sync.Map
}

func (h *handler) pathExist(path string) bool {
	_, err := h.content.Open(path)
	return err == nil
}

func (h *handler) serveContent(basePaths ...string) http.Handler {
	key := filepath.Join(basePaths...)
	if rtn, ok := h.cacheFS.Load(key); ok {
		return rtn.(http.Handler)
	}

	rtn := h.content.ToFileServer(basePaths...)
	h.cacheFS.Store(key, rtn)
	return rtn
}

func (h *handler) Refresh() {
	h.cacheFS.Range(func(key, _ any) bool {
		h.cacheFS.Delete(key)
		return true
	})
}

func (h *handler) ServeAssets(middleware func(http.Handler) http.Handler, next http.Handler) http.Handler {
	assets := middleware(h.serveContent())
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.pathExist(r.URL.Path) {
			assets.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (h *handler) ServeFaviconDashboard() http.Handler {
	return h.serveContent("dashboard")

}

func (h *handler) GetIndex() ([]byte, error) {
	path := filepath.Join("dashboard", "index.html")
	f, err := h.content.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
