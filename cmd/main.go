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

func augment_schema(_input Schema, _links []VirtualLink) Schema {
	for _, link := range _links {
		_input = augment_tables(_input, link)
	}
	return _input
}

func augment_tables(_input Schema, _link VirtualLink) Schema {
	new_tables := []Table{}
	for _, table := range _input.Tables {
		if table.Name == _link.source_table {
			table = augment_columns(table, _link)
		}
		new_tables = append(new_tables, table)
	}
	_input.Tables = new_tables
	return _input
}

func augment_columns(_table Table, _links VirtualLink) Table {
	new_columns := []Column{}
	for _, column := range _table.Columns {
		if column.Name == _links.source_column {
			column.Reference = &Reference{
				Table: _links.referenced_table,
				Column: _links.referenced_column,
			}
			column.Key = "VIRTUAL"
		}
		new_columns = append(new_columns, column)
	}
	_table.Columns = new_columns
	return _table
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
