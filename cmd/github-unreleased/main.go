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
  --forks                 Also show unreleased commits in forks
  --no-tags               Also show repositories with no releases
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
	var unreleasedData unreleased.UnreleasedCommits

	if args["<repository>"] != nil {
		reposlug := args["<repository>"].(string)
		unreleasedData, err = ur.GetUnreleasedForRepo(reposlug)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		unreleasedData, err = ur.GetUnreleasedForCurrentUser(args["--forks"].(bool))
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	for _, unreleasedRepoCommits := range unreleasedData {
		// we store all the commits with their 4 fields SHA, Message, Name, URL
		data := make([][]string, 0, len(unreleasedRepoCommits.Commits))

		for _, commit := range unreleasedRepoCommits.Commits {
			data = append(data, []string{
				commit.Sha[0:10], commit.Author, commit.Message, commit.URL})
		}
		if len(unreleasedRepoCommits.Commits) > 0 {
			fmt.Printf("\n ==> There are %d commits since tag %q for %q\n",
				len(unreleasedRepoCommits.Commits), unreleasedRepoCommits.Tag,
				unreleasedRepoCommits.RepoSlug)
			printTable([]string{"SHA", "Author", "Message", "URL"}, data)
		} else if unreleasedRepoCommits.Tag == "" && args["--no-tags"].(bool) == true {
			fmt.Printf("\n ==> No tags for  %q\n", unreleasedRepoCommits.RepoSlug)
		}

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
