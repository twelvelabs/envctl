package models

import (
	"fmt"
	"net/url"
	"sort"
)

// URL is lightweight wrapper around [url.URL] that
// adds a few convenience methods.
type URL struct {
	*url.URL
}

// Default returns the `default` query param,
// and whether that params exists in the URL.
func (u URL) Default() (string, bool) {
	q := u.Query()
	return q.Get("default"), q.Has("default")
}

// Value is a value to inject into the environment.
type Value string

// String returns v as a string.
func (v Value) String() string {
	return string(v)
}

// URL returns v as a URL.
func (v Value) URL() URL {
	parsed, _ := url.Parse(v.String())
	if parsed == nil {
		parsed = &url.URL{}
	}
	return URL{URL: parsed}
}

// Vars is a map of key/value pairs
// to use as environment variables.
type Vars map[string]Value

// Environ returns a list of key/value pairs in the form
// of "key=value" (similar to [os.Environ]).
func (vars Vars) Environ() []string {
	pairs := []string{}
	for k, v := range vars {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	// Ensure stable sort order for tests.
	sort.Strings(pairs)
	return pairs
}

// Get returns the value for key and true if key exists
// in the collection.
func (vars Vars) Get(key string) (Value, bool) {
	value, found := vars[key]
	return value, found
}

func (vars Vars) Map() map[string]string {
	out := map[string]string{}
	for k, v := range vars {
		out[k] = string(v)
	}
	return out
}

func (vars Vars) Pluck(keys ...string) Vars {
	plucked := Vars{}
	for _, k := range keys {
		if v, ok := vars[k]; ok {
			plucked[k] = v
		}
	}
	return plucked
}
