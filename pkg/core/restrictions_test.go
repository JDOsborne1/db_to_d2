package core

import (
	"reflect"
	"testing"
)

// TestRestrictionsIdentity is a basic test that checks that the identity function returns the same value as the input.
func TestRestrictionsIdentity(t *testing.T) {
	//Test data
	input := Schema{
		Tables: []Table{
			{
				Name: "test",
				Columns: []Column{
					{
						Name: "test",
					},
				},
			},
		},
	}
	expected := input

	//Execute test
	actual := Restrict(input, Standard)

	//Compare actual to expected
	if !reflect.DeepEqual(actual, expected) {
		t.Log("Identity function failed to return the same value as the input.")
		t.Fail()
	}
}
