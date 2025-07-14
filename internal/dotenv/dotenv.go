package dotenv

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/twelvelabs/termite/fsutil"

	"github.com/twelvelabs/envctl/internal/models"
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
		path:        path,
		quoteStyle:  QuoteStyleDouble,
		escapeStyle: EscapeStyleDefault,
	}
}

type DotEnvService struct {
	path        string
	quoteStyle  QuoteStyle
	escapeStyle EscapeStyle
}

func (s *DotEnvService) WithQuoteStyle(qs QuoteStyle) *DotEnvService {
	s.quoteStyle = qs
	return s
}

func (s *DotEnvService) WithEscapeStyle(es EscapeStyle) *DotEnvService {
	s.escapeStyle = es
	return s
}

func (s *DotEnvService) Create(vars models.Vars, args []string) (models.Vars, []string, CleanupFunc, error) { //nolint:lll
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
	vars[DotEnvPathVar] = models.Value(dotEnvPath)
	for idx, arg := range args {
		args[idx] = strings.ReplaceAll(arg, DotEnvPathVar, dotEnvPath)
	}

	// Write `vars` to the dotenv file.
	// Using godotenv because it knows how to properly escape values.
	err = s.Write(vars, dotEnvPath)
	if err != nil {
		return vars, args, nil, err
	}

	cleanup := func() error {
		return os.Remove(dotEnvPath)
	}
	return vars, args, cleanup, nil
}

// Marshal serializes the given vars as a dotenv-formatted string.
// Each line is in the format: KEY=VALUE where VALUE is determined
// by the configured quote and escape styles.
func (s *DotEnvService) Marshal(vars models.Vars) string {
	envMap := vars.Map()
	lines := make([]string, 0, len(envMap))
	for k, v := range envMap {
		if d, err := strconv.Atoi(v); err == nil {
			lines = append(lines, fmt.Sprintf(`%s=%d`, k, d))
		} else {
			lines = append(lines, fmt.Sprintf(`%s=%s`, k, s.quote(v)))
		}
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

func (s *DotEnvService) quote(value string) string {
	value = s.escape(value)
	switch s.quoteStyle {
	case QuoteStyleDouble:
		return fmt.Sprintf(`"%s"`, value)
	case QuoteStyleSingle:
		return fmt.Sprintf(`'%s'`, value)
	default:
		return value
	}
}

const (
	doubleQuoteCharsCompose = "\\\n\r\t\""
	doubleQuoteCharsDefault = doubleQuoteCharsCompose + "!$`"
	singleQuoteCharsDefault = "'"
)

func (s *DotEnvService) specialChars() string {
	switch s.quoteStyle {
	case QuoteStyleDouble:
		if s.escapeStyle == EscapeStyleCompose {
			return doubleQuoteCharsCompose
		}
		return doubleQuoteCharsDefault
	case QuoteStyleSingle:
		return singleQuoteCharsDefault
	default:
		return ""
	}
}

func (s *DotEnvService) escape(value string) string {
	for _, c := range s.specialChars() {
		toReplace := "\\" + string(c)
		if c == '\n' {
			toReplace = `\n`
		}
		if c == '\r' {
			toReplace = `\r`
		}
		value = strings.ReplaceAll(value, string(c), toReplace)
	}
	return value
}

// Write serializes the given vars and writes them to path.
func (s *DotEnvService) Write(vars models.Vars, path string) error {
	f, err := os.Create(path) //nolint: gosec
	if err != nil {
		return err
	}
	defer f.Close()

	content := s.Marshal(vars)
	_, err = f.WriteString(content + "\n")
	if err != nil {
		return err
	}

	return f.Sync()
}

func defaultDotEnvDir() string {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		envDir := filepath.Join(homeDir, ".cache", "envctl")
		return envDir
	}
	return os.TempDir()
}
