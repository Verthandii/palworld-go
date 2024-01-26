//go:build windows

package config

import (
	"path/filepath"
)

var (
	gameConfigFile = filepath.Join("Pal", "Saved", "Config", "WindowsServer", "PalWorldSettings.ini")
)
