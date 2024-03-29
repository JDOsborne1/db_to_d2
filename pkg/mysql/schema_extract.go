package mysql

import (
	"core"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

// TODO - Merge this definition with the one in Configuration (Likely needs a config module)

const (
	env_user     = "D2_TARGET_DB_USER"
	env_password = "D2_TARGET_DB_PASSWORD"
	env_host     = "D2_TARGET_DB_HOST"
	env_port     = "D2_TARGET_DB_PORT"
	env_name     = "D2_TARGET_DB_NAME"
)

// connect_to_db connects to the database specified by the environment variables
// - D2_TARGET_DB_USER,
// - D2_TARGET_DB_PASSWORD,
// - D2_TARGET_DB_HOST,
// - D2_TARGET_DB_PORT,
// - D2_TARGET_DB_NAME
// The database must be a MySQL database.
// Returns a pointer to the database and an error if the connection failed.
func connect_to_db() (*sql.DB, error) {
	user := viper.GetString(env_user)
	password := viper.GetString(env_password)
	host := viper.GetString(env_host)
	port := viper.GetString(env_port)
	dbname := viper.GetString(env_name)
	essential_vars := []string{user, password, host, port, dbname}
	for _, v := range essential_vars {
		if v == "" {
			return nil, fmt.Errorf("missing environment variable:" + v)
		}
	}

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

// This function is used to extract the schema from the database.
// it uses the information_schema to get the table and column information.
// Currently it only supports the schema 'testdb'. This will be changed in the future.
// TODO: Make this function consider the schema optional
func information_schema_from(_db *sql.DB, _schema string) *sql.Rows {

	// Retrieve the table and column information from the information schema
	rows, err := _db.Query(`
	SELECT C.TABLE_NAME, C.COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY, EXTRA, KC.TABLE_NAME, KC.COLUMN_NAME 
	FROM INFORMATION_SCHEMA.COLUMNS C
	LEFT JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE KC
	ON C.COLUMN_NAME = REFERENCED_COLUMN_NAME AND C.TABLE_NAME = REFERENCED_TABLE_NAME 
	WHERE C.TABLE_SCHEMA = '` + _schema + `'
	ORDER BY C.TABLE_NAME, KC.ORDINAL_POSITION;
	`)
	if err != nil {
		panic(err.Error())
	}

	return rows
}

// This function is used to extract the schema from the database.
// Once connected to the information schema of the database, it loops through each row and builds the data structure.
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

// Extract_schema retrieves the schema from the database and returns it as a data structure
// Currently this panics if there is an error, but it is intended to eventually return a tuple of (schema, error)
// TODO: Return a tuple of (schema, error)
func Extract_schema() core.Schema {
	// Connect to the database
	db, err := connect_to_db()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Retrieve the schema from the database
	rows := information_schema_from(db, viper.GetString(env_name))

	// Build the data structure
	schema := structured_schema_from(rows)

	return schema
}
