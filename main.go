package main

import (
	"os"

	"github.com/rancher/steve/pkg/debug"
	stevecli "github.com/rancher/steve/pkg/server/cli"
	"github.com/rancher/steve/pkg/version"
	"github.com/rancher/wrangler/v3/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/cnrancher/kube-explorer/internal/server"
)

var (
	config      stevecli.Config
	debugconfig debug.Config
)

func main() {
	app := cli.NewApp()
	app.Name = "kube-explorer"
	app.Version = version.FriendlyVersion()
	app.Usage = ""
	app.Flags = joinFlags(
		stevecli.Flags(&config),
		debug.Flags(&debugconfig),
		server.Flags(),
	)
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(_ *cli.Context) error {
	ctx := signals.SetupSignalContext()
	debugconfig.MustSetupDebug()
	s, err := server.ToServer(ctx, &config, false)
	if err != nil {
		return err
	}
	return s.ListenAndServe(ctx, config.HTTPSListenPort, config.HTTPListenPort, nil)
}

func joinFlags(flags ...[]cli.Flag) []cli.Flag {
	var rtn []cli.Flag
	for _, flag := range flags {
		rtn = append(rtn, flag...)
	}
	return rtn
}
