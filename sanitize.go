package planetscale

import (
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"
)

func sanitize(query string, args []any) (string, error) {
	if strings.Count(query, "?") != len(args) {
		return "", errors.New("invalid amount of query args")
	}

	var f strings.Builder
	i := 0

	for _, c := range query {
		if c != '?' {
			f.WriteRune(c)
			continue
		}

		if err := sanitizeArg(args[i], &f); err != nil {
			return "", err
		}

		i++
	}

	return f.String(), nil
}

func sanitizeArg(arg any, f *strings.Builder) error {
	switch a := arg.(type) {
	case nil:
		f.WriteString("null")
	case bool:
		f.WriteString(strconv.FormatBool(a))
	case string:
		f.WriteRune('\'')
		f.WriteString(escapeString(a))
		f.WriteRune('\'')
	case int:
		f.WriteString(strconv.FormatInt(int64(a), 10))
	case int8:
		f.WriteString(strconv.FormatInt(int64(a), 10))
	case int16:
		f.WriteString(strconv.FormatInt(int64(a), 10))
	case int32:
		f.WriteString(strconv.FormatInt(int64(a), 10))
	case int64:
		f.WriteString(strconv.FormatInt(a, 10))
	case float32:
		f.WriteString(strconv.FormatFloat(float64(a), 'f', -1, 32))
	case float64:
		f.WriteString(strconv.FormatFloat(a, 'f', -1, 64))
	case time.Time:
		f.WriteString(a.UTC().Format("2006-01-02T15:04:05.999Z07:00"))
	case []byte:
		f.WriteString(`'\x`)
		f.WriteString(hex.EncodeToString(a))
		f.WriteRune('\'')
	default:
		return errors.New("invalid type passed to query args")
	}

	return nil
}

func escapeString(input string) string {
	// escapes important values
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

	return r.Replace(input)
}
