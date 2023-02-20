package main

import (
	"database/sql"
	"fmt"
	// "strings"

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
	db, err := sql.Open("mysql", "root:example_password@tcp(localhost:3306)/testdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()


	// Retrieve the table and column information from the information schema
	rows, err := db.Query(`
	SELECT C.TABLE_NAME, C.COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY, EXTRA, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME 
FROM INFORMATION_SCHEMA.COLUMNS C
LEFT JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE KC
ON C.COLUMN_NAME = REFERENCED_COLUMN_NAME AND C.TABLE_NAME = REFERENCED_TABLE_NAME 
WHERE C.TABLE_SCHEMA = 'testdb'
ORDER BY TABLE_NAME, KC.ORDINAL_POSITION;`)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var schema Schema
	var currentTable string
	var currentColumns []Column

	// Loop through each row and build the data structure
	for rows.Next() {
		var tableName, columnName, columnType, isNullable, columnKey, extra, referencedTable, referencedColumn sql.NullString
		err := rows.Scan(&tableName, &columnName, &columnType, &isNullable, &columnKey, &extra, &referencedTable, &referencedColumn)
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
			}
		}

		currentColumns = append(currentColumns, column)
	}

	// Add the last table to the schema
	if len(currentColumns) > 0 {
		schema.Tables = append(schema.Tables, Table{Name: currentTable, Columns: currentColumns})
	}

	fmt.Println(schema)

	// // Create a map of namespace prefixes and URIs
	// prefixes := map[string]string{
	// 	"rdf": "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
	// 	"rdfs": "http://www.w3.org/2000/01/rdf-schema#",
	// 	"owl": "http://www.w3.org/2002/07/owl#",
	// 	"db": "http://example.com/database#",
	// }

	// // Initialize the triples string
	// triples := ""

	// // Add namespace declarations
	// for prefix, uri := range prefixes {
	// 	triples += fmt.Sprintf("@prefix %s: <%s> .\n", prefix, uri)
	// }

	// // Add triples for each table and column
	// for _, table := range schema.Tables {
	// 	// Add the table as a resource
	// 	triples += fmt.Sprintf("\ndb:%s a owl:Class ;\n", table.Name)

	// 	// Add triples for each column
	// 	for _, column := range table.Columns {
	// 		// Add the column as a resource
	// 		triples += fmt.Sprintf("\tdb:%s a owl:DatatypeProperty ;\n", column.Name)

	// 		// Add triples for the column type
	// 		triples += fmt.Sprintf("\t\trdfs:range \"%s\"^^xsd:string ;\n", column.Type)

	// 		// Add triples for the nullable flag
	// 		if column.Nullable {
	// 			triples += "\t\trdfs:subClassOf [ a owl:Restriction ; owl:onProperty owl:maxCardinality ; owl:cardinality \"0\"^^xsd:nonNegativeInteger ] ;\n"
	// 		}

	// 		// Add triples for the primary key
	// 		if column.Key == "PRI" {
	// 			triples += "\t\trdfs:subClassOf [ a owl:Restriction ; owl:onProperty owl:maxCardinality ; owl:cardinality \"1\"^^xsd:nonNegativeInteger ] ;\n"
	// 		}

	// 		// Add triples for the foreign key reference
	// 		if column.Reference != nil {
	// 			triples += fmt.Sprintf("\t\trdfs:range db:%s ;\n", column.Reference.Table)
	// 			triples += fmt.Sprintf("\t\tdb:%s rdfs:subPropertyOf [ a owl:ObjectProperty ; owl:inverseOf db:%s ] ;\n", column.Name, column.Reference.Column)
	// 		}

	// 		// Add triples for the extra flag
	// 		if column.Extra != "" {
	// 			triples += fmt.Sprintf("\t\tdb:%s rdfs:comment \"%s\"^^xsd:string ;\n", column.Name, column.Extra)
	// 		}

    //         // Add a semicolon after each column
    //         triples += ";\n"
	// 	}

	// 	// Remove the trailing semicolon after the last column
	// 	triples = strings.TrimSuffix(triples, ";\n")

	// 	// Add a closing bracket for the table
	// 	triples += ".\n"
	// }

	// // Print the triples to the console
	// fmt.Println(triples)
}
