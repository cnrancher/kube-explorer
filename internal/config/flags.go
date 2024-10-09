package config

import (
	"github.com/urfave/cli"
)

var InsecureSkipTLSVerify bool
var SystemDefaultRegistry string
var APIUIVersion = "1.1.11"
var ShellPodImage string
var BindAddress string

func Flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:        "insecure-skip-tls-verify",
			Destination: &InsecureSkipTLSVerify,
		},
		cli.StringFlag{
			Name:        "system-default-registry",
			Destination: &SystemDefaultRegistry,
		},
		cli.StringFlag{
			Name:        "pod-image",
			Destination: &ShellPodImage,
			Value:       "rancher/shell:v0.2.1-rc.7",
		},
		cli.StringFlag{
			Name:        "apiui-version",
			Hidden:      true,
			Destination: &APIUIVersion,
			Value:       APIUIVersion,
		},
		cli.StringFlag{
			Name:        "bind-address",
			Destination: &BindAddress,
			Usage:       `Bind address with url format. The supported schemes are unix, tcp and namedpipe, e.g. unix:///path/to/kube-explorer.sock or namedpipe:/\.\pipe\kube-explorer`,
		},
	}
}
