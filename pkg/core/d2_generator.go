package core

import (
	"fmt"
	"strings"
)

// Converts a schema to a string containing corresponding d2 definitions.
func Schema_to_d2(schema Schema, _groups []TableGroup) string {
	var builder strings.Builder
	groupings := make(map[string][]Table)
	table_group_check := make(map[string]bool)
	// Extracting table groups
	for _, table := range schema.Tables {
		for _, group := range _groups {
			in_set := in_set(table.Name, group.Tables)
			if in_set {
				groupings[group.Tag] = append(groupings[group.Tag], table)
				table_group_check[table.Name] = true
			}
		}
	}

	// Write ungrouped table definitions
	for _, table := range schema.Tables {
		if !table_group_check[table.Name] {
			builder.WriteString(table_to_d2(table))
		}
	}

	// Write grouped table definitions
	for group, tables := range groupings {

		builder.WriteString(group + ": { \n")

		for _, table := range tables {
			builder.WriteString(table_to_d2(table))
		}
		builder.WriteString("}\n")

	}

	// Write labels for groups
	for _, group := range _groups {
		if group.Name != "" {
			builder.WriteString(group.Tag + " : " + "\"" + group.Name + "\"" + "\n")
		}
	}

	// Write foreign key relationships
	for _, table := range schema.Tables {
		for _, column := range table.Columns {
			if column.Reference != nil {
				// This is a vile hack which has the unfortunate trait of making the program work.
				// Only this arrangement allows for multiple virtual connections from a common reference
				// point without inverting all the original foreign keys from the information schema.
				// This should be rectified in a re-write of the internals of the mysql or virtual package
				if column.Key == "VIRTUAL" {
					builder.WriteString(fmt.Sprintf(
						"%s.%s -> %s.%s",
						wrap_name_in_group(column.Reference.Table, _groups),
						column.Reference.Column,
						wrap_name_in_group(table.Name, _groups),
						column.Name))
				} else {
					builder.WriteString(fmt.Sprintf(
						"%s.%s -> %s.%s",
						wrap_name_in_group(table.Name, _groups),
						column.Name,
						wrap_name_in_group(column.Reference.Table, _groups),
						column.Reference.Column))
				}
				builder.WriteString(" {\n")
				builder.WriteString("  target-arrowhead: {shape: cf-many}\n")
				if column.Key == "VIRTUAL" {
					builder.WriteString("  style: {stroke-dash: 3}")
				}
				builder.WriteString("}")
				builder.WriteString("\n\n")
			}
		}
	}

	return builder.String()
}

// Creates the d2 definition for a single table
func table_to_d2(_table Table) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s: {\n  shape: sql_table\n", _table.Name))

	for _, column := range _table.Columns {
		builder.WriteString(fmt.Sprintf("  %s: %s", column.Name, column.Type))

		switch column.Key {
		case "PRI":
			builder.WriteString(" {constraint: primary_key}")
		case "MUL":
			builder.WriteString(" {constraint: foreign_key}")
		case "UNK":
			builder.WriteString(" {constraint: unique}")
		case "VIRTUAL":
			builder.WriteString(" {constraint: foreign_key}")
		}

		builder.WriteString("\n")
	}

	builder.WriteString("}\n\n")

	return builder.String()
}
