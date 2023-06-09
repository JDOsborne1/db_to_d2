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

// get_virtual_links returns the virtual links for the program.
// These are currently set by a file specified by the VIRTUAL_LINKS_PATH environment variable.
// The file should be a json array of virtual links. See the virtual package for more information.
func get_virtual_links() []virtual.VirtualLink {
	links_reader, err := os.Open(viper.GetString("VIRTUAL_LINKS_PATH"))
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
	table_groups := []core.TableGroup{}
	table_groups_reader, err := os.Open(viper.GetString("TABLE_GROUPS_PATH"))
	if err != nil {
		//TODO: Log error, or bubble up instead of printing to console
		fmt.Println("Failed to read table groups file")
	}

	table_groups, err = read_table_groups(table_groups_reader)

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
	return viper.GetString("DESIGNATED_USER")
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
		use_virtual_links: viper.GetBool("VIRTUAL_LINKS"),
		use_table_groups:  viper.GetBool("TABLE_GROUPS"),
		restrictor_type:   viper.GetString("RESTRICTOR_TYPE"),
		db_source_type:    viper.GetString("D2_TARGET_DB_TYPE"),
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
		viper.RegisterAlias("VIRTUAL_LINKS", "VirtualLinks")
	}
	if *virtual_links_path != "" {
		viper.RegisterAlias("VIRTUAL_LINKS_PATH", "VirtualLinksPath")
	}
	if *table_groups != "" {
	viper.RegisterAlias("TABLE_GROUPS", "TableGroups")
	}

	if *table_groups_path != "" {
	viper.RegisterAlias("TABLE_GROUPS_PATH", "TableGroupsPath")
	}
	if *restrictor_type != "" {
	viper.RegisterAlias("RESTRICTOR_TYPE", "RestrictorType")
	}
	if *db_user != "" {
	viper.RegisterAlias("D2_TARGET_DB_USER", "D2TargetDbUser")
	}
	if *db_password != "" {
		viper.RegisterAlias("D2_TARGET_DB_PASSWORD", "D2TargetDbPassword")
	}
	if *db_host != "" {
	viper.RegisterAlias("D2_TARGET_DB_HOST", "D2TargetDbHost")
	}

	if *db_port != "" {
	viper.RegisterAlias("D2_TARGET_DB_PORT", "D2TargetDbPort")
	}
	if *db_name != "" {
	viper.RegisterAlias("D2_TARGET_DB_NAME", "D2TargetDbName")
	}
	if *db_type != "" {
	viper.RegisterAlias("D2_TARGET_DB_TYPE", "D2TargetDbType")
	}
	if *db_designated_user != "" {
	viper.RegisterAlias("DESIGNATED_USER", "DesignatedUser")
	}
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