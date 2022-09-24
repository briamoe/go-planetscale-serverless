package planetscale

import (
	"errors"
	"strings"
)

func sanitize(query string, args []string) (string, error) {
	if strings.Count(query, "?") != len(args) {
		return "", errors.New("invalid amount of query args")
	}

	// escape important values
	r := strings.NewReplacer(
		"\"", "\\\"",
		"'", "\\'",
		"\n", "\\n",
		"\r", "\\r",
		"\t", "\\t",
		"\\", "\\\\",
		"\x00", "\\0",
		"\b", "\\b",
		"\x1a", "\\Z",
	)

	var f strings.Builder
	i := 0

	for _, c := range query {
		if c != '?' {
			f.WriteRune(c)
			continue
		}

		f.WriteRune('\'')
		f.WriteString(r.Replace(args[i]))
		f.WriteRune('\'')

		i++
	}

	return f.String(), nil
}
