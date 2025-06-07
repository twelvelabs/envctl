package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentService_Get(t *testing.T) {
	config, err := NewTestConfig()
	assert.NoError(t, err)
	config.Environments = []Environment{
		{
			Name: "common",
			Values: map[string]string{
				"ONE": "common-one",
				"TWO": "common-two",
			},
		},
		{
			Name: "local",
			Extends: []string{
				"common",
			},
			Values: map[string]string{
				"THREE": "local-three",
			},
		},
		{
			Name: "staging",
			Extends: []string{
				"common",
			},
			Values: map[string]string{
				"TWO":   "staging-two",
				"THREE": "staging-three",
			},
		},
	}

	envSvc := NewEnvironmentService(config)

	_, err = envSvc.Get("nope")
	assert.ErrorIs(t, err, ErrUnknownEnvironment)

	env, err := envSvc.Get("common")
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{
		"ONE": "common-one",
		"TWO": "common-two",
	}, env.Values)

	env, err = envSvc.Get("local")
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{
		"ONE":   "common-one",
		"TWO":   "common-two",
		"THREE": "local-three",
	}, env.Values)

	env, err = envSvc.Get("staging")
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{
		"ONE":   "common-one",
		"TWO":   "staging-two",
		"THREE": "staging-three",
	}, env.Values)
}

func TestEnvironmentService_Get_WhenCircularDependency(t *testing.T) {
	config, err := NewTestConfig()
	assert.NoError(t, err)
	config.Environments = []Environment{
		{
			Name: "aaa",
			Extends: []string{
				"bbb",
			},
			Values: map[string]string{},
		},
		{
			Name: "bbb",
			Extends: []string{
				"aaa",
			},
			Values: map[string]string{},
		},
	}

	envSvc := NewEnvironmentService(config)

	_, err = envSvc.Get("aaa")
	assert.ErrorIs(t, err, ErrCircularDependency)
}

func TestEnvironmentService_Get_WhenUnknownDependency(t *testing.T) {
	config, err := NewTestConfig()
	assert.NoError(t, err)
	config.Environments = []Environment{
		{
			Name: "aaa",
			Extends: []string{
				"nope",
			},
			Values: map[string]string{},
		},
	}

	envSvc := NewEnvironmentService(config)

	_, err = envSvc.Get("aaa")
	assert.ErrorIs(t, err, ErrUnknownEnvironment)
}
