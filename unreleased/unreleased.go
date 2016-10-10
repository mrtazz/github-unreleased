// Package unreleased contains the logic for checking for unreleased changes
package unreleased

import (
	"fmt"
	"github.com/mrtazz/github-unreleased/config"
	"github.com/mrtazz/github-unreleased/logger"
)

// Checker is the client struct to init and run methods on to get
// unreleased commits
type Checker struct {
	cfg     *config.Config
	fetcher FetcherInterface
}

// NewCheckerWithConfig returns an Unreleased type with the passed in
// configuration set
func NewCheckerWithConfig(cfg *config.Config) (ret *Checker, err error) {
	ret = &Checker{}
	// we ignore the error return here as GetConfigValue returns an empty string on err
	token, _ := cfg.GetConfigValue("token")
	ret.fetcher = NewFetcher(token)

	return ret, err
}

// GetUnreleasedForRepo returns an UnreleasedCommits list
func (ur *Checker) GetUnreleasedForRepo(slug string) (ret RepoCommits, err error) {

	ret = RepoCommits{}

	ret.Slug = slug

	tags, err := ur.fetcher.FetchTags(slug)
	if err != nil {
		return ret, err
	}
	if len(tags) == 0 {
		return ret, nil
	}
	logger.Debug(fmt.Sprintf("Got tag %q for %q", tags[0].Name, slug))
	unreleasedCommits, err := ur.fetcher.CompareCommits(slug, tags[0].Name, "HEAD")
	if err != nil {
		return ret, err
	}

	ret.Tag = tags[0]

	logger.Debug(fmt.Sprintf("Found %d unreleased commits", len(unreleasedCommits)))

	for _, commit := range unreleasedCommits {
		logger.Debug(fmt.Sprintf("Parsing %+v", commit))
		ret.Commits = append(ret.Commits, Commit{
			Sha:     commit.Sha,
			Message: commit.Message,
			Author:  commit.Author,
			URL:     commit.URL})
	}

	return
}

// GetUnreleasedForCurrentUser returns unreleased commits from all repo for
// the current user, optionally including forks
func (ur *Checker) GetUnreleasedForCurrentUser(showForks bool) (ret Commits,
	err error) {

	allRepos, err := ur.fetcher.FetchRepos("owner", 50)
	if err != nil {
		logger.Info(fmt.Sprintf("Unable to get repos: %s", err.Error()))
		return
	}

	for _, repo := range allRepos {
		if repo.IsFork == true && showForks == false {
			continue
		}
		repoCommits, err := ur.GetUnreleasedForRepo(repo.Slug)
		if err != nil {
			logger.Info(fmt.Sprintf("Unable to get data for %q: %s",
				repo.Slug, err.Error()))
		} else {
			ret = append(ret, repoCommits)
		}
	}

	return
}
