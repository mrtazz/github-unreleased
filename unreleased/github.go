// Package unreleased implements the unreleased checker logic and this file
// deals with the GitHub interface parts
//go:generate mockgen -source=github.go -package=unreleased -destination=github_mock.go
package unreleased

import (
	"fmt"
	gogithub "github.com/google/go-github/github"
	"github.com/mrtazz/github-unreleased/logger"
	"golang.org/x/oauth2"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Tag represents a git tag with the important information we need to identify
// and order tags
type Tag struct {
	Name string
	Date time.Time
}

// Tags is the collection type for Tag
type Tags []Tag

// implementation methods for Tag to provide the sort interface
func (t Tags) Len() int {
	return len(t)
}
func (t Tags) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t Tags) Less(i, j int) bool {
	return t[i].Date.Before(t[j].Date)
}

// Commit is the exported type to represent commit data
type Commit struct {
	Sha     string
	Message string
	Author  string
	URL     string
}

// RepoCommits is the struct to represent a repo with it's most
// recent tag and all unreleased commits
type RepoCommits struct {
	Slug    string
	Tag     Tag
	Commits []Commit
}

// Commits is a slice of UnreleasedRepoCommits
type Commits []RepoCommits

// Repository represents the minimal amount of information we care about for a
// GitHub repository
type Repository struct {
	Slug   string
	IsFork bool
}

// Repositories is a slice of Repository structs
type Repositories []Repository

// FetcherInterface defines the interface to satisfy for fetching information about
// repositories
type FetcherInterface interface {
	FetchTags(slug string) (Tags, error)
	FetchRepos(user string, affiliation string, perPage int) (Repositories,
		error)
	CompareCommits(slug string, base string, head string) ([]Commit, error)
}

// Fetcher implements FetcherInterface for GitHub
type Fetcher struct {
	ghClient *gogithub.Client
}

// NewFetcher creates a Fetcher for GitHub with the given token
func NewFetcher(token string) *Fetcher {
	return &Fetcher{ghClient: buildGithubClient(token)}
}

// FetchTags implements tag fetching for GitHub
func (f *Fetcher) FetchTags(slug string) (ret Tags, err error) {
	ret = make(Tags, 0, 10)
	slugParts := strings.Split(slug, "/")
	opt := &gogithub.ListOptions{PerPage: 50}
	for {
		tags, resp, err := f.ghClient.Repositories.ListTags(slugParts[0],
			slugParts[1], opt)
		if err != nil {
			return ret, err
		}
		for _, ghTag := range tags {

			commit, _, err := f.ghClient.Git.GetCommit(slugParts[0],
				slugParts[1], *ghTag.Commit.SHA)
			logger.Debug(fmt.Sprintf("Fetched %q for %q", commit, ghTag))
			if err != nil {
				return ret, err
			}
			ret = append(ret, Tag{Name: *ghTag.Name, Date: *commit.Author.Date})
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	sort.Sort(sort.Reverse(ret))
	return ret, err
}

// FetchRepos implements repo fetching for GitHub
func (f *Fetcher) FetchRepos(user string,
	affiliation string, perPage int) (ret Repositories, err error) {
	ret = Repositories{}
	opt := &gogithub.RepositoryListOptions{
		Affiliation: affiliation,
		ListOptions: gogithub.ListOptions{PerPage: perPage}}

	for {
		repos, resp, err := f.ghClient.Repositories.List(user, opt)
		if err != nil {
			return ret, err
		}
		for _, repo := range repos {
			ret = append(ret, Repository{
				Slug:   fmt.Sprintf("%s/%s", *repo.Owner.Login, *repo.Name),
				IsFork: *repo.Fork})

		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	return ret, err
}

// CompareCommits returns a slice of Commit structs from the diff between base and
// head
func (f *Fetcher) CompareCommits(slug string,
	base string, head string) (ret []Commit, err error) {

	ret = make([]Commit, 0)
	slugParts := strings.Split(slug, "/")
	comparison, _, err := f.ghClient.Repositories.CompareCommits(slugParts[0],
		slugParts[1], base, head)

	for _, commit := range comparison.Commits {
		logger.Debug(fmt.Sprintf("Parsing %+v", commit))
		ret = append(ret, Commit{
			Sha:     *commit.SHA,
			Message: *commit.Commit.Message,
			Author:  *commit.Commit.Author.Name,
			URL:     *commit.HTMLURL})
	}

	return ret, err
}

func buildGithubClient(token string) *gogithub.Client {
	var tc *http.Client
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = oauth2.NewClient(oauth2.NoContext, ts)
	}
	return gogithub.NewClient(tc)
}
