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
	rows := retrieve_information_schema()
	defer rows.Close()

	schema := generate_schema_from_sql_rows(rows)

	fmt.Println(schema)

	// triples := generate_turtle_from_schema(schema)
	
	// fmt.Println(triples)

	d2 := generateD2Diagram(schema)

	fmt.Println(d2)
}
