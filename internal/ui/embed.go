//go:build embed

package ui

import (
	"embed"
)

// content holds our static web server content.
//
//go:embed all:ui/*
var staticContent embed.FS
