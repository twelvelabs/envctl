package core

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/twelvelabs/envctl/internal/models"
)

func TestEnvironmentService_Get(t *testing.T) {
	config, err := NewTestConfig()
	assert.NoError(t, err)
	config.Environments = []Environment{
		{
			Name: "common",
			Vars: models.Vars{
				"ONE": "common-one",
				"TWO": "common-two",
			},
		},
		{
			Name: "local",
			Extends: []string{
				"common",
			},
			Vars: models.Vars{
				"THREE": "local-three",
			},
		},
		{
			Name: "staging",
			Extends: []string{
				"common",
			},
			Vars: models.Vars{
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
	assert.Equal(t, models.Vars{
		"ONE": "common-one",
		"TWO": "common-two",
	}, env.Vars)

	env, err = envSvc.Get("local")
	assert.NoError(t, err)
	assert.Equal(t, models.Vars{
		"ONE":   "common-one",
		"TWO":   "common-two",
		"THREE": "local-three",
	}, env.Vars)

	env, err = envSvc.Get("staging")
	assert.NoError(t, err)
	assert.Equal(t, models.Vars{
		"ONE":   "common-one",
		"TWO":   "staging-two",
		"THREE": "staging-three",
	}, env.Vars)
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
			Vars: models.Vars{},
		},
		{
			Name: "bbb",
			Extends: []string{
				"aaa",
			},
			Vars: models.Vars{},
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
			Vars: models.Vars{},
		},
	}

	envSvc := NewEnvironmentService(config)

	_, err = envSvc.Get("aaa")
	assert.ErrorIs(t, err, ErrUnknownEnvironment)
}
