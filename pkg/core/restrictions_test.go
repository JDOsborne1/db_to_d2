package core

import (
	"reflect"
	"testing"
)

var example_schema = Schema{
	Tables: []Table{
		{
			Name: "users",
			Columns: []Column{
				{
					Name:     "id",
					Type:     "int",
					Nullable: false,
					Key:      "PRIMARY KEY",
					Extra:    "AUTO_INCREMENT",
				},
				{
					Name:     "username",
					Type:     "varchar(255)",
					Nullable: false,
					Key:      "UNIQUE KEY",
					Extra:    "",
				},
				{
					Name:     "email",
					Type:     "varchar(255)",
					Nullable: false,
					Key:      "UNIQUE KEY",
					Extra:    "",
				},
			},
		},
		{
			Name: "posts",
			Columns: []Column{
				{
					Name:     "id",
					Type:     "int",
					Nullable: false,
					Key:      "PRIMARY KEY",
					Extra:    "AUTO_INCREMENT",
				},
				{
					Name:     "title",
					Type:     "varchar(255)",
					Nullable: false,
					Key:      "",
					Extra:    "",
				},
				{
					Name:     "content",
					Type:     "text",
					Nullable: true,
					Key:      "",
					Extra:    "",
				},
				{
					Name:     "user_id",
					Type:     "int",
					Nullable: false,
					Key:      "FOREIGN KEY",
					Extra:    "",
					Reference: &Reference{
						Table:    "users",
						Column:   "id",
						OnDelete: "CASCADE",
						OnUpdate: "CASCADE",
					},
				},
			},
		},
		{
			Name: "meta",
			Columns: []Column{
				{
					Name:     "meta_key",
					Type:     "varchar(255)",
					Nullable: false,
					Key:      "",
					Extra:    "",
				},
				{
					Name:     "meta_value",
					Type:     "varchar(255)",
					Nullable: false,
					Key:      "",
					Extra:    "",
				},
			},
		},
	},
}

func get_table_names(_input Schema) []string {
	var output []string
	for _, table := range _input.Tables {
		output = append(output, table.Name)
	}
	return output
}

func get_column_names(_input Table) []string {
	var output []string
	for _, column := range _input.Columns {
		output = append(output, column.Name)
	}
	return output
}

// TestRestrictionsIdentity is a basic test that checks that the identity function returns the same value as the input.
func TestRestrictionsIdentity(t *testing.T) {
	//Execute test
	actual := Restrict(example_schema, Standard)

	//Compare actual to expected
	if !reflect.DeepEqual(actual, example_schema) {
		t.Log("Identity function failed to return the same value as the input.")
		t.Fail()
	}
}

func TestRestrictionsMinimalist(t *testing.T) {
	actual := Restrict(example_schema, Minimalist)
	names := get_table_names(actual)
	if !(equal_set(names, []string{"users", "posts"})) {
		t.Log("Minimalist cleared inappropriate tables.")
		t.Fail()
	}

	if !(equal_set(get_column_names(actual.Tables[0]), []string{"id", "username", "email"})) {
		t.Log("Minimalist cleared columns in a table where all the columns are keys.")
		t.Fail()
	}

	if !(equal_set(get_column_names(actual.Tables[1]), []string{"id", "user_id"})) {
		t.Log("Minimalist didn't clear the right columns in a table which has partial keys.")
		t.Fail()
	}

	t.Log(actual.Tables)
	t.Fail()
}
