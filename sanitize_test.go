package planetscale

import (
	"errors"
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
				t.Fatalf("wanted %s, but got %s", want, sanitized)
			}
		})
	}

	t.Run("arg length", func(t *testing.T) {
		e := errors.New("invalid amount of query args")
		_, err := sanitize(query, []string{"", ""})

		if err.Error() != e.Error() {
			t.Fatalf("wanted %s, but got %s", e, err)
		}
	})
}
