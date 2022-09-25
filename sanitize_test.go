package planetscale

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestSanitize(t *testing.T) {
	query := "SELECT * FROM test WHERE id=?"
	args := map[any]string{
		nil:         "null",
		true:        "true",
		1:           "1",
		1.1:         "1.1",
		time.Time{}: time.Time{}.Format("2006-01-02T15:04:05.999Z07:00"),

		"\"":   "'\\\"'",
		"'":    "'\\''",
		"\n":   "'\\n'",
		"\r":   "'\\r'",
		"\t":   "'\\t'",
		"\\":   "'\\\\'",
		"\x00": "'\\0'",
		"\b":   "'\\b'",
		"\x1a": "'\\Z'",
	}

	for k, v := range args {
		t.Run(fmt.Sprintf("sanitize %s", k), func(t *testing.T) {
			sanitized, err := sanitize(query, []any{k})
			if err != nil {
				t.Fatal(err)
			}

			want := fmt.Sprintf("SELECT * FROM test WHERE id=%v", v)

			if !strings.EqualFold(sanitized, want) {
				t.Fatalf("wanted %s, but got %s", want, sanitized)
			}
		})
	}

	t.Run("arg length", func(t *testing.T) {
		e := errors.New("invalid amount of query args")
		_, err := sanitize(query, []any{"", ""})

		if err.Error() != e.Error() {
			t.Fatalf("wanted %s, but got %s", e, err)
		}
	})
}
