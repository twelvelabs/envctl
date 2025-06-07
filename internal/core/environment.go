package core

import (
	"errors"
	"fmt"
	"maps"

	mapset "github.com/deckarep/golang-set/v2"
)

var (
	ErrCircularDependency = errors.New("circular dependency")
	ErrUnknownEnvironment = errors.New("unknown environment")
)

type Environment struct {
	Name    string            `yaml:"name"`
	Extends []string          `yaml:"extends"`
	Values  map[string]string `yaml:"values"`
}

func NewEnvironmentService(config *Config) *EnvironmentService {
	lookup := map[string]Environment{}
	for _, env := range config.Environments {
		lookup[env.Name] = env
	}

	return &EnvironmentService{
		config: config,
		lookup: lookup,
		getter: DefaultGetter,
	}
}

type EnvironmentService struct {
	config *Config
	lookup map[string]Environment
	getter Getter
}

func (s *EnvironmentService) Get(id string) (Environment, error) {
	env, err := s.get(id)
	if err != nil {
		return env, err
	}

	env, err = s.expand(env, nil)
	if err != nil {
		return env, err
	}

	return env, nil
}

func (s *EnvironmentService) get(id string) (Environment, error) {
	// Most likely this is an environment defined in the config file.
	// If so, just return out of the lookup.
	if env, ok := s.lookup[id]; ok {
		return env, nil
	}

	env := Environment{}
	return env, fmt.Errorf("%w: %s", ErrUnknownEnvironment, id)
}

func (s *EnvironmentService) expand(env Environment, seen mapset.Set[string]) (Environment, error) {
	// Check for cycles.
	if seen == nil {
		seen = mapset.NewSet[string]()
	}
	if seen.Contains(env.Name) {
		return env, ErrCircularDependency
	}
	seen.Add(env.Name)

	// Collect the values from each ancestor in the `extends` list.
	// Later entries may overwrite values from previous ones.
	values := map[string]string{}
	for _, id := range env.Extends {
		ancestor, err := s.get(id)
		if err != nil {
			return env, err
		}
		// Recursively expand each ancestor.
		ancestor, err = s.expand(ancestor, seen)
		if err != nil {
			return env, err
		}
		maps.Copy(values, ancestor.Values)
	}
	// Finally, overwrite w/ values from the child and update.
	maps.Copy(values, env.Values)
	env.Values = values
	return env, nil
}
