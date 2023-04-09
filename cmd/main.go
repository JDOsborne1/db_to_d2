package main

import (
	"core"
	"fmt"
	"mysql"
	"virtual"
)

type Table core.Table
type Column core.Column
type Reference core.Reference
type Schema core.Schema
type TableGroup core.TableGroup


type options struct {
	use_virtual_links bool
	use_table_groups  bool
	restrictor_type   string
}

func main() {
	// Connect to the MySQL database
	db, err := mysql.Connect_to_db()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db_schema := mysql.Information_schema_from(db)
	defer db_schema.Close()

	schema := mysql.Structured_schema_from(db_schema)

	options := get_options()

	if options.use_virtual_links {
		schema = virtual.Augment_schema(schema, get_virtual_links())
	}

	switch options.restrictor_type {
	case "user":
		permission_restrictor := mysql.Restrict_to_table_for_user(db, get_designated_user())
		schema = core.Restrict(schema, permission_restrictor)
	case "minimal":
		schema = core.Restrict(schema, core.Minimalist)
	}

	var d2 string
	if options.use_table_groups {
		d2 = core.Schema_to_d2(schema, get_table_groups())
	} else {
		d2 = core.Schema_to_d2(schema, []core.TableGroup{})
	}

	fmt.Println(d2)

}
