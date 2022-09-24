package planetscale

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type postTest struct {
	Test string `json:"test"`
}

type testPostClient struct{}

func (c *testPostClient) RoundTrip(req *http.Request) (*http.Response, error) {
	res := &postTest{
		Test: "test",
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

func TestPostHTTP(t *testing.T) {
	c := &Connection{
		Config: &Config{},
		client: &http.Client{
			Transport: &testPostClient{},
		},
	}

	p, err := postHTTP[postTest](c, "", &postTest{})
	if err != nil {
		t.Fatalf("got error: %s", err)
	}

	if p.Test != "test" {
		t.Fatalf("wanted %+v, but got %+v", "test", p.Test)
	}
}
