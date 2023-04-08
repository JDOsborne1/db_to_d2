package main

import (
	"fmt"
	"strings"
)

func schema_to_d2(schema Schema, _minimalist bool, _group TableGroup) string {
	var builder strings.Builder

	ungrouped_tables := []Table{}
	grouped_tables := []Table{}
	// Extracting table groups 
	for _, table := range schema.Tables {
		in_set := false
		for _, group_table := range _group.Tables {
			if table.Name == group_table {
				in_set = true
			}
		}
		if !in_set {
			ungrouped_tables = append(ungrouped_tables, table)
		}
		if in_set {
			grouped_tables = append(grouped_tables, table)
		}
	}


	// Write table definitions
	for _, table := range ungrouped_tables {
		builder.WriteString(table_to_d2(table, _minimalist))
	}

	builder.WriteString(_group.Name  + ": { \n")

	// Write table definitions
	for _, table := range grouped_tables {
		builder.WriteString(table_to_d2(table, _minimalist))
	}

	builder.WriteString("}\n")


	// Write foreign key relationships
	for _, table := range schema.Tables {
		for _, column := range table.Columns {
			if column.Reference != nil {
				builder.WriteString(fmt.Sprintf("%s.%s -> %s.%s", wrap_name_in_group(table.Name, _group), column.Name, wrap_name_in_group(column.Reference.Table, _group), column.Reference.Column))
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

		} else if column.Key == "VIRTUAL" {
			builder.WriteString(" {constraint: foreign_key}")
		}

		builder.WriteString("\n")
	}

	builder.WriteString("}\n\n")

	return builder.String()
}
