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

func postHTTP[T interface{}](conn *Connection, action string, data interface{}) (*T, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(urlFormat, conn.Config.Host, action), bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(conn.Config.Username, conn.Config.Password)

	req.Header.Set("User-Agent", "go-planetscale-serverless/0.2.0")
	req.Header.Set("Content-Type", "application/json")

	res, err := conn.client.Do(req)
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
