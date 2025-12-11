// Package main implements the Linux package repository plugin for ReleasePilot.
package main

import (
	"github.com/felixgeelhaar/release-pilot/pkg/plugin"
)

func main() {
	plugin.Serve(&LinuxPkgPlugin{})
}
