package content

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"sync"
)

const (
	defaultIndex = "https://releases.rancher.com/dashboard/latest/index.html"
)

func NewExternal(getIndex func() string) Handler {
	return &externalIndexHandler{
		getIndexFunc: getIndex,
	}
}

var (
	insecureClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	_ Handler = &externalIndexHandler{}
)

type externalIndexHandler struct {
	sync.RWMutex
	getIndexFunc    func() string
	current         string
	downloadSuccess *bool
}

func (u *externalIndexHandler) ServeAssets(_ func(http.Handler) http.Handler, next http.Handler) http.Handler {
	return next
}

func (u *externalIndexHandler) ServeFaviconDashboard() http.Handler {
	return http.NotFoundHandler()
}

func (u *externalIndexHandler) GetIndex() ([]byte, error) {
	if u.canDownload() {
		var buffer bytes.Buffer
		if err := serveIndex(&buffer, u.current); err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	}
	return nil, errors.New("external index is not available")
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

func (u *externalIndexHandler) canDownload() bool {
	u.RLock()
	rtn := u.downloadSuccess
	u.RUnlock()
	if rtn != nil {
		return *rtn
	}

	return u.refresh()
}

func (u *externalIndexHandler) refresh() bool {
	u.Lock()
	defer u.RUnlock()

	u.current = u.getIndexFunc()
	if u.current == "" {
		u.current = defaultIndex
	}
	t := serveIndex(io.Discard, u.current) == nil
	u.downloadSuccess = &t
	return t
}

func (u *externalIndexHandler) Refresh() {
	_ = u.refresh()
}
