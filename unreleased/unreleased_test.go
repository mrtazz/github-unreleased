package unreleased

import (
	"fmt"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	gogithub "github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"

	"github.com/mrtazz/github-unreleased/config"
)

func TestSimpleCheckCreation(t *testing.T) {

	cfg, _ := config.NewConfigFromFile("../fixtures/exampleConfig.ini")

	c, _ := NewCheckerWithConfig(cfg)

	assert.NotEqual(t, nil, c)
}

func TestMockedFetchRepos(t *testing.T) {
	// mock setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFetcher := NewMockFetcherInterface(mockCtrl)
	mockFetcher.EXPECT().FetchRepos("", "owner", 50).Return(
		Repositories{Repository{Slug: "mrtazz/test", IsFork: false}}, nil)
	mockFetcher.EXPECT().CompareCommits("mrtazz/test",
		"0.1.0", "HEAD").Return(
		[]Commit{{Sha: "123",
			Message: "test",
			Author:  "mrtazz",
			URL:     "https://github.com/mrtazz/test"}}, nil)
	mockFetcher.EXPECT().FetchTags("mrtazz/test").Return(
		Tags{
			{Name: "0.1.0", Date: time.Unix(1476152199, 0)},
			{Name: "0.0.9", Date: time.Unix(1476142199, 0)}}, nil)

	c := &Checker{fetcher: mockFetcher}
	ret, err := c.GetUnreleasedForCurrentUser(false)

	// assertions
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(ret))
}

func TestMockedFetchReposFailed(t *testing.T) {
	// mock setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFetcher := NewMockFetcherInterface(mockCtrl)
	mockFetcher.EXPECT().FetchRepos("", "owner", 50).Return(
		Repositories{}, fmt.Errorf("failed to get repo"))

	c := &Checker{fetcher: mockFetcher}
	ret, err := c.GetUnreleasedForCurrentUser(false)

	// assertions
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "failed to get repo", err.Error())
	assert.Equal(t, Commits{}, ret)
}

func TestMockedFetchReposIsFork(t *testing.T) {
	// mock setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFetcher := NewMockFetcherInterface(mockCtrl)
	mockFetcher.EXPECT().FetchRepos("", "owner", 50).Return(
		Repositories{Repository{Slug: "mrtazz/test", IsFork: true}}, nil)

	c := &Checker{fetcher: mockFetcher}
	ret, err := c.GetUnreleasedForCurrentUser(false)

	// assertions
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(ret))
}

func TestMockedUnreleasedForRepoFailed(t *testing.T) {
	// mock setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFetcher := NewMockFetcherInterface(mockCtrl)
	mockFetcher.EXPECT().FetchRepos("", "owner", 50).Return(
		Repositories{Repository{Slug: "mrtazz/test", IsFork: false}}, nil)

	mockFetcher.EXPECT().FetchTags("mrtazz/test").Return(
		Tags{}, fmt.Errorf("no tags"))

	c := &Checker{fetcher: mockFetcher}
	ret, err := c.GetUnreleasedForCurrentUser(false)

	// assertions
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(ret))
}

func TestMockedNoTags(t *testing.T) {
	// mock setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFetcher := NewMockFetcherInterface(mockCtrl)
	mockFetcher.EXPECT().FetchRepos("", "owner", 50).Return(
		Repositories{Repository{Slug: "mrtazz/test", IsFork: false}}, nil)

	mockFetcher.EXPECT().FetchTags("mrtazz/test").Return(
		Tags{}, nil)

	c := &Checker{fetcher: mockFetcher}
	ret, err := c.GetUnreleasedForCurrentUser(false)

	// assertions
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(ret))
}

func TestMockedCompareCommitsFailed(t *testing.T) {
	// mock setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFetcher := NewMockFetcherInterface(mockCtrl)
	mockFetcher.EXPECT().FetchRepos("", "owner", 50).Return(
		Repositories{Repository{Slug: "mrtazz/test", IsFork: false}}, nil)
	mockFetcher.EXPECT().CompareCommits("mrtazz/test",
		"0.1.0", "HEAD").Return(
		[]Commit{}, fmt.Errorf("compare failed"))
	mockFetcher.EXPECT().FetchTags("mrtazz/test").Return(
		Tags{{Name: "0.1.0", Date: time.Now()}}, nil)

	c := &Checker{fetcher: mockFetcher}
	ret, err := c.GetUnreleasedForCurrentUser(false)

	// assertions
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(ret))
}

// some integration tests which call to actual GitHub
func TestSimpleGitHubFetchRepos(t *testing.T) {

	f := &Fetcher{
		ghClient: gogithub.NewClient(nil)}

	repos, err := f.FetchRepos("mrtazz", "owner", 50)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, len(repos) > 1)
}

func TestSimpleGitHubCompareCommits(t *testing.T) {

	f := &Fetcher{
		ghClient: gogithub.NewClient(nil)}

	commits, err := f.CompareCommits("mrtazz/checkmake",
		"e27165fa24654c003d55ca85cc36195b768f1324",
		"55457d1dd6198f2c1f7c95151855accb39a50198")

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(commits))
}

func TestSimpleGitHubFetchTags(t *testing.T) {

	f := &Fetcher{
		ghClient: gogithub.NewClient(nil)}

	tags, err := f.FetchTags("mrtazz/restclient-cpp")

	assert.Equal(t, nil, err)
	assert.Equal(t, true, len(tags) > 1)
}
