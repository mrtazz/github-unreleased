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

type UnreleasedRepoCommits struct {
	RepoSlug string
	Tag      string
	Commits  []Commit
}

type UnreleasedCommits []UnreleasedRepoCommits

type UnreleasedChecker struct {
	cfg      *config.Config
	ghClient *gogithub.Client
}

// NewUnreleasedWithConfig returns an Unreleased type with the passed in
// configuration set
func NewUnreleasedWithConfig(cfg *config.Config) (ret *UnreleasedChecker, err error) {
	ret = &UnreleasedChecker{}
	// we ignore the error return here as GetConfigValue returns an empty string on err
	token, _ := cfg.GetConfigValue("token")
	ret.ghClient, err = buildGithubClient(token)

	return ret, err
}

// GetUnreleasedForRepo returns an UnreleasedCommits list
func (ur *UnreleasedChecker) GetUnreleasedForRepo(slug string) (ret UnreleasedCommits,
	err error) {

	ret = UnreleasedCommits{UnreleasedRepoCommits{}}

	ret[0].RepoSlug = slug
	slugParts := strings.Split(slug, "/")

	tags, _, err := ur.ghClient.Repositories.ListTags(slugParts[0],
		slugParts[1], nil)
	if err != nil {
		return ret, err
	}
	if len(tags) == 0 {
		return ret, nil
	}
	logger.Debug(fmt.Sprintf("Got tag %q for %q/%q", *tags[0].Name, slugParts[0], slugParts[1]))
	comparison, _, err := ur.ghClient.Repositories.CompareCommits(slugParts[0],
		slugParts[1], *tags[0].Name, "HEAD")
	if err != nil {
		return ret, err
	}

	ret[0].Tag = *tags[0].Name

	logger.Debug(fmt.Sprintf("Found %d unreleased commits", len(comparison.Commits)))

	for _, commit := range comparison.Commits {
		logger.Debug(fmt.Sprintf("Parsing %+v", commit))
		ret[0].Commits = append(ret[0].Commits, Commit{
			Sha:     *commit.SHA,
			Message: *commit.Commit.Message,
			Author:  *commit.Commit.Author.Name,
			URL:     *commit.HTMLURL})
	}

	return
}

func (ur *UnreleasedChecker) GetUnreleasedForCurrentUser(showForks bool) (ret UnreleasedCommits,
	err error) {

	opt := &gogithub.RepositoryListOptions{
		Affiliation: "owner",
		ListOptions: gogithub.ListOptions{PerPage: 50}}

	var allRepos []*gogithub.Repository
	for {
		repos, resp, err := ur.ghClient.Repositories.List("", opt)
		if err != nil {
			return ret, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	for _, repo := range allRepos {
		slug := fmt.Sprintf("%s/%s", *repo.Owner.Login, *repo.Name)
		if *repo.Fork == true && showForks == false {
			continue
		}
		repoCommits, err := ur.GetUnreleasedForRepo(slug)
		if err != nil {
			logger.Info(fmt.Sprintf("Unable to get data for %q: %s",
				slug, err.Error()))
		} else {
			ret = append(ret, repoCommits[0])
		}
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
