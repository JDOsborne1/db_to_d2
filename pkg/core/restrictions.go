package core

type Restrictor func(Table, Column) bool

func Minimalist(_table Table, _column Column) bool {
	return _column.Key == ""
}

func standard(_table Table, _column Column) bool {
	return false
}

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
