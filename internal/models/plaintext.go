package models

import (
	"errors"
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var (
	ErrInvalidJSON = errors.New("invalid json")
)

func ReadPlaintext(plaintext string, url URL) (string, error) {
	fragment := url.Fragment
	if fragment != "" {
		// Ensure non-zero.
		if plaintext == "" {
			plaintext = "{}"
		}
		// Ensure valid.
		if !gjson.Valid(plaintext) {
			return "", fmt.Errorf("%w: %s", ErrInvalidJSON, url.String())
		}
		result := gjson.Get(plaintext, fragment)
		plaintext = result.String()
	}
	def, hasDef := url.Default()
	if plaintext == "" && hasDef {
		plaintext = def
	}
	return plaintext, nil
}

func WritePlaintext(plaintext string, url URL, updated string) (string, error) {
	fragment := url.Fragment
	if fragment != "" {
		// Ensure non-zero.
		if plaintext == "" {
			plaintext = "{}"
		}
		// Ensure valid.
		if !gjson.Valid(plaintext) {
			return "", fmt.Errorf("%w: %s", ErrInvalidJSON, url.String())
		}
		// Update.
		updated, _ = sjson.Set(plaintext, fragment, updated)
	}
	return updated, nil
}
