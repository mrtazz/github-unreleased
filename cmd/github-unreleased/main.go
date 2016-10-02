package main

import (
	"fmt"
	"os"

	"github.com/mrtazz/github-unreleased/config"
	"github.com/mrtazz/github-unreleased/logger"
	"github.com/mrtazz/github-unreleased/unreleased"

	docopt "github.com/docopt/docopt-go"
	"github.com/olekukonko/tablewriter"
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

	ur, err := unreleased.NewUnreleasedWithConfig(cfg)

	if args["<repository>"] != nil {
		reposlug := args["<repository>"].(string)
		commits, err := ur.GetUnreleasedForRepo(reposlug)
		if err != nil {
			fmt.Println(err.Error())
		}
		// we store all the commits with their 4 fields SHA, Message, Name, URL
		data := make([][]string, 0, len(commits))

		for _, commit := range commits {
			data = append(data, []string{
				commit.Sha[0:10], commit.Author, commit.Message, commit.URL})
		}
		fmt.Printf("\n ==> There are %d commits since the last release for %q\n",
			len(commits), reposlug)
		printTable([]string{"SHA", "Author", "Message", "URL"}, data)
	}

}

func printTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)

	if len(header) > 0 {
		table.SetHeader(header)
	}

	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetAutoWrapText(true)

	table.AppendBulk(data)
	table.Render()
}
