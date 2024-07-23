package ui

import (
	"net/http"

	"github.com/cnrancher/kube-explorer/internal/ui/content"
	"github.com/rancher/apiserver/pkg/middleware"
	"github.com/sirupsen/logrus"
)

type StringSetting func() string
type BoolSetting func() bool

func StaticSetting[T any](input T) func() T {
	return func() T {
		return input
	}
}

type Handler struct {
	contentHandlers map[string]content.Handler
	pathSetting     func() string
	indexSetting    func() string
	releaseSetting  func() bool
	offlineSetting  func() string
	middleware      func(http.Handler) http.Handler
	indexMiddleware func(http.Handler) http.Handler
}

type Options struct {
	// The location on disk of the UI files
	Path StringSetting
	// The HTTP URL of the index file to download
	Index StringSetting
	// Whether or not to run the UI offline, should return true/false/dynamic/embed
	Offline StringSetting
	// Whether or not is it release, if true UI will run offline if set to dynamic
	ReleaseSetting BoolSetting
}

func NewUIHandler(opts *Options) *Handler {
	if opts == nil {
		opts = &Options{}
	}

	h := &Handler{
		contentHandlers: make(map[string]content.Handler),
		indexSetting:    opts.Index,
		offlineSetting:  opts.Offline,
		pathSetting:     opts.Path,
		releaseSetting:  opts.ReleaseSetting,
		middleware: middleware.Chain{
			middleware.Gzip,
			middleware.FrameOptions,
			middleware.CacheMiddleware("json", "js", "css"),
		}.Handler,
		indexMiddleware: middleware.Chain{
			middleware.Gzip,
			middleware.NoCache,
			middleware.FrameOptions,
			middleware.ContentType,
		}.Handler,
	}

	if h.indexSetting == nil {
		h.indexSetting = StaticSetting("")
	}

	if h.offlineSetting == nil {
		h.offlineSetting = StaticSetting("dynamic")
	}

	if h.pathSetting == nil {
		h.pathSetting = StaticSetting("")
	}

	if h.releaseSetting == nil {
		h.releaseSetting = StaticSetting(false)
	}

	h.contentHandlers["embed"] = content.NewEmbedded(staticContent, "ui")
	h.contentHandlers["false"] = content.NewExternal(h.indexSetting)
	h.contentHandlers["true"] = content.NewFilepath(h.pathSetting)

	return h
}

func (h *Handler) content() content.Handler {
	offline := h.offlineSetting()
	if handler, ok := h.contentHandlers[offline]; ok {
		return handler
	}
	embedHandler := h.contentHandlers["embed"]
	filepathHandler := h.contentHandlers["true"]
	externalHandler := h.contentHandlers["false"]
	// default to dynamic
	switch {
	case h.pathSetting() != "":
		if _, err := filepathHandler.GetIndex(); err == nil {
			return filepathHandler
		}
		fallthrough
	case h.releaseSetting():
		// release must use embed first
		return embedHandler
	default:
		// try embed
		if _, err := embedHandler.GetIndex(); err == nil {
			return embedHandler
		}
		return externalHandler
	}
}

func (h *Handler) ServeAssets(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.content().ServeAssets(h.middleware, next).ServeHTTP(w, r)
	})
}

func (h *Handler) ServeFaviconDashboard() http.Handler {
	return h.middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.content().ServeFaviconDashboard().ServeHTTP(w, r)
	}))
}

func (h *Handler) IndexFile() http.Handler {
	return h.indexMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rtn, err := h.content().GetIndex()
		if err != nil {
			logrus.Warnf("failed to serve index with error %v", err)
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(rtn)
	}))
}
