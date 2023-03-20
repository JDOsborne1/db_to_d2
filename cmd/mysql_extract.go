package main

import (
	"fmt"
	"os"
	
	"database/sql"
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

func structured_schema_from(_rows *sql.Rows) Schema {
	var schema Schema
	var currentTable string
	var currentColumns []Column

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
				Table:  referencedTable.String,
				Column: referencedColumn.String,
			}
		}

		currentColumns = append(currentColumns, column)
	}

	// Add the last table to the schema
	if len(currentColumns) > 0 {
		schema.Tables = append(schema.Tables, Table{Name: currentTable, Columns: currentColumns})
	}

	return schema
}


type UserPermission struct {
	User   string
	Table  string
	Column string
	Select bool
}

func GetUserPermissions(db *sql.DB) ([]UserPermission, error) {
	var permissions []UserPermission

	rows, err := db.Query(`
		SELECT
			Grantee,
			Table_Name,
			Column_Name,
			(SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMN_PRIVILEGES
				WHERE TABLE_NAME = C.Table_Name AND COLUMN_NAME = C.Column_Name AND PRIVILEGE_TYPE = 'SELECT'
				AND GRANTEE = CP.GRANTEE) AS Has_Select_Permission
		FROM
			INFORMATION_SCHEMA.COLUMNS C
			INNER JOIN INFORMATION_SCHEMA.TABLES T ON C.Table_Name = T.Table_Name
			LEFT JOIN INFORMATION_SCHEMA.COLUMN_PRIVILEGES CP ON C.Table_Name = CP.Table_Name AND C.Column_Name = CP.Column_Name
		WHERE
			T.TABLE_TYPE = 'BASE TABLE'
		ORDER BY
			Grantee, Table_Name, Column_Name;
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user string
		var table string
		var column string
		var hasSelectPermission sql.NullInt64

		err = rows.Scan(&user, &table, &column, &hasSelectPermission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, UserPermission{
			User:   user,
			Table:  table,
			Column: column,
			Select: hasSelectPermission.Valid && hasSelectPermission.Int64 > 0,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return permissions, nil
}
