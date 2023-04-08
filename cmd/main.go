package main

import (
	"fmt"
)

// Column represents a single column in a database table.
type Column struct {
	// Name is the name of the column.
	Name string

	// Type is the data type of the column.
	Type string

	// Nullable is a boolean value that indicates whether the column can be null.
	Nullable bool

	// Key is a string that indicates whether the column is a primary key, foreign key, or other type of key.
	Key string

	// Extra is a string that contains any additional information about the column, such as AUTO_INCREMENT or DEFAULT.
	Extra string

	// Reference is a pointer to a Reference struct that represents a foreign key constraint, if any.
	Reference *Reference
}

// Reference represents a foreign key constraint on a column.
type Reference struct {
	// Table is the name of the table that the foreign key references.
	Table string

	// Column is the name of the column that the foreign key references.
	Column string

	// OnDelete is a string that specifies the action to take when a row in the referenced table is deleted.
	OnDelete string

	// OnUpdate is a string that specifies the action to take when a row in the referenced table is updated.
	OnUpdate string
}

// Table represents a database table.
type Table struct {
	// Name is the name of the table.
	Name string

	// Columns is a slice of Column structs that represent the columns in the table.
	Columns []Column
}

// Schema represents a database schema.
type Schema struct {
	// Tables is a slice of Table structs that represent the tables in the schema.
	Tables []Table

	// Indexes is a slice of strings that represent the indexes in the schema.
	Indexes []string
}

// TableGroup represents a group of tables in a database schema.
type TableGroup struct {
	// Tag is a string that identifies the group.
	Tag string

	// Name is the name of the group.
	Name string

	// Tables is a slice of strings that represent the names of the tables in the group.
	Tables []string
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

	designated_user := "'testuser'@'%'"

	table_level_permissions, err := get_table_level_permissions(db)
	if err != nil {
		panic(err)
	}

	column_level_permissions, err := get_column_level_permissions(db)
	if err != nil {
		panic(err)
	}

	permission_restrictor := permission_driven_restrictor(table_level_permissions, column_level_permissions, designated_user)

	augmented_schema := augment_schema(schema, links)
	d2 := schema_to_d2(augmented_schema, permission_restrictor, []TableGroup{table_group1, table_group2})

	fmt.Println(d2)

}
