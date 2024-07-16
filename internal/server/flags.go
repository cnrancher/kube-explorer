package server

import (
	"github.com/urfave/cli"
)

var insecureSkipTLSVerify bool

func Flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:        "insecure-skip-tls-verify",
			Destination: &insecureSkipTLSVerify,
		},
	}
}
