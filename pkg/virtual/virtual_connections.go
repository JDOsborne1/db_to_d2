package virtual

import (
	"core"
)

type VirtualLink struct {
	SourceTable      string `json:"source_table,omitempty"`
	SourceColumn     string `json:"source_column,omitempty"`
	ReferencedTable  string `json:"referenced_table,omitempty"`
	ReferencedColumn string `json:"referenced_column,omitempty"`
}

func Augment_schema(_input core.Schema, _links []VirtualLink) core.Schema {
	for _, link := range _links {
		_input = augment_tables(_input, link)
	}
	return _input
}

func augment_tables(_input core.Schema, _link VirtualLink) core.Schema {
	new_tables := []core.Table{}
	for _, table := range _input.Tables {
		if table.Name == _link.SourceTable {
			table = augment_columns_source(table, _link)
		}
		if table.Name == _link.ReferencedTable {
			table = augment_columns_reference(table, _link)
		}

		new_tables = append(new_tables, table)
	}
	_input.Tables = new_tables
	return _input
}

func augment_columns_source(_table core.Table, _links VirtualLink) core.Table {
	new_columns := []core.Column{}
	for _, column := range _table.Columns {
		if column.Name == _links.SourceColumn {
			column.Reference = &core.Reference{
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

func augment_columns_reference(_table core.Table, _links VirtualLink) core.Table {
	new_columns := []core.Column{}
	for _, column := range _table.Columns {
		if column.Name == _links.ReferencedColumn {
			column.Key = "VIRTUAL"
			column.Extra = "Virtual link from " + _links.SourceTable + "." + _links.SourceColumn
		}
		new_columns = append(new_columns, column)
	}
	_table.Columns = new_columns
	return _table
}
