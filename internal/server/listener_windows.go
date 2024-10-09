//go:build windows
// +build windows

package server

import (
	"fmt"
	"net"
	"net/url"

	"github.com/Microsoft/go-winio"
	"github.com/cnrancher/kube-explorer/internal/config"
)

func ensureListener() (net.Listener, string, error) {
	if config.BindAddress == "" {
		return nil, "", nil
	}
	u, err := url.Parse(config.BindAddress)
	if err != nil {
		return nil, "", err
	}
	switch u.Scheme {
	case "":
		return nil, config.BindAddress, nil
	case "tcp":
		return nil, u.Host, nil
	case "namedpipe":
		listener, err := winio.ListenPipe(u.Path, nil)
		return listener, u.Path, err
	default:
		return nil, "", fmt.Errorf("Unsupported scheme %s, only tcp and namedpipe are supported in windows", u.Scheme)
	}
}
