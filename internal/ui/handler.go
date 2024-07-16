package ui

import (
	"crypto/tls"
	"io"
	"net/http"
	"sync"

	"github.com/rancher/apiserver/pkg/middleware"
	"github.com/sirupsen/logrus"
)

const (
	defaultPath = "./ui"
)

var (
	insecureClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)

type StringSetting func() string
type BoolSetting func() bool

type Handler struct {
	pathSetting     func() string
	indexSetting    func() string
	releaseSetting  func() bool
	offlineSetting  func() string
	middleware      func(http.Handler) http.Handler
	indexMiddleware func(http.Handler) http.Handler

	downloadOnce    sync.Once
	downloadSuccess bool
}

type Options struct {
	// The location on disk of the UI files
	Path StringSetting
	// The HTTP URL of the index file to download
	Index StringSetting
	// Whether or not to run the UI offline, should return true/false/dynamic
	Offline StringSetting
	// Whether or not is it release, if true UI will run offline if set to dynamic
	ReleaseSetting BoolSetting
}

func NewUIHandler(opts *Options) *Handler {
	if opts == nil {
		opts = &Options{}
	}

	h := &Handler{
		indexSetting:   opts.Index,
		offlineSetting: opts.Offline,
		pathSetting:    opts.Path,
		releaseSetting: opts.ReleaseSetting,
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
		h.indexSetting = func() string {
			return "https://releases.rancher.com/dashboard/latest/index.html"
		}
	}

	if h.offlineSetting == nil {
		h.offlineSetting = func() string {
			return "dynamic"
		}
	}

	if h.pathSetting == nil {
		h.pathSetting = func() string {
			return defaultPath
		}
	}

	if h.releaseSetting == nil {
		h.releaseSetting = func() bool {
			return false
		}
	}

	return h
}

func (u *Handler) path() (path string, isURL bool) {
	switch u.offlineSetting() {
	case "dynamic":
		if u.releaseSetting() {
			return u.pathSetting(), false
		}
		if u.canDownload(u.indexSetting()) {
			return u.indexSetting(), true
		}
		return u.pathSetting(), false
	case "true":
		return u.pathSetting(), false
	default:
		return u.indexSetting(), true
	}
}

func (u *Handler) canDownload(url string) bool {
	u.downloadOnce.Do(func() {
		if err := serveIndex(io.Discard, url); err == nil {
			u.downloadSuccess = true
		} else {
			logrus.Errorf("Failed to download %s, falling back to packaged UI", url)
		}
	})
	return u.downloadSuccess
}

func serveIndex(resp io.Writer, url string) error {
	r, err := insecureClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	_, err = io.Copy(resp, r.Body)
	return err
}
