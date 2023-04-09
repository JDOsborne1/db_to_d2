package mysql

import (
	"core"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func connect_to_db() (*sql.DB, error) {
	user := os.Getenv("D2_TARGET_DB_USER")
	password := os.Getenv("D2_TARGET_DB_PASSWORD")
	host := os.Getenv("D2_TARGET_DB_HOST")
	port := os.Getenv("D2_TARGET_DB_PORT")
	dbname := os.Getenv("D2_TARGET_DB_NAME")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("Error using dsn:", dataSourceName)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Error using dsn:", dataSourceName)
		return nil, err
	}

	return db, nil
}

func information_schema_from(_db *sql.DB) *sql.Rows {

	// Retrieve the table and column information from the information schema
	rows, err := _db.Query(`
	SELECT C.TABLE_NAME, C.COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY, EXTRA, KC.TABLE_NAME, KC.COLUMN_NAME 
	FROM INFORMATION_SCHEMA.COLUMNS C
	LEFT JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE KC
	ON C.COLUMN_NAME = REFERENCED_COLUMN_NAME AND C.TABLE_NAME = REFERENCED_TABLE_NAME 
	WHERE C.TABLE_SCHEMA = 'testdb'
	ORDER BY C.TABLE_NAME, KC.ORDINAL_POSITION;
	`)
	if err != nil {
		panic(err.Error())
	}

	return rows
}

func structured_schema_from(_rows *sql.Rows) core.Schema {
	var schema core.Schema
	var currentTable string
	var currentColumns []core.Column

	// Loop through each row and build the data structure
	for _rows.Next() {
		var tableName, columnName, columnType, isNullable, columnKey, extra, referencedTable, referencedColumn sql.NullString
		err := _rows.Scan(&tableName, &columnName, &columnType, &isNullable, &columnKey, &extra, &referencedTable, &referencedColumn)
		if err != nil {
			panic(err.Error())
		}
		// If the table name has changed, add the current table to the schema and start a new one
		if tableName.String != currentTable {
			if len(currentColumns) > 0 {
				schema.Tables = append(schema.Tables, core.Table{Name: currentTable, Columns: currentColumns})
			}
			currentTable = tableName.String
			currentColumns = []core.Column{}
		}

		// Create a new column and add it to the current table
		column := core.Column{
			Name:     columnName.String,
			Type:     columnType.String,
			Nullable: isNullable.String == "YES",
			Key:      columnKey.String,
			Extra:    extra.String,
		}

		if referencedTable.Valid && referencedColumn.Valid {
			column.Reference = &core.Reference{
				Table:  referencedTable.String,
				Column: referencedColumn.String,
			}
		}

		currentColumns = append(currentColumns, column)
	}

	// Add the last table to the schema
	if len(currentColumns) > 0 {
		schema.Tables = append(schema.Tables, core.Table{Name: currentTable, Columns: currentColumns})
	}

	return schema
}

func Extract_schema() core.Schema {
	// Connect to the database
	db, err := connect_to_db()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Retrieve the schema from the database
	rows := information_schema_from(db)

	// Build the data structure
	schema := structured_schema_from(rows)

	return schema
}
