package virtual

import (
	"core"
	"reflect"
	"testing"
)

// Test that the virtual links are added to the schema
func TestAugment_schema(t *testing.T) {
	inital_schema := core.Schema{
		Tables: []core.Table{
			{
				Name: "table1",
				Columns: []core.Column{
					{
						Name: "col1",
					},
					{
						Name: "col2",
					},
				},
			},
			{
				Name: "table2",
				Columns: []core.Column{
					{
						Name: "col1",
					},
					{
						Name: "col2",
					},
				},
			},
		},
	}
	experimental_links := []VirtualLink{
		{
			SourceTable:      "table1",
			SourceColumn:     "col1",
			ReferencedTable:  "table2",
			ReferencedColumn: "col1",
		},
	}

	expected := core.Schema{
		Tables: []core.Table{
			{
				Name: "table1",
				Columns: []core.Column{
					{
						Name: "col1",
						Key:  core.Virtual_key,
						Reference: &core.Reference{
							Table:  "table2",
							Column: "col1",
						},
						Extra: "Virtual link to table2.col1",
					},
					{
						Name: "col2",
					},
				},
			},
			{
				Name: "table2",
				Columns: []core.Column{
					{
						Name:  "col1",
						Key:   core.Virtual_key,
						Extra: "Virtual link from table1.col1",
					},
					{
						Name: "col2",
					},
				},
			},
		},
	}

	actual := Augment_schema(inital_schema, experimental_links)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

}
