package core

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigFromPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name      string
		args      args
		want      *Config
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "should return valid config",
			args: args{
				path: filepath.Join("testdata", "config", "valid.yaml"),
			},
			want: &Config{
				ConfigPath: filepath.Join("testdata", "config", "valid.yaml"),
				Debug:      true,
			},
			assertion: assert.NoError,
		},
		{
			name: "should return error if malformed",
			args: args{
				path: filepath.Join("testdata", "config", "malformed.yaml"),
			},
			want:      nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfigFromPath(tt.args.path)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
