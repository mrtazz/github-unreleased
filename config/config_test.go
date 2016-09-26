// Package config testing
package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleConfig(t *testing.T) {

	cfg, err := NewConfigFromFile("../fixtures/exampleConfig.ini")

	require.Equal(t, nil, err, "Parsing of the fixture config file should have worked.")

	val, _ := cfg.GetConfigValue("bla")

	assert.Equal(t, "bla", val)
}

func TestFailConfig(t *testing.T) {

	_, err := NewConfigFromFile("../fixtures/idontexist.ini")

	assert.NotEqual(t, nil, err)
}
