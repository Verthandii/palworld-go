//go:build windows

package config

import "path"

var (
	gameConfigFile = path.Join("Pal", "Saved", "Config", "WindowsServer", "PalWorldSettings.ini")
)
