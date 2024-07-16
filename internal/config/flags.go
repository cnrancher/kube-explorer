package config

import (
	"github.com/urfave/cli"
)

var InsecureSkipTLSVerify bool
var SystemDefaultRegistry string

var ShellPodImage string

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
	}
}
