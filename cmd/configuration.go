package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func get_virtual_links() []VirtualLink {
	links := []VirtualLink{}
	links_json, err := os.ReadFile(os.Getenv("VIRTUAL_LINKS_PATH"))
	if err != nil {
		//TODO: Log error, or bubble up instead of printing to console
		fmt.Println("Failed to read virtual links file")
	}
	json.Unmarshal(links_json, &links)
	return links
}

func get_table_groups() []TableGroup {
	table_groups := []TableGroup{}
	table_groups_json, err := os.ReadFile(os.Getenv("TABLE_GROUPS_PATH"))
	if err != nil {
		//TODO: Log error, or bubble up instead of printing to console
		fmt.Println("Failed to read table groups file")
	}
	json.Unmarshal(table_groups_json, &table_groups)
	return table_groups
}