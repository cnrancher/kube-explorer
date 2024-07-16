//go:build !embed

package ui

import (
	"net/http"
	"os"
	"path/filepath"
)

func (u *Handler) ServeAsset() http.Handler {
	return u.middleware(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		http.FileServer(http.Dir(u.pathSetting())).ServeHTTP(rw, req)
	}))
}

func (u *Handler) ServeFaviconDashboard() http.Handler {
	return u.middleware(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		http.FileServer(http.Dir(filepath.Join(u.pathSetting(), "dashboard"))).ServeHTTP(rw, req)
	}))
}

func (u *Handler) IndexFileOnNotFound() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// we ignore directories here because we want those to come from the CDN when running in that mode
		if stat, err := os.Stat(filepath.Join(u.pathSetting(), req.URL.Path)); err == nil && !stat.IsDir() {
			u.ServeAsset().ServeHTTP(rw, req)
		} else {
			u.IndexFile().ServeHTTP(rw, req)
		}
	})
}

func (u *Handler) IndexFile() http.Handler {
	return u.indexMiddleware(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if path, isURL := u.path(); isURL {
			_ = serveIndex(rw, path)
		} else {
			http.ServeFile(rw, req, filepath.Join(path, "index.html"))
		}
	}))
}
