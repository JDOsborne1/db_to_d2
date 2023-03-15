package main

import (
	"fmt"
	"strings"
)

func generateD2Diagram(schema Schema) string {
    var builder strings.Builder
    
    // Write table definitions
    for _, table := range schema.Tables {
        builder.WriteString(fmt.Sprintf("%s: {\n  shape: sql_table\n", table.Name))
        
        for _, column := range table.Columns {
            builder.WriteString(fmt.Sprintf("  %s: %s", column.Name, column.Type))
            
            if column.Key == "PRI" {
                builder.WriteString(" {constraint: primary_key}")
            }
                                  
            builder.WriteString("\n")
        }
        
        builder.WriteString("}\n\n")
    }
    
    // Write foreign key relationships
    for _, table := range schema.Tables {
        for _, column := range table.Columns {
            if column.Reference != nil {
                builder.WriteString(fmt.Sprintf("%s.%s -> %s.%s\n\n", table.Name, column.Name, column.Reference.Table, column.Reference.Column))
            }
        }
    }
    
    return builder.String()
}