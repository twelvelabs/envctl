package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvVars_Environ(t *testing.T) {
	vars := EnvVars{
		"AAA": "something",
		"BBB": "something with spaces",
		"CCC": "something \"quoted\"",
		"DDD": "something\nmultiline",
	}
	require.Equal(t, []string{
		"AAA=something",
		"BBB=something with spaces",
		"CCC=something \"quoted\"",
		"DDD=something\nmultiline",
	}, vars.Environ())
}

func TestEnvironmentService_Get(t *testing.T) {
	config, err := NewTestConfig()
	assert.NoError(t, err)
	config.Environments = []Environment{
		{
			Name: "common",
			Vars: EnvVars{
				"ONE": "common-one",
				"TWO": "common-two",
			},
		},
		{
			Name: "local",
			Extends: []string{
				"common",
			},
			Vars: EnvVars{
				"THREE": "local-three",
			},
		},
		{
			Name: "staging",
			Extends: []string{
				"common",
			},
			Vars: EnvVars{
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
	assert.Equal(t, EnvVars{
		"ONE": "common-one",
		"TWO": "common-two",
	}, env.Vars)

	env, err = envSvc.Get("local")
	assert.NoError(t, err)
	assert.Equal(t, EnvVars{
		"ONE":   "common-one",
		"TWO":   "common-two",
		"THREE": "local-three",
	}, env.Vars)

	env, err = envSvc.Get("staging")
	assert.NoError(t, err)
	assert.Equal(t, EnvVars{
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
			Vars: EnvVars{},
		},
		{
			Name: "bbb",
			Extends: []string{
				"aaa",
			},
			Vars: EnvVars{},
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
			Vars: EnvVars{},
		},
	}

	envSvc := NewEnvironmentService(config)

	_, err = envSvc.Get("aaa")
	assert.ErrorIs(t, err, ErrUnknownEnvironment)
}
