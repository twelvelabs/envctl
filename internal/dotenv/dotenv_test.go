package dotenv

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/twelvelabs/envctl/internal/models"
)

func TestDotEnvService(t *testing.T) {
	dir := t.TempDir()
	svc := NewDotEnvService(dir)

	vars, args, cleanup, err := svc.Create(
		models.Vars{
			"AAA": "something",
			"BBB": "something with spaces",
			"CCC": "something \"quoted\"",
			"DDD": "something\nmultiline",
			"EEE": "123",
		},
		[]string{"something", "--file='" + DotEnvPathVar + "'"},
	)
	require.NoError(t, err)

	// dotenv file path should be added to `vars` and `args`.
	dotenvPath := vars[DotEnvPathVar]
	require.NotEmpty(t, dotenvPath.String())
	require.FileExists(t, dotenvPath.String())
	require.Equal(t, fmt.Sprintf("--file='%s'", dotenvPath), args[1])

	// dotenv file content should be properly escaped.
	buf, err := os.ReadFile(dotenvPath.String())
	require.NoError(t, err)
	lines := strings.Split(string(buf), "\n")
	require.Equal(t, []string{
		"AAA=\"something\"",
		"BBB=\"something with spaces\"",
		"CCC=\"something \\\"quoted\\\"\"",
		"DDD=\"something\\nmultiline\"",
		"EEE=123",
		fmt.Sprintf("%s=\"%s\"", DotEnvPathVar, dotenvPath),
		"",
	}, lines)

	// Cleanup function should delete the dotenv file.
	err = cleanup()
	require.NoError(t, err)
	require.NoFileExists(t, dotenvPath.String())
}

func TestDotEnvService_WhenSingleQuote(t *testing.T) {
	dir := t.TempDir()
	svc := NewDotEnvService(dir).WithQuoteStyle(QuoteStyleSingle)

	vars, _, _, err := svc.Create(
		models.Vars{
			"AAA": "something",
			"BBB": "something with spaces",
			"CCC": "something 'quoted'",
			"DDD": "something\nmultiline",
			"EEE": "123",
		},
		[]string{},
	)
	require.NoError(t, err)

	dotenvPath := vars[DotEnvPathVar]
	require.NotEmpty(t, dotenvPath.String())
	require.FileExists(t, dotenvPath.String())

	// dotenv file content should be properly escaped.
	buf, err := os.ReadFile(dotenvPath.String())
	require.NoError(t, err)
	lines := strings.Split(string(buf), "\n")
	require.Equal(t, []string{
		"AAA='something'",
		"BBB='something with spaces'",
		"CCC='something \\'quoted\\''",
		"DDD='something",
		"multiline'",
		"EEE=123",
		fmt.Sprintf("%s='%s'", DotEnvPathVar, dotenvPath),
		"",
	}, lines)
}

func TestDotEnvService_WhenNoQuotes(t *testing.T) {
	dir := t.TempDir()
	svc := NewDotEnvService(dir).WithQuoteStyle(QuoteStyleNone)

	vars, _, _, err := svc.Create(
		models.Vars{
			"AAA": "something",
			"BBB": "something with spaces",
			"CCC": "something 'quoted'",
			"DDD": "something\nmultiline",
			"EEE": "123",
		},
		[]string{},
	)
	require.NoError(t, err)

	dotenvPath := vars[DotEnvPathVar]
	require.NotEmpty(t, dotenvPath.String())
	require.FileExists(t, dotenvPath.String())

	// dotenv file content should be properly escaped.
	buf, err := os.ReadFile(dotenvPath.String())
	require.NoError(t, err)
	lines := strings.Split(string(buf), "\n")
	require.Equal(t, []string{
		"AAA=something",
		"BBB=something with spaces",
		"CCC=something 'quoted'",
		"DDD=something",
		"multiline",
		"EEE=123",
		fmt.Sprintf("%s=%s", DotEnvPathVar, dotenvPath),
		"",
	}, lines)
}

func TestDotEnvService_WhenEmptyArg(t *testing.T) {
	svc := NewDotEnvService("")
	require.Equal(t, defaultDotEnvDir(), svc.path)
}
