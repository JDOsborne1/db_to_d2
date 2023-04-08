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

type TableGroup struct {
	Tag    string
	Name   string
	Tables []string
}

func in_set(_element string, _set []string) bool {
	for _, element := range _set {
		if _element == element {
			return true
		}
	}
	return false
}

func wrap_name_in_group(_table_name string, _grouping []TableGroup) string {
	for _, group := range _grouping {
		if in_set(_table_name, group.Tables) {
			return group.Tag + "." + _table_name
		}
	}
	return _table_name
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
		source_table:      "comments",
		source_column:     "content",
		referenced_table:  "posts",
		referenced_column: "content",
	})

	table_group1 := TableGroup{
		Tag:    "ugc",
		Tables: []string{"comments", "posts"},
		Name:   "User Generated Content",
	}

	table_group2 := TableGroup{
		Tag:    "pii",
		Tables: []string{"users"},
		Name:   "Personally Identifiable Information",
	}

	augmented_schema := augment_schema(schema, links)
	d2 := schema_to_d2(augmented_schema, false, []TableGroup{table_group1, table_group2})

	fmt.Println(d2)
}
