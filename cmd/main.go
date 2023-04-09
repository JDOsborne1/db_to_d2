package main

import (
	"core"
	"fmt"
	"mysql"
	"virtual"
)

// options is a struct that contains the options for the program.
// TODO: Add validation for options
type options struct {
	use_virtual_links bool
	use_table_groups  bool
	restrictor_type   string // "user" or "minimal"
	db_source_type    string // "mysql"
}

func main() {

	options := get_options()

	var schema core.Schema

	switch options.db_source_type {
	case "mysql":
		schema = mysql.Extract_schema()
	default:
		fmt.Println("Invalid db_source_type")
		return
	}

	if options.use_virtual_links {
		schema = virtual.Augment_schema(schema, get_virtual_links())
	}

	switch options.restrictor_type {
	case "user":
		permission_restrictor := mysql.Restrict_to_table_for_user(get_designated_user())
		schema = core.Restrict(schema, permission_restrictor)
	case "minimal":
		schema = core.Restrict(schema, core.Minimalist)
	default:
		fmt.Println("Invalid restrictor_type, using default")
	}

	var d2 string
	if options.use_table_groups {
		d2 = core.Schema_to_d2(schema, get_table_groups())
	} else {
		d2 = core.Schema_to_d2(schema, []core.TableGroup{})
	}

	fmt.Println(d2)

}
