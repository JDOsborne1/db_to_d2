package main

type VirtualLink struct {
	SourceTable      string `json:"source_table,omitempty"`
	SourceColumn     string `json:"source_column,omitempty"`
	ReferencedTable  string `json:"referenced_table,omitempty"`
	ReferencedColumn string `json:"referenced_column,omitempty"`
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
		if table.Name == _link.SourceTable {
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
		if column.Name == _links.SourceColumn {
			column.Reference = &Reference{
				Table:  _links.ReferencedTable,
				Column: _links.ReferencedColumn,
			}
			column.Key = "VIRTUAL"
			column.Extra = "Virtual link to " + _links.ReferencedTable + "." + _links.ReferencedColumn
		}
		new_columns = append(new_columns, column)
	}
	_table.Columns = new_columns
	return _table
}
