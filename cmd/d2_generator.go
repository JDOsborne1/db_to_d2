package main

import (
	"fmt"
	"strings"
)

func schema_to_d2(schema Schema, _minimalist bool) string {
	var builder strings.Builder

	// Write table definitions
	for _, table := range schema.Tables {
		builder.WriteString(table_to_d2(table, _minimalist))
	}

	// Write foreign key relationships
	for _, table := range schema.Tables {
		for _, column := range table.Columns {
			if column.Reference != nil {
				builder.WriteString(fmt.Sprintf("%s.%s -> %s.%s", table.Name, column.Name, column.Reference.Table, column.Reference.Column))
				builder.WriteString("{target-arrowhead: {shape: cf-many}}")
				builder.WriteString("\n\n")
			}
		}
	}

	return builder.String()
}

func table_to_d2(_table Table, _minimalist bool) string {
	var builder strings.Builder


			builder.WriteString(fmt.Sprintf("%s: {\n  shape: sql_table\n", _table.Name))

		for _, column := range _table.Columns {
			if column.Key == "" && _minimalist {
				continue
			}

			builder.WriteString(fmt.Sprintf("  %s: %s", column.Name, column.Type))

			if column.Key == "PRI" {
				builder.WriteString(" {constraint: primary_key}")
			} else if column.Key == "MUL" {
				builder.WriteString(" {constraint: foreign_key}")
				} else if column.Key == "UNK" {
				builder.WriteString(" {constraint: unique}")

			}

			builder.WriteString("\n")
		}

		builder.WriteString("}\n\n")


		return builder.String()
}