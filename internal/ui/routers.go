package ui

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func New(opt *Options) (http.Handler, APIUI) {
	vue := NewUIHandler(opt)
	router := mux.NewRouter()
	router.UseEncodedPath()

	router.Handle("/", http.RedirectHandler("/dashboard/", http.StatusFound))
	router.Handle("/dashboard", http.RedirectHandler("/dashboard/", http.StatusFound))
	router.Handle("/dashboard/", vue.IndexFile())
	router.Handle("/favicon.png", vue.ServeFaviconDashboard())
	router.Handle("/favicon.ico", vue.ServeFaviconDashboard())
	router.PathPrefix("/dashboard/").Handler(vue.ServeAssets(vue.IndexFile()))
	router.PathPrefix("/api-ui/").Handler(vue.ServeAssets(http.NotFoundHandler()))
	router.PathPrefix("/k8s/clusters/local").HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		url := strings.TrimPrefix(req.URL.Path, "/k8s/clusters/local")
		if url == "" {
			url = "/"
		}
		http.Redirect(rw, req, url, http.StatusFound)
	})

	return router, apiUI(opt)
}
