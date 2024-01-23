package main

import (
	"core"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"virtual"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	env_user             = "D2_TARGET_DB_USER"
	env_password         = "D2_TARGET_DB_PASSWORD"
	env_host             = "D2_TARGET_DB_HOST"
	env_port             = "D2_TARGET_DB_PORT"
	env_name             = "D2_TARGET_DB_NAME"
	env_type             = "D2_TARGET_DB_TYPE"
	env_links            = "VIRTUAL_LINKS"
	env_links_path       = "VIRTUAL_LINKS_PATH"
	env_groups           = "TABLE_GROUPS"
	env_groups_path      = "TABLE_GROUPS_PATH"
	env_restriction      = "RESTRICTOR_TYPE"
	env_perspective_user = "DESIGNATED_USER"
)

// get_virtual_links returns the virtual links for the program.
// These are currently set by a file specified by the VIRTUAL_LINKS_PATH environment variable.
// The file should be a json array of virtual links. See the virtual package for more information.
func get_virtual_links() []virtual.VirtualLink {
	links_reader, err := os.Open(viper.GetString(env_links_path))
	if err != nil {
		//TODO: Log error, or bubble up instead of printing to console
		fmt.Println("Failed to open virtual links file")
	}
	links, err := read_virtual_links(links_reader)
	if err != nil {
		fmt.Println("Failed to read virtual links file")
	}

	return links
}

func read_virtual_links(_input io.Reader) ([]virtual.VirtualLink, error) {
	links := []virtual.VirtualLink{}
	links_json, err := io.ReadAll(_input)
	if err != nil {
		return links, err
	}
	err = json.Unmarshal(links_json, &links)
	return links, err
}

// get_table_groups returns the table groups for the program.
// These are currently set by a file specified by the TABLE_GROUPS_PATH environment variable.
// The file should be a json array of table groups. See the core package for more information.
func get_table_groups() []core.TableGroup {
	table_groups_reader, err := os.Open(viper.GetString(env_groups_path))
	if err != nil {
		//TODO: Log error, or bubble up instead of printing to console
		fmt.Println("Failed to read table groups file")
	}

	table_groups, err := read_table_groups(table_groups_reader)

	if err != nil {
		fmt.Println("Failed to read table groups file")
	}
	return table_groups
}

func read_table_groups(_input io.Reader) ([]core.TableGroup, error) {
	table_groups := []core.TableGroup{}
	table_groups_json, err := io.ReadAll(_input)
	if err != nil {
		return table_groups, err
	}
	err = json.Unmarshal(table_groups_json, &table_groups)
	return table_groups, err
}

// get_designated_user returns the designated user for the program. Set by the DESIGNATED_USER environment variable.
// This is used to restrict the schema to the tables that the designated user has access to. See the mysql package for more information.
func get_designated_user() string {
	return viper.GetString(env_perspective_user)
}

// options is a struct that contains the options for the program.
// TODO: Add validation for options
type options struct {
	use_virtual_links bool
	use_table_groups  bool
	restrictor_type   string // "user" or "minimal"
	db_source_type    string // "mysql"
}

// get_options returns the options for the program. These are set by calling the relevant viper methods.
// A notable consequence of this is that the variables are called at time of use, rather than at time of definition.
// This shouldn't matter in the case of this program, but it's worth noting for future use.
// Additionally, the options are set by the following methods, in order of precedence:
//
// 1. Command line flags
//
// 2. Environment variables
//
// 3. Default values
func get_options() options {

	register_commandline_flags()
	register_environent_variables()
	//TODO: Add validation for options
	return options{
		use_virtual_links: viper.GetBool(env_links),
		use_table_groups:  viper.GetBool(env_groups),
		restrictor_type:   viper.GetString(env_restriction),
		db_source_type:    viper.GetString(env_type),
	}
}

func register_commandline_flags() {
	virtual_links := pflag.String("VirtualLinks", "", "Use virtual links (true/false)")
	virtual_links_path := pflag.String("VirtualLinksPath", "", "Path to virtual links file")
	table_groups := pflag.String("TableGroups", "", "Use table groups (true/false)")
	table_groups_path := pflag.String("TableGroupsPath", "", "Path to table groups file")
	restrictor_type := pflag.String("RestrictorType", "", "Restrictor type (minimal/user/none)")

	db_user := pflag.String("D2TargetDbUser", "", "db login user ")
	db_password := pflag.String("D2TargetDbPassword", "", "db login password")
	db_host := pflag.String("D2TargetDbHost", "", "db login host")
	db_port := pflag.String("D2TargetDbPort", "", "db login port")
	db_name := pflag.String("D2TargetDbName", "", "db login name")
	db_type := pflag.String("D2TargetDbType", "", "db login type")
	db_designated_user := pflag.String("DesignatedUser", "", "User to investigate, format: 'username'@'hostname'")

	pflag.Parse()
	if *virtual_links != "" {
		viper.RegisterAlias(env_links, "VirtualLinks")
	}
	if *virtual_links_path != "" {
		viper.RegisterAlias(env_links_path, "VirtualLinksPath")
	}
	if *table_groups != "" {
		viper.RegisterAlias(env_groups, "TableGroups")
	}

	if *table_groups_path != "" {
		viper.RegisterAlias(env_groups_path, "TableGroupsPath")
	}
	if *restrictor_type != "" {
		viper.RegisterAlias(env_restriction, "RestrictorType")
	}
	if *db_user != "" {
		viper.RegisterAlias(env_user, "D2TargetDbUser")
	}
	if *db_password != "" {
		viper.RegisterAlias(env_password, "D2TargetDbPassword")
	}
	if *db_host != "" {
		viper.RegisterAlias(env_host, "D2TargetDbHost")
	}

	if *db_port != "" {
		viper.RegisterAlias(env_port, "D2TargetDbPort")
	}
	if *db_name != "" {
		viper.RegisterAlias(env_name, "D2TargetDbName")
	}
	if *db_type != "" {
		viper.RegisterAlias(env_type, "D2TargetDbType")
	}
	if *db_designated_user != "" {
		viper.RegisterAlias(env_perspective_user, "DesignatedUser")
	}
	viper.BindPFlags(pflag.CommandLine)
}

func register_environent_variables() {
	viper.BindEnv(env_user)
	viper.BindEnv(env_password)
	viper.BindEnv(env_host)
	viper.BindEnv(env_port)
	viper.BindEnv(env_name)
	viper.BindEnv(env_type)

	viper.BindEnv(env_links)
	viper.BindEnv(env_links_path)

	viper.BindEnv(env_groups)
	viper.BindEnv(env_groups_path)

	viper.BindEnv(env_restriction)

	viper.BindEnv(env_perspective_user)

}
