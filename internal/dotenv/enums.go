package dotenv

// EscapeStyle determines how to escape double-quoted string values.
type EscapeStyle string

const (
	// Default escaping behavior.
	EscapeStyleDefault EscapeStyle = "default"
	// Docker Compose escaping behavior.
	EscapeStyleCompose EscapeStyle = "compose"
)

// QuoteStyle determines how to quote dotenv values.
type QuoteStyle string

const (
	QuoteStyleNone   QuoteStyle = "none"   // Do not quote values.
	QuoteStyleSingle QuoteStyle = "single" // Single quote values.
	QuoteStyleDouble QuoteStyle = "double" // Double quote values.
)
