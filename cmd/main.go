package main

import (
	"fmt"
)

type Column struct {
	Name      string
	Type      string
	Nullable  bool
	Key       string
	Extra     string
	Reference *Reference
}

type Reference struct {
	Table    string
	Column   string
	OnDelete string
	OnUpdate string
}

type Table struct {
	Name    string
	Columns []Column
}

type Schema struct {
	Tables  []Table
	Indexes []string
}

type VirtualLink struct {
	source_table string
	source_column string
	referenced_table string
	referenced_column string
}

func augment_schema_with_virtual(_input Schema, _links VirtualLink) Schema {
	return _input
}

func main() {
	// Connect to the MySQL database
	db, err := connect_to_db()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db_schema := information_schema_from(db)
	defer db_schema.Close()

	schema := structured_schema_from(db_schema)

	augmented_schema := augment_schema_with_virtual(schema, VirtualLink{})

	d2 := schema_to_d2(augmented_schema, false)

	fmt.Println(d2)
}
