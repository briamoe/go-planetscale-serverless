package planetscale

import (
	"encoding/json"
)

type resultField struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Table string `json:"table"`

	OrgTable string `json:"orgTable"`
	Database string `json:"database"`
	OrgName  string `json:"orgName"`

	ColumnLength int `json:"columnLength"`
	Charset      int `json:"charset"`
	Flags        int `json:"flags"`
}

type resultRow struct {
	Lengths []string
	Values  string
}

type executeResult struct {
	Fields []*resultField `json:"fields"`
	Rows   []*resultRow   `json:"rows"`

	RowsAffected int64 `json:"rowsAffected,string"`
	InsertID     int64 `json:"insertId,string"`
}

type executeResponse struct {
	Session json.RawMessage   `json:"session"`
	Result  *executeResult    `json:"result"`
	Error   *planetscaleError `json:"error"`
}

type executeRequest struct {
	Query string `json:"query"`
}

type Executed struct {
	// Final statement passed to PlanetScale.
	Statement string

	// Amount of rows affected by a query.
	RowsAffected int64

	// ID of a insert. Blank if no insert is performed.
	InsertID int64

	fields []*resultField
	rows   []*resultRow
}

func (c *Connection) Execute(query string, args ...string) (*Executed, error) {
	// TODO: sanitize & add args
	q := query

	r, err := postHTTP[executeResponse](c.Config, "Execute", &executeRequest{Query: q})
	if err != nil {
		return nil, err
	}
	if r.Error != nil {
		return nil, r.Error
	}

	return &Executed{
		Statement: q,

		RowsAffected: r.Result.RowsAffected,
		InsertID:     r.Result.InsertID,

		fields: r.Result.Fields,
		rows:   r.Result.Rows,
	}, nil
}
