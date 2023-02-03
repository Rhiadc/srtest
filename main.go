package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	cepschemaregistry "github.com/rhiadc/srtest/cep_schema_client"
	"github.com/riferrei/srclient"
)

type ComplexType struct {
	ID   int    `json:"id"`
	Age  int    `json:"age"`
	Name string `json:"name"`
}

// Define the Avro schema for the JSON data

// Define the HTTP handler function to handle the POST request
func handlePost(w http.ResponseWriter, r *http.Request) {
	topic := "teste"
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON data from the request body
	var user map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error parsing JSON:", err)
		user = map[string]interface{}{"ID": 1, "Name": "Gopher"}
	}

	src := srclient.CreateSchemaRegistryClient("http://localhost:8085")
	schemaRegistry := cepschemaregistry.NewSchemaRegistryClient(src)

	schema, err := schemaRegistry.GetOrCreateSchema(topic)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error creating the schema:", err)
		return
	}

	recordValue, err := schemaRegistry.EncodeMessageIntoAvroAndInsertSchemaID(user, schema)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error encoding message:", err)
		return
	}

	// Send the Avro data back as the response
	w.WriteHeader(http.StatusOK)
	w.Write(recordValue)
}

func main() {
	// Define the HTTP endpoint for handling POST requests
	http.HandleFunc("/", handlePost)

	// Start the HTTP server
	err := http.ListenAndServe(":8080", nil)
	log.Println("Server running on port 8080...")
	if err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}
}
