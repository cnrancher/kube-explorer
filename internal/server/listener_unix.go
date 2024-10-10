//go:build unix
// +build unix

package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/cnrancher/kube-explorer/internal/config"
	"github.com/sirupsen/logrus"
)

var _ net.Listener = &closerListener{}

type closerListener struct {
	listener net.Listener
	lockFile *os.File
}

func (l *closerListener) Accept() (net.Conn, error) {
	return l.listener.Accept()
}

func (l *closerListener) Close() error {
	return errors.Join(
		l.listener.Close(),
		l.lockFile.Close(),
		os.RemoveAll(l.lockFile.Name()),
	)
}

func (l *closerListener) Addr() net.Addr {
	return l.listener.Addr()
}

func ensureListener(ctx context.Context) (net.Listener, string, error) {
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
	case "unix":
		listener, err := createCloserListener(ctx, u.Path)
		if err != nil {
			return nil, "", err
		}
		return listener, u.Path, err
	default:
		return nil, "", fmt.Errorf("Unsupported scheme %s, only tcp and unix are supported in UNIX OS", u.Scheme)
	}
}

func createCloserListener(ctx context.Context, socketPath string) (net.Listener, error) {
	lockFilePath := getLockFileName(socketPath)
	lockFile, err := os.OpenFile(lockFilePath, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	lockErr := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if lockErr != nil {
		return nil, fmt.Errorf("Socket file %s is in use, exiting", socketPath)
	}

	if _, err := os.Stat(socketPath); err == nil {
		logrus.Infof("Removing stale socket file %s", socketPath)
		_ = os.Remove(socketPath)
	}

	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "unix", socketPath)
	if err != nil {
		return nil, err
	}

	return &closerListener{
		listener: listener,
		lockFile: lockFile,
	}, nil
}

func getLockFileName(socketPath string) string {
	return strings.TrimSuffix(socketPath, filepath.Ext(socketPath)) + ".lock"
}
