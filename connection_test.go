package planetscale

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type testConnectionClient struct{}

func (c *testConnectionClient) RoundTrip(req *http.Request) (*http.Response, error) {
	s := []byte(`{}`)
	res := &createSessionResponse{
		Session: (*json.RawMessage)(&s),
	}

	m, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader(m)),
	}, nil
}

func TestConnection(t *testing.T) {
	_, err := NewConnection(&Config{
		Username: "",
		Password: "",
		Host:     "",

		Transport: &testConnectionClient{},
	})

	if err != nil {
		t.Fatalf("got error: %s", err)
	}
}
