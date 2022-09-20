package planetscale

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const urlFormat = "https://%s/psdb.v1alpha1.Database/%s"

func parseBytes[T interface{}](b []byte) (*T, error) {
	var t *T
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, err
	}

	return t, nil
}

func postHTTP[T interface{}](config *Config, action string, data interface{}) (*T, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(urlFormat, config.Host, action), bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(config.Username, config.Password)

	req.Header.Set("User-Agent", "pscale-serverless-go/0.1.0")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		perr, err := parseBytes[planetscaleError](b)

		if err != nil {
			return nil, err
		} else {
			return nil, perr
		}
	}

	res.Body.Close()
	return parseBytes[T](b)
}
