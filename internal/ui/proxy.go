package ui

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"

	"github.com/rancher/apiserver/pkg/urlbuilder"
	"k8s.io/apimachinery/pkg/util/proxy"
)

type RoundTripFunc func(*http.Request) (*http.Response, error)

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}

func proxyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme := urlbuilder.GetScheme(r)
		host := urlbuilder.GetHost(r, scheme)
		pathPrepend := r.Header.Get(urlbuilder.PrefixHeader)

		if scheme == r.URL.Scheme && host == r.URL.Host && pathPrepend == "" {
			next.ServeHTTP(w, r)
			return
		}

		proxyRoundtrip := proxy.Transport{
			Scheme:      scheme,
			Host:        host,
			PathPrepend: pathPrepend,
			RoundTripper: RoundTripFunc(func(r *http.Request) (*http.Response, error) {
				rw := &dummyResponseWriter{
					next:   w,
					header: make(http.Header),
				}
				next.ServeHTTP(rw, r)
				return rw.getResponse(r), nil
			}),
		}
		//proxyRoundtripper will write the response in RoundTrip func
		resp, _ := proxyRoundtrip.RoundTrip(r)
		responseToWriter(resp, w)
	})

}

var _ http.ResponseWriter = &dummyResponseWriter{}
var _ http.Hijacker = &dummyResponseWriter{}

type dummyResponseWriter struct {
	next http.ResponseWriter

	header     http.Header
	body       bytes.Buffer
	statusCode int
}

// Hijack implements http.Hijacker.
func (drw *dummyResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := drw.next.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("")
}

// Header implements the http.ResponseWriter interface.
func (drw *dummyResponseWriter) Header() http.Header {
	return drw.header
}

// Write implements the http.ResponseWriter interface.
func (drw *dummyResponseWriter) Write(b []byte) (int, error) {
	return drw.body.Write(b)
}

// WriteHeader implements the http.ResponseWriter interface.
func (drw *dummyResponseWriter) WriteHeader(statusCode int) {
	drw.statusCode = statusCode
}

// GetStatusCode returns the status code written to the response.
func (drw *dummyResponseWriter) GetStatusCode() int {
	if drw.statusCode == 0 {
		return 200
	}
	return drw.statusCode
}

func (drw *dummyResponseWriter) getResponse(req *http.Request) *http.Response {
	return &http.Response{
		Status:     strconv.Itoa(drw.GetStatusCode()),
		StatusCode: drw.GetStatusCode(),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Request:    req,
		Header:     drw.header,
		Body:       io.NopCloser(&drw.body),
	}
}

func responseToWriter(resp *http.Response, writer http.ResponseWriter) {
	for k, v := range resp.Header {
		writer.Header()[k] = v
	}
	if resp.StatusCode != http.StatusOK {
		writer.WriteHeader(resp.StatusCode)
	}
	_, _ = io.Copy(writer, resp.Body)
}
