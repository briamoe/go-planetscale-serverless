package planetscale

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type castTest struct {
	IntTest    int     `ps:"intTest"`
	UIntTest   uint    `ps:"uintTest"`
	FloatTest  float32 `ps:"floatTest"`
	StringTest string  `ps:"stringTest"`
}

type testSelectClient struct{}

func (c *testSelectClient) RoundTrip(req *http.Request) (*http.Response, error) {
	res := &executeResponse{
		Session: (json.RawMessage)([]byte(`{}`)),
		Result: &executeResult{
			Fields: []*resultField{
				{
					Name:  "intTest",
					Type:  "INT32",
					Table: "test",
				},
				{
					Name:  "uintTest",
					Type:  "UINT32",
					Table: "test",
				},
				{
					Name:  "floatTest",
					Type:  "FLOAT32",
					Table: "test",
				},
				{
					Name:  "stringTest",
					Type:  "STRING",
					Table: "test",
				},
			},
			Rows: []*resultRow{
				{
					Lengths: []string{"2", "1", "3", "4"},
					Values:  "LTExMS4xdGVzdA==", // base64 encoded -111.1test
				},
			},
		},
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

type testInsertClient struct{}

func (c *testInsertClient) RoundTrip(req *http.Request) (*http.Response, error) {
	res := &executeResponse{
		Session: (json.RawMessage)([]byte(`{}`)),
		Result: &executeResult{
			RowsAffected: 1,
			InsertID:     1,
		},
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

func TestExecute(t *testing.T) {
	conn, err := NewConnection(&Config{
		Username: "",
		Password: "",
		Host:     "",

		Transport: &testConnectionClient{},
	})
	if err != nil {
		t.Fatalf("got error: %s", err)
	}

	t.Run("select", func(t *testing.T) {
		nconn := &Connection{
			Config:  conn.Config,
			Session: conn.Session,
			client: &http.Client{
				Transport: &testSelectClient{},
			},
		}

		e, err := nconn.Execute("SELECT * FROM test")
		if err != nil {
			t.Fatalf("got error: %s", err)
		}

		var c []*castTest
		if err := e.Decode(&c); err != nil {
			t.Fatalf("got error: %s", err)
		}

		want := []*castTest{
			{
				IntTest:    -1,
				UIntTest:   1,
				FloatTest:  1.1,
				StringTest: "test",
			},
		}

		if !reflect.DeepEqual(c[0], want[0]) {
			t.Fatalf("wanted %+v, but got %+v", want[0], c[0])
		}
	})

	t.Run("insert", func(t *testing.T) {
		nconn := &Connection{
			Config:  conn.Config,
			Session: conn.Session,
			client: &http.Client{
				Transport: &testInsertClient{},
			},
		}

		e, err := nconn.Execute("INSERT INTO test (id) VALUES ('1')")
		if err != nil {
			t.Fatalf("got error: %s", err)
		}

		if e.InsertID != 1 {
			t.Fatalf("wanted %v, but got %v", 1, e.InsertID)
		}

		if e.RowsAffected != 1 {
			t.Fatalf("wanted %v, but got %v", 1, e.RowsAffected)
		}
	})
}
