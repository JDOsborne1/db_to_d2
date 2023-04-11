package core

// Restrictor is a function that takes a table and a column and returns a boolean
// indicating whether the column should be excluded from the schema
// This is used to restrict the schema to a subset of the database in the Restrict function
type Restrictor func(Table, Column) bool

// Minimalist is a Restrictor that excludes all columns that are not keys
// This can be used to restrict the schema to only the tables and columns that are
// necessary to represent the relationships between tables. Often, this is all that is needed
func Minimalist(_table Table, _column Column) bool {
	return _column.Key == ""
}

// Standard is a Restrictor that excludes no columns. This is the default behavior.
func Standard(_table Table, _column Column) bool {
	return false
}

// Restrict takes a schema and a restrictor and returns a new schema that excludes all columns
// for which the restrictor returns true.
func Restrict(_schema Schema, _restrictor Restrictor) Schema {
	new_tables := []Table{}
	for _, table := range _schema.Tables {
		new_columns := []Column{}
		for _, column := range table.Columns {
			if _restrictor(table, column) {
				continue
			}
			new_columns = append(new_columns, column)
		}
		if len(new_columns) != 0 {
			table.Columns = new_columns
			new_tables = append(new_tables, table)
		}
	}
	_schema.Tables = new_tables
	return _schema
}
