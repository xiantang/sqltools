package sqltools

import (
	"reflect"
	"strings"
	"testing"
)

func Test_isWhitespace(t *testing.T) {
	type args struct {
		ch rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test", args{ch: ' '}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isWhitespace(tt.args.ch); got != tt.want {
				t.Errorf("isWhitespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isLetter(t *testing.T) {
	type args struct {
		ch rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test", args{ch: ' '}, false},
		{"test", args{ch: 'a'}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLetter(tt.args.ch); got != tt.want {
				t.Errorf("isLetter() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Ensure the parser can parse strings into Statement ASTs.
func TestParser_ParseSelectStatement(t *testing.T) {
	var tests = []struct {
		s    string
		stmt *SelectStatement
		err  string
	}{
		// Single field statement
		{
			s: `SELECT name FROM tbl`,
			stmt: &SelectStatement{
				Fields:    []string{"name"},
				TableName: "tbl",
			},
		},

		// Multi-field statement
		{
			s: `SELECT first_name, last_name, age FROM my_table`,
			stmt: &SelectStatement{
				Fields:    []string{"first_name", "last_name", "age"},
				TableName: "my_table",
			},
		},

		// Select all statement
		{
			s: `SELECT * FROM my_table`,
			stmt: &SelectStatement{
				Fields:    []string{"*"},
				TableName: "my_table",
			},
		},

		// Errors
		{s: `foo`, err: `found "foo", expected SELECT or ALTER`},
		{s: `SELECT !`, err: `found "!", expected field`},
		{s: `SELECT field xxx`, err: `found "xxx", expected FROM`},
		{s: `SELECT field FROM *`, err: `found "*", expected table name`},
	}

	for i, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			stmt, err := NewParser(strings.NewReader(tt.s)).Parse()
			if !reflect.DeepEqual(tt.err, errstring(err)) {
				t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
			} else if tt.err == "" && !reflect.DeepEqual(tt.stmt, stmt) {
				t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.stmt, stmt)
			}
		})

	}
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

// Ensure the parser can parse strings into Statement ASTs.
func TestParser_ParseAlterStatement(t *testing.T) {
	var tests = []struct {
		s    string
		stmt *AlterStatement
		err  string
	}{
		// Single field statement
		{
			s: `ALTER TABLE table_name DROP COLUMN column_name;`,
			stmt: &AlterStatement{
				Option: DROP,
				Column: ColumnStatement{
					ColumnName: "column_name",
				},
				TableName: "table_name",
			},
		},
		{
			s: `alter table xxxxx add xxxx2 varchar(255) null`,
			stmt: &AlterStatement{
				Option: ADD,
				Column: ColumnStatement{
					ColumnName: "xxxx2",
					DataType:   VARCHAR,
					Length:     255,
					Nullable:   true,
				},
				TableName: "xxxxx",
			},
		},
	}
	for i, tt := range tests {

		stmt, err := NewParser(strings.NewReader(tt.s)).Parse()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.stmt, stmt) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.stmt, stmt)
		}

	}
}
