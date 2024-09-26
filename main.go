package main

import (
	"os"

	"github.com/cnrancher/kube-explorer/internal/version"
	"github.com/rancher/steve/pkg/debug"
	stevecli "github.com/rancher/steve/pkg/server/cli"
	"github.com/rancher/wrangler/v3/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	keconfig "github.com/cnrancher/kube-explorer/internal/config"
	"github.com/cnrancher/kube-explorer/internal/server"
)

func main() {
	app := cli.NewApp()
	app.Name = "kube-explorer"
	app.Version = version.FriendlyVersion()
	app.Usage = ""
	app.Flags = joinFlags(
		stevecli.Flags(&keconfig.Steve),
		debug.Flags(&keconfig.Debug),
		keconfig.Flags(),
	)
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(_ *cli.Context) error {
	ctx := signals.SetupSignalContext()
	keconfig.Debug.MustSetupDebug()
	s, err := server.ToServer(ctx, &keconfig.Steve, false)
	if err != nil {
		return err
	}
	return s.ListenAndServe(ctx, keconfig.Steve.HTTPSListenPort, keconfig.Steve.HTTPListenPort, nil)
}

func joinFlags(flags ...[]cli.Flag) []cli.Flag {
	var rtn []cli.Flag
	for _, flag := range flags {
		rtn = append(rtn, flag...)
	}
	return rtn
}
