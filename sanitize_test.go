package planetscale

import (
	"fmt"
	"strings"
	"testing"
)

func TestSanitize(t *testing.T) {
	query := "SELECT * FROM test WHERE id=?"
	args := map[string]string{
		"\"":   "\\\"",
		"'":    "\\'",
		"\n":   "\\n",
		"\r":   "\\r",
		"\t":   "\\t",
		"\\":   "\\\\",
		"\x00": "\\0",
		"\b":   "\\b",
		"\x1a": "\\Z",
	}

	for k, v := range args {
		t.Run(fmt.Sprintf("sanitize %s", k), func(t *testing.T) {
			sanitized, err := sanitize(query, []string{k})
			if err != nil {
				t.Fatal(err)
			}

			want := fmt.Sprintf("SELECT * FROM test WHERE id='%s'", v)

			if !strings.EqualFold(sanitized, want) {
				t.Fatalf("sanitize %s: want %s, got %s", k, want, sanitized)
			}
		})
	}
}
