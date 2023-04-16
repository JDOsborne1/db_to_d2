package main

import (
	"core"
	"fmt"
	"mysql"
	"virtual"

	"github.com/spf13/viper"
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

	viper.BindEnv("D2_TARGET_DB_USER", "D2_TARGET_DB_USER")
	viper.BindEnv("D2_TARGET_DB_PASSWORD", "D2_TARGET_DB_PASSWORD")
	viper.BindEnv("D2_TARGET_DB_HOST", "D2_TARGET_DB_HOST")
	viper.BindEnv("D2_TARGET_DB_PORT", "D2_TARGET_DB_PORT")
	viper.BindEnv("D2_TARGET_DB_NAME", "D2_TARGET_DB_NAME")
	viper.BindEnv("D2_TARGET_DB_TYPE", "D2_TARGET_DB_TYPE")

	viper.BindEnv("VIRTUAL_LINKS", "VIRTUAL_LINKS")
	viper.BindEnv("VIRTUAL_LINKS_PATH", "VIRTUAL_LINKS_PATH")

	viper.BindEnv("TABLE_GROUPS", "TABLE_GROUPS")
	viper.BindEnv("TABLE_GROUPS_PATH", "TABLE_GROUPS_PATH")

	viper.BindEnv("RESTRICTOR_TYPE", "RESTRICTOR_TYPE")

	viper.BindEnv("DESIGNATED_USER", "DESIGNATED_USER")


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
		permission_restrictor, err := mysql.Restrict_to_table_for_user(get_designated_user())
		if err != nil {
			fmt.Println("Error getting permission_restrictor, using default")
			break
		}
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
