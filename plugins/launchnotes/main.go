// Package main implements the LaunchNotes plugin for ReleasePilot.
package main

import (
	"github.com/felixgeelhaar/release-pilot/pkg/plugin"
)

func main() {
	plugin.Serve(&LaunchNotesPlugin{})
}
