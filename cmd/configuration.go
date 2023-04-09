package main

import (
	"core"
	"encoding/json"
	"fmt"
	"os"
	"virtual"
)

func get_virtual_links() []virtual.VirtualLink {
	links := []virtual.VirtualLink{}
	links_json, err := os.ReadFile(os.Getenv("VIRTUAL_LINKS_PATH"))
	if err != nil {
		//TODO: Log error, or bubble up instead of printing to console
		fmt.Println("Failed to read virtual links file")
	}
	json.Unmarshal(links_json, &links)
	return links
}

func get_table_groups() []core.TableGroup {
	table_groups := []core.TableGroup{}
	table_groups_json, err := os.ReadFile(os.Getenv("TABLE_GROUPS_PATH"))
	if err != nil {
		//TODO: Log error, or bubble up instead of printing to console
		fmt.Println("Failed to read table groups file")
	}
	json.Unmarshal(table_groups_json, &table_groups)
	return table_groups
}

func get_designated_user() string {
	return os.Getenv("DESIGNATED_USER")
}

func get_options() options {
	return options{
		use_virtual_links: os.Getenv("USE_VIRTUAL_LINKS") == "true",
		use_table_groups:  os.Getenv("USE_TABLE_GROUPS") == "true",
		restrictor_type:   os.Getenv("RESTRICTOR_TYPE"),
	}
}
