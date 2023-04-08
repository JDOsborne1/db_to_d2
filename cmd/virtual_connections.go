package main

type VirtualLink struct {
	source_table      string
	source_column     string
	referenced_table  string
	referenced_column string
}

func augment_schema(_input Schema, _links []VirtualLink) Schema {
	for _, link := range _links {
		_input = augment_tables(_input, link)
	}
	return _input
}

func augment_tables(_input Schema, _link VirtualLink) Schema {
	new_tables := []Table{}
	for _, table := range _input.Tables {
		if table.Name == _link.source_table {
			table = augment_columns(table, _link)
		}
		new_tables = append(new_tables, table)
	}
	_input.Tables = new_tables
	return _input
}

func augment_columns(_table Table, _links VirtualLink) Table {
	new_columns := []Column{}
	for _, column := range _table.Columns {
		if column.Name == _links.source_column {
			column.Reference = &Reference{
				Table:  _links.referenced_table,
				Column: _links.referenced_column,
			}
			column.Key = "VIRTUAL"
			column.Extra = "Virtual link to " + _links.referenced_table + "." + _links.referenced_column
		}
		new_columns = append(new_columns, column)
	}
	_table.Columns = new_columns
	return _table
}
