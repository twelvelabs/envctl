package core

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDotEnvService(t *testing.T) {
	dir := t.TempDir()
	svc := NewDotEnvService(dir)

	vars, args, cleanup, err := svc.Create(
		EnvVars{
			"AAA": "something",
			"BBB": "something with spaces",
			"CCC": "something \"quoted\"",
			"DDD": "something\nmultiline",
		},
		[]string{"something", "--file='" + DotEnvPathVar + "'"},
	)
	require.NoError(t, err)

	// dotenv file path should be added to `vars` and `args`.
	dotenvPath := vars[DotEnvPathVar]
	require.NotEmpty(t, dotenvPath)
	require.FileExists(t, dotenvPath)
	require.Equal(t, fmt.Sprintf("--file='%s'", dotenvPath), args[1])

	// dotenv file content should be properly escaped.
	buf, err := os.ReadFile(dotenvPath) //nolint:gosec
	require.NoError(t, err)
	lines := strings.Split(string(buf), "\n")
	require.Equal(t, []string{
		"AAA=\"something\"",
		"BBB=\"something with spaces\"",
		"CCC=\"something \\\"quoted\\\"\"",
		"DDD=\"something\\nmultiline\"",
	}, lines[0:4])
	require.Contains(t, lines[4], fmt.Sprintf("%s=\"%s\"", DotEnvPathVar, dotenvPath))

	// Cleanup function should delete the dotenv file.
	err = cleanup()
	require.NoError(t, err)
	require.NoFileExists(t, dotenvPath)
}

func TestDotEnvService_WhenEmptyArg(t *testing.T) {
	svc := NewDotEnvService("")
	require.Equal(t, defaultDotEnvDir(), svc.path)
}
