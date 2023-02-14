package main 

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
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
	db, err := sql.Open("mysql", "username:password@tcp(host:port)/database")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Retrieve the table and column information from the information schema
	rows, err := db.Query("SELECT TABLE_NAME, COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY, EXTRA, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME, UPDATE_RULE, DELETE_RULE FROM INFORMATION_SCHEMA.COLUMNS LEFT JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE ON COLUMN_NAME = REFERENCED_COLUMN_NAME AND TABLE_NAME = REFERENCED_TABLE_NAME WHERE TABLE_SCHEMA = DATABASE() ORDER BY TABLE_NAME, ORDINAL_POSITION")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var schema Schema
	var currentTable string
	var currentColumns []Column

	// Loop through each row and build the data structure
	for rows.Next() {
		var tableName, columnName, columnType, isNullable, columnKey, extra, referencedTable, referencedColumn, onUpdate, onDelete sql.NullString
		err := rows.Scan(&tableName, &columnName, &columnType, &isNullable, &columnKey, &extra, &referencedTable, &referencedColumn, &onUpdate, &onDelete)
		if err != nil {
			panic(err.Error())
		}

		// If the table name has changed, add the current table to the schema and start a new one
		if tableName.String != currentTable {
			if len(currentColumns) > 0 {
				schema.Tables = append(schema.Tables, Table{Name: currentTable, Columns: currentColumns})
			}
			currentTable = tableName.String
			currentColumns = []Column{}
		}

		// Create a new column and add it to the current table
		column := Column{
			Name:     columnName.String,
			Type:     columnType.String,
			Nullable: isNullable.String == "YES",
			Key:      columnKey.String,
			Extra:    extra.String,
		}

		if referencedTable.Valid && referencedColumn.Valid {
			column.Reference = &Reference{
				Table:    referencedTable.String,
				Column:   referencedColumn.String,
				OnUpdate: onUpdate.String,
				OnDelete: onDelete.String,
			}
		}

		currentColumns = append(currentColumns, column)
	}

	// Add the last table to the schema
	if len(currentColumns) > 0 {
		schema.Tables = append(schema.Tables, Table{Name: currentTable, Columns: currentColumns})
	}

	fmt.Println(schema)
}
