package core

// The d2 tests can be done by using strings.TrimSpace() to remove the whitespace from the output of the d2 generator.
// This should remove the influence of formatting differences between the two strings. This is a simple way to test that that two are equal.
// In future, we may want to use a d2 formatter to make the output more consistent.

import (
	"strings"
	"testing"
)

func Test_d2_generator(t *testing.T) {
	schema := Schema{
		Tables: []Table{
			{
				Name: "table1",
				Columns: []Column{
					{
						Name: "col1",
						Type: "int",
					},
					{
						Name: "col2",
						Type: "varchar",
					},
				},
			},
			{
				Name: "table2",
				Columns: []Column{
					{
						Name: "col1",
						Type: "int",
					},
					{
						Name: "col2",
						Type: "varchar",
					},
				},
			},
		},
	}

expected := `table1: {
  shape: sql_table
  col1: int
  col2: varchar
}

table2: {
  shape: sql_table
  col1: int
  col2: varchar
}`
	expected = strings.TrimSpace(expected)

	actual := strings.TrimSpace(Schema_to_d2(schema, []TableGroup{}))

	if actual != expected {
		t.Errorf("Expected \n%s\n, got \n%s", expected, actual)
	}

}
