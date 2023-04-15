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
    },
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