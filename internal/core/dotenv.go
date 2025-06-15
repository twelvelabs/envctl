package core

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"

	"github.com/twelvelabs/termite/fsutil"
)

var (
	DotEnvPathVar = "ENVCTL_DOTENV_PATH"
)

type CleanupFunc func() error

func NewDotEnvService(path string) *DotEnvService {
	if path == "" {
		path = defaultDotEnvDir()
	}
	return &DotEnvService{
		path: path,
	}
}

type DotEnvService struct {
	path string
}

func (s *DotEnvService) Create(vars EnvVars, args []string) (EnvVars, []string, CleanupFunc, error) {
	// Ensure we have a dir to write dotenv files to.
	if err := fsutil.EnsureDirWritable(s.path); err != nil {
		return vars, args, nil, err
	}

	// Create a temp dotenv file there.
	f, err := os.CreateTemp(s.path, "")
	if err != nil {
		return vars, args, nil, err
	}
	defer f.Close()

	// Set the path in both vars and args.
	dotEnvPath := f.Name()
	vars[DotEnvPathVar] = dotEnvPath
	for idx, arg := range args {
		args[idx] = strings.ReplaceAll(arg, DotEnvPathVar, dotEnvPath)
	}

	// Write `vars` to the dotenv file.
	// Using godotenv because it knows how to properly escape values.
	err = godotenv.Write(vars, dotEnvPath)
	if err != nil {
		return vars, args, nil, err
	}

	cleanup := func() error {
		return os.Remove(dotEnvPath)
	}
	return vars, args, cleanup, nil
}

func defaultDotEnvDir() string {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		envDir := filepath.Join(homeDir, ".cache", "envctl")
		return envDir
	}
	return os.TempDir()
}
