package planetscale

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
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

	// Count of rows returned.
	Count int

	// Amount of rows affected by a query.
	RowsAffected int64

	// ID of a insert. Blank if no insert is performed.
	InsertID int64

	fields []*resultField
	rows   []*resultRow
}

func (c *Connection) Execute(query string, args ...string) (*Executed, error) {
	q, err := sanitize(query, args)
	if err != nil {
		return nil, err
	}

	r, err := postHTTP[executeResponse](c, "Execute", &executeRequest{Query: q})
	if err != nil {
		return nil, err
	}
	if r.Error != nil {
		return nil, r.Error
	}

	return &Executed{
		Statement: q,

		Count: len(r.Result.Rows),

		RowsAffected: r.Result.RowsAffected,
		InsertID:     r.Result.InsertID,

		fields: r.Result.Fields,
		rows:   r.Result.Rows,
	}, nil
}

func (e *Executed) Decode(out interface{}) error {
	// checks if the interface is a pointer
	if reflect.TypeOf(out).Kind() != reflect.Pointer {
		return errors.New("interface is not a pointer")
	}

	// grabs the value pointed at and checks if it's an array
	s := reflect.Indirect(reflect.ValueOf(out))
	if s.Kind() != reflect.Slice {
		return errors.New("type is not a slice")
	}

	// inner type of the slice
	it := s.Type().Elem().Elem()

	// creates a map of 'ps' tags to the associated field index in the type
	m := make(map[string]int)
	for i := 0; i < it.NumField(); i++ {
		f := it.Field(i)

		v := f.Tag.Get("ps")
		if v == "" {
			continue
		}

		m[v] = i
	}

	// sets the length of the slice to 0
	s.SetLen(0)

	for _, r := range e.rows {
		// decodes all values from base64
		v, err := base64.StdEncoding.DecodeString(r.Values)
		if err != nil {
			return err
		}

		// converts the base64 to a string
		sv := string(v)

		// offset for the lengths
		o := int64(0)

		// creates a new instance of the inner type
		rv := reflect.New(it)

		for li, l := range r.Lengths {
			// converts the length of the field to int64
			n, err := strconv.ParseInt(l, 10, 64)
			if err != nil {
				continue
			}

			// checks if there's a field in the type associated with the one from the table, and casts if it dooes
			if i, ok := m[e.fields[li].Name]; ok {
				cast(sv[o:o+n], reflect.Indirect(rv).Field(i), e.fields[li])
			}

			o += n
		}

		// appends the new instance to the slice and sets it in the output
		s.Set(reflect.Append(s, rv))
	}

	return nil
}

// TODO: make this less vulnerable to panics
func cast(data string, value reflect.Value, field *resultField) error {
	if data == "" {
		return nil
	}

	switch field.Type {
	case "INT8", "INT16", "INT24", "INT32", "INT64":
		n, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		value.SetInt(n)
	case "UINT8", "UINT16", "UINT24", "UINT32", "UINT64":
		n, err := strconv.ParseUint(data, 10, 64)
		if err != nil {
			return nil
		}

		value.SetUint(n)
	case "FLOAT32", "FLOAT64":
		n, err := strconv.ParseFloat(data, value.Type().Bits())
		if err != nil {
			return nil
		}

		value.SetFloat(n)
	default:
		value.Set(reflect.ValueOf(data))
	}

	return nil
}
