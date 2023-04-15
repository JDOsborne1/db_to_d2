package main

import (
	"core"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"virtual"
)

// get_virtual_links returns the virtual links for the program.
// These are currently set by a file specified by the VIRTUAL_LINKS_PATH environment variable.
// The file should be a json array of virtual links. See the virtual package for more information.
func get_virtual_links() []virtual.VirtualLink {
	links_reader, err := os.Open(os.Getenv("VIRTUAL_LINKS_PATH"))
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
	table_groups_reader, err := os.Open(os.Getenv("TABLE_GROUPS_PATH"))
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
	return os.Getenv("DESIGNATED_USER")
}

// get_options returns the options for the program. These are currently set by environment variables
func get_options() options {
	//TODO: Add validation for options
	return options{
		use_virtual_links: os.Getenv("VIRTUAL_LINKS") == "true",
		use_table_groups:  os.Getenv("TABLE_GROUPS") == "true",
		restrictor_type:   os.Getenv("RESTRICTOR_TYPE"),
		db_source_type:    os.Getenv("D2_TARGET_DB_TYPE"),
	}
}
