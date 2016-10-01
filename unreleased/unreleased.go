// Package unreleased contains the logic for checking for unreleased changes
package unreleased

import (
	gogithub "github.com/google/go-github/github"
	"github.com/mrtazz/github-unreleased/config"
	"golang.org/x/oauth2"
	"time"
)

type Commit struct {
	sha     string
	message string
	author  string
	date    time.Time
}

type UnreleasedCommits []Commit

type Unreleased struct {
	cfg      *config.Config
	ghClient *github.Client
}

// NewUnreleasedWithConfig returns an Unreleased type with the passed in
// configuration set
func NewUnreleasedWithConfig(cfg *config.Config) {
	ret := &Unreleased{}
	token, err := cfg.GetConfigValue("token")
	if err == nil {
		ret.ghClient = buildGithubClient(token)
	} else {
		ret.ghClient = buildGithubClient(nil)
	}

	return ret
}

// GetUnreleasedForRepo returns a list of
func (ur *Unreleased) GetUnreleasedForRepo(slug string) UnreleasedCommits {
}

// GetUnreleasedForRepo returns a list of
func (ur *Unreleased) GetUnreleasedForRepo(slug string) UnreleasedCommits {
}

func buildGithubClient(token string) (*gogithub.Client, error) {
	var tc oauth.Client
	if token != nil {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = oauth2.NewClient(oauth2.NoContext, ts)
	}

	return gogithub.NewClient(tc), nil
}
