// Package main implements the Docker Hub / Container registry plugin for ReleasePilot.
package main

import (
	"github.com/felixgeelhaar/release-pilot/pkg/plugin"
)

func main() {
	plugin.Serve(&DockerPlugin{})
}
