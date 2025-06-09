package core

import (
	"maps"
	"net/url"

	"github.com/twelvelabs/envctl/internal/resolvers/env"
)

var (
	Resolvers = map[string]Resolver{
		"env": env.NewEnvResolver(),
	}
)

type Resolver interface {
	Resolve(url *url.URL) (string, error)
}

func NewResolverService(config *Config, resolvers map[string]Resolver) *ResolverService {
	return &ResolverService{
		config:    config,
		resolvers: resolvers,
	}
}

type ResolverService struct {
	config    *Config
	resolvers map[string]Resolver
}

func (s *ResolverService) ResolveVars(vars EnvVars) (EnvVars, error) {
	resolved := maps.Clone(vars)
	for key, val := range vars {
		// Try to parse the value into a URL.
		u, _ := url.Parse(val)
		if u == nil || u.Scheme == "" {
			continue
		}
		// Check if it's a resolvable protocol.
		resolver, ok := s.resolvers[u.Scheme]
		if !ok {
			continue
		}
		// If so, try to resolve.
		val, err := resolver.Resolve(u)
		if err != nil {
			return resolved, err
		}
		resolved[key] = val
	}
	return resolved, nil
}
