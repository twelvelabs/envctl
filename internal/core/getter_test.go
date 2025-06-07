package core

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultGetter(t *testing.T) {
	srcPath := filepath.Join("testdata", "config", "valid.yaml")
	dstPath := filepath.Join(t.TempDir(), "valid.yaml")
	assert.NoFileExists(t, dstPath)
	err := DefaultGetter(t.Context(), srcPath, dstPath)
	assert.NoError(t, err)
	assert.FileExists(t, dstPath)
}

func TestDefaultGetter_WhenError(t *testing.T) {
	err := DefaultGetter(t.Context(), "does/not/exist", "does/not/exist")
	assert.ErrorContains(t, err, "no such file or directory")
}

func TestMockGetter(t *testing.T) {
	getter := NewMockGetter(func() error {
		return nil
	})

	err := getter.Get(t.Context(), "src/path", "dst/path")
	assert.NoError(t, err)

	assert.Equal(t, true, getter.Called)
	assert.Equal(t, "src/path", getter.Src)
	assert.Equal(t, "dst/path", getter.Dst)
}
