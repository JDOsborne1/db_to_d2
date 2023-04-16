package main

import (
	"core"
	"fmt"
	"mysql"
	"virtual"

	"github.com/spf13/pflag"
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

func register_commandline_flags() {
	pflag.String("VirtualLinks", "false", "Use virtual links")
	pflag.String("VirtualLinksPath", "example_virtual_links.json", "Path to virtual links file")
	pflag.String("TableGroups", "false", "Use table groups")
	pflag.String("TableGroupsPath", "example_table_groups.json", "Path to table groups file")
	pflag.String("RestrictorType", "", "Restrictor type")
	
	pflag.String("D2TargetDbUser", "", "D2 target db user")
	pflag.String("D2TargetDbPassword", "", "D2 target db password")
	pflag.String("D2TargetDbHost", "", "D2 target db host")
	pflag.String("D2TargetDbPort", "", "D2 target db port")
	pflag.String("D2TargetDbName", "", "D2 target db name")
	pflag.String("D2TargetDbType", "", "D2 target db type")
	pflag.String("DesignatedUser", "", "Designated user")
	
	pflag.Parse()
	viper.RegisterAlias("VIRTUAL_LINKS", "VirtualLinks")
	viper.RegisterAlias("VIRTUAL_LINKS_PATH", "VirtualLinksPath")
	viper.RegisterAlias("TABLE_GROUPS", "TableGroups")
	viper.RegisterAlias("TABLE_GROUPS_PATH", "TableGroupsPath")
	viper.RegisterAlias("RESTRICTOR_TYPE", "RestrictorType")
	
	viper.RegisterAlias("D2_TARGET_DB_USER", "D2TargetDbUser")
	viper.RegisterAlias("D2_TARGET_DB_PASSWORD", "D2TargetDbPassword")
	viper.RegisterAlias("D2_TARGET_DB_HOST", "D2TargetDbHost")
	viper.RegisterAlias("D2_TARGET_DB_PORT", "D2TargetDbPort")
	viper.RegisterAlias("D2_TARGET_DB_NAME", "D2TargetDbName")
	viper.RegisterAlias("D2_TARGET_DB_TYPE", "D2TargetDbType")
	viper.RegisterAlias("DESIGNATED_USER", "DesignatedUser")
	
	viper.BindPFlags(pflag.CommandLine)

}

func register_environent_variables() {
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

}

func main() {

	register_commandline_flags()
	register_environent_variables()


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
