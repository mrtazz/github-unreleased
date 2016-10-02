// Package unreleased contains the logic for checking for unreleased changes
package unreleased

import (
	"fmt"
	gogithub "github.com/google/go-github/github"
	"github.com/mrtazz/github-unreleased/config"
	"github.com/mrtazz/github-unreleased/logger"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

type Commit struct {
	Sha     string
	Message string
	Author  string
	URL     string
}

type UnreleasedCommits []Commit

type Unreleased struct {
	cfg      *config.Config
	ghClient *gogithub.Client
}

// NewUnreleasedWithConfig returns an Unreleased type with the passed in
// configuration set
func NewUnreleasedWithConfig(cfg *config.Config) (ret *Unreleased, err error) {
	ret = &Unreleased{}
	// we ignore the error return here as GetConfigValue returns an empty string on err
	token, _ := cfg.GetConfigValue("token")
	ret.ghClient, err = buildGithubClient(token)

	return ret, err
}

// GetUnreleasedForRepo returns an UnreleasedCommits list
func (ur *Unreleased) GetUnreleasedForRepo(slug string) (ret UnreleasedCommits,
	err error) {

	slugParts := strings.Split(slug, "/")

	tags, _, err := ur.ghClient.Repositories.ListTags(slugParts[0],
		slugParts[1], nil)
	if err != nil {
		return ret, err
	}
	if len(tags) == 0 {
		return ret, fmt.Errorf("no tags")
	}
	logger.Debug(fmt.Sprintf("Got tag %q for %q/%q", *tags[0].Name, slugParts[0], slugParts[1]))
	comparison, _, err := ur.ghClient.Repositories.CompareCommits(slugParts[0],
		slugParts[1], *tags[0].Name, "HEAD")
	if err != nil {
		return ret, err
	}

	logger.Debug(fmt.Sprintf("Found %d unreleased commits", len(comparison.Commits)))

	for _, commit := range comparison.Commits {
		logger.Debug(fmt.Sprintf("Parsing %+v", commit))
		ret = append(ret, Commit{
			Sha:     *commit.SHA,
			Message: *commit.Commit.Message,
			Author:  *commit.Commit.Author.Name,
			URL:     *commit.HTMLURL})
	}

	return
}

func buildGithubClient(token string) (*gogithub.Client, error) {
	var tc *http.Client
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = oauth2.NewClient(oauth2.NoContext, ts)
	}

	return gogithub.NewClient(tc), nil
}
