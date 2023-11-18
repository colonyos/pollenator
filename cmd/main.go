package main

import (
	"github.com/colonyos/pollinator/internal/cli"
	"github.com/colonyos/pollinator/pkg/build"
)

var (
	BuildVersion string = ""
	BuildTime    string = ""
)

func main() {
	build.BuildVersion = BuildVersion
	build.BuildTime = BuildTime
	cli.Execute()
}
