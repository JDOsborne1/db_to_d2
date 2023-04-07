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
	db, err := connect_to_db()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db_schema := information_schema_from(db)
	defer db_schema.Close()

	schema := structured_schema_from(db_schema)

	var links []VirtualLink
	links = append(links, VirtualLink{
		source_table: "comments",
		source_column: "content",
		referenced_table: "posts",
		referenced_column: "content",
	})

	augmented_schema := augment_schema(schema, links)
	d2 := schema_to_d2(augmented_schema, true)

	fmt.Println(d2)
}
