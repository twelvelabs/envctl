package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadPlaintext(t *testing.T) {
	json := `{"AAA": "123", "BBB": "456"}`

	plaintext, err := ReadPlaintext(json, Value("foo://bar").URL())
	require.NoError(t, err)
	require.Equal(t, json, plaintext)

	plaintext, err = ReadPlaintext(json, Value("foo://bar#BBB").URL())
	require.NoError(t, err)
	require.Equal(t, "456", plaintext)

	plaintext, err = ReadPlaintext(json, Value("foo://bar#LOL.WAT").URL())
	require.NoError(t, err)
	require.Equal(t, "", plaintext)

	plaintext, err = ReadPlaintext("nope", Value("foo://bar#BBB").URL())
	require.ErrorIs(t, err, ErrInvalidJSON)
	require.Equal(t, "", plaintext)
}

func TestWritePlaintext(t *testing.T) {
	json := `{"AAA": "123", "BBB": "456"}`

	plaintext, err := WritePlaintext(json, Value("foo://bar").URL(), "howdy")
	require.NoError(t, err)
	require.Equal(t, "howdy", plaintext)

	plaintext, err = WritePlaintext(json, Value("foo://bar#BBB").URL(), "howdy")
	require.NoError(t, err)
	require.Equal(t, `{"AAA": "123", "BBB": "howdy"}`, plaintext)

	plaintext, err = WritePlaintext("", Value("foo://bar#BBB").URL(), "howdy")
	require.NoError(t, err)
	require.Equal(t, `{"BBB":"howdy"}`, plaintext)

	plaintext, err = WritePlaintext(json, Value("foo://bar#LOL.WAT").URL(), "howdy")
	require.NoError(t, err)
	require.Equal(t, `{"AAA": "123", "BBB": "456","LOL":{"WAT":"howdy"}}`, plaintext)

	plaintext, err = WritePlaintext("nope", Value("foo://bar#BBB").URL(), "howdy")
	require.ErrorIs(t, err, ErrInvalidJSON)
	require.Equal(t, "", plaintext)
}
