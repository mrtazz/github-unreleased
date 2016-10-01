package main

import (
	"fmt"
	"os"

	"github.com/mrtazz/github-unreleased/config"
	"github.com/mrtazz/github-unreleased/logger"
	"github.com/mrtazz/github-unreleased/unreleased"

	docopt "github.com/docopt/docopt-go"
)

var (
	usage = `github-unreleased.

  Usage:
  github-unreleased [options] [<repository>]
  github-unreleased -h | --help
  github-unreleased --version

  Options:
  -h --help               Show this screen.
  --version               Show version.
  --debug                 Enable debug mode
`

	version   = ""
	buildTime = ""
	builder   = ""
	goversion = ""

	configPath = fmt.Sprintf("%s/.github-unreleased.ini", os.Getenv("HOME"))
)

func main() {

	args, err := docopt.Parse(usage, nil, true,
		fmt.Sprintf("%s %s built at %s by %s with %s",
			"github-unreleased", version, buildTime, builder, goversion), false)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	if args["--debug"] != false {
		logger.SetLogLevel(logger.DebugLevel)
	}

	var cfg *config.Config
	if _, err := os.Stat(configPath); err == nil {
		cfg, err = config.NewConfigFromFile(configPath)
		if err != nil {
			logger.Error(fmt.Sprintf("Unable to parse config file at '%q': %q",
				configPath, err.Error()))
		}
	}
	fmt.Println(cfg)

	if args["<repository>"] != false {
		fmt.Println(args)
	}

	os.Exit(1)
}
