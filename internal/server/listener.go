package server

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/cnrancher/kube-explorer/internal/config"
	dynamicserver "github.com/rancher/dynamiclistener/server"
	"github.com/rancher/steve/pkg/server"
	"github.com/sirupsen/logrus"
)

func Serve(ctx context.Context, server *server.Server) error {
	listener, ipOrPath, err := ensureListener(ctx)
	if err != nil {
		return err
	}
	if listener != nil {
		defer listener.Close()
		return serveSocket(ctx, ipOrPath, listener, server)
	}
	return server.ListenAndServe(ctx, config.Steve.HTTPSListenPort, config.Steve.HTTPListenPort, &dynamicserver.ListenOpts{
		BindHost: ipOrPath,
	})
}

func serveSocket(ctx context.Context, socketPath string, listener net.Listener, handler http.Handler) error {
	logger := logrus.StandardLogger()
	errorLog := log.New(logger.WriterLevel(logrus.DebugLevel), "", log.LstdFlags)
	socketServer := &http.Server{
		Handler:  handler,
		ErrorLog: errorLog,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
	go func() {
		logrus.Infof("Listening on %s", socketPath)
		err := socketServer.Serve(listener)
		if err != http.ErrServerClosed && err != nil {
			logrus.Fatalf("https server failed: %v", err)
		}
	}()
	go func() {
		<-ctx.Done()
		_ = socketServer.Shutdown(context.Background())
	}()
	<-ctx.Done()
	return ctx.Err()
}
