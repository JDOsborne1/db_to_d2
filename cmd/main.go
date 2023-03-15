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

func main() {
	// Connect to the MySQL database
	db := connect_to_db()
	defer db.Close()
	db_schema := information_schema_from(db)
	defer db_schema.Close()

	schema := structured_schema_from(db_schema)

	d2 := schema_to_d2(schema)

	fmt.Println(d2)
}
