package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// JSONEnvelope wraps all JSON CLI output in a consistent structure.
type JSONEnvelope struct {
	SchemaVersion int         `json:"schema_version"`
	Command       string      `json:"command"`
	Data          interface{} `json:"data,omitempty"`
	Metadata      interface{} `json:"metadata,omitempty"`
	Error         *JSONError  `json:"error,omitempty"`
}

// JSONError represents a structured error in JSON output.
type JSONError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// OutputFormatter handles rendering CLI output in either
// human-readable (tabwriter) or JSON format.
type OutputFormatter struct {
	w      io.Writer
	isJSON bool
}

// NewOutputFormatter creates a formatter that writes to w.
// When jsonMode is true, output is structured JSON.
func NewOutputFormatter(w io.Writer, jsonMode bool) *OutputFormatter {
	return &OutputFormatter{w: w, isJSON: jsonMode}
}

// IsJSON reports whether the formatter is in JSON mode.
func (f *OutputFormatter) IsJSON() bool {
	return f.isJSON
}

// WriteJSON writes a successful JSON response envelope.
func (f *OutputFormatter) WriteJSON(command string, data interface{}, metadata interface{}) error {
	env := JSONEnvelope{
		SchemaVersion: 1,
		Command:       command,
		Data:          data,
		Metadata:      metadata,
	}
	enc := json.NewEncoder(f.w)
	enc.SetIndent("", "  ")
	return enc.Encode(env)
}

// WriteJSONError writes a JSON error response envelope.
func (f *OutputFormatter) WriteJSONError(command string, code int, message, detail string) error {
	env := JSONEnvelope{
		SchemaVersion: 1,
		Command:       command,
		Error: &JSONError{
			Code:    code,
			Message: message,
			Detail:  detail,
		},
	}
	enc := json.NewEncoder(f.w)
	enc.SetIndent("", "  ")
	return enc.Encode(env)
}

// TableWriter returns a tabwriter configured for aligned columnar output.
// Callers must call Flush on the returned writer when done.
func (f *OutputFormatter) TableWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(f.w, 0, 0, 2, ' ', 0)
}

// Writef writes a formatted string to the output writer.
func (f *OutputFormatter) Writef(format string, args ...interface{}) error {
	_, err := fmt.Fprintf(f.w, format, args...)
	return err
}
