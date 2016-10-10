package unreleased

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mrtazz/github-unreleased/config"
)

func TestSimpleCheckCreation(t *testing.T) {

	cfg, _ := config.NewConfigFromFile("../fixtures/exampleConfig.ini")

	c, _ := NewCheckerWithConfig(cfg)

	assert.NotEqual(t, nil, c)
}
