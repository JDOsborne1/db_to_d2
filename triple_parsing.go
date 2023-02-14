package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func generateERDFromSPARQL(endpointURL string) {
	// Construct the SPARQL query to extract the triples
	query := `SELECT ?subject ?predicate ?object WHERE {?subject ?predicate ?object .}`
	url := fmt.Sprintf("%s?query=%s&format=json", endpointURL, query)

	// Send the HTTP GET request to the SPARQL endpoint
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// Read the response body into a byte slice
	data, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Parse the JSON response into a map
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		panic(err)
	}

	// Initialize the D2 string
	d2 := "@startuml\n"

	// Add nodes for each table
	tables := make(map[string]bool)
	for _, binding := range result["results"].(map[string]interface{})["bindings"].([]interface{}) {
		subject := binding.(map[string]interface{})["subject"].(map[string]interface{})["value"].(string)
		predicate := binding.(map[string]interface{})["predicate"].(map[string]interface{})["value"].(string)
		object := binding.(map[string]interface{})["object"].(map[string]interface{})["value"].(string)

		// Ignore triples with rdf:type predicate
		if predicate == "http://www.w3.org/1999/02/22-rdf-syntax-ns#type" {
			continue
		}

		// Extract the table and column names from the subject and predicate URIs
		subjectParts := strings.Split(subject, "#")
		table := subjectParts[1]
		predicateParts := strings.Split(predicate, "#")
		column := predicateParts[1]

		// Add the table to the D2 string if it hasn't been added already
		if !tables[table] {
			d2 += fmt.Sprintf("class %s {\n", table)
			tables[table] = true
		}

		// Add a node for the column
		d2 += fmt.Sprintf("\t%s\n", column)

		// If the object is a foreign key reference, add a relationship line
		if strings.HasPrefix(object, "http://example.com/database#") {
			referenceParts := strings.Split(object, "#")
			referencedTable := referenceParts[1]
			d2 += fmt.Sprintf("\t%s -- %s\n", column, referencedTable)
		}
	}

	// Close the class definition for each table
	for range tables {
		d2 += "}\n"
	}

	// Add a closing bracket for the D2 string
	d2 += "@enduml\n"

	// Print the D2 string to the console
	fmt.Println(d2)
}
