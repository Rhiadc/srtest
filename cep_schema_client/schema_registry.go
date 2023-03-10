package cepsr

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"

	"github.com/riferrei/srclient"
)

const schemaDef = `{
	"type": "record",
	"name": "User",
	"fields": [
		{"name": "id", "type": "int"},
		{"name": "name", "type": "string"}
	]
}`

type SchemaRegistry struct {
	Client *srclient.SchemaRegistryClient
}

func NewSchemaRegistryClient(client *srclient.SchemaRegistryClient) *SchemaRegistry {
	return &SchemaRegistry{Client: client}
}

func (sr SchemaRegistry) GetOrCreateSchema(topic string) (*srclient.Schema, error) {
	schema, err := sr.Client.GetLatestSchema(topic)
	log.Fatal(schema)
	if schema == nil {
		schema, err = sr.Client.CreateSchema(topic, schemaDef, srclient.Avro)
		return schema, err
	}
	return schema, err
}

func (sr SchemaRegistry) EncodeMessageIntoAvroAndInsertSchemaID(message interface{}, schema *srclient.Schema) ([]byte, error) {
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))

	value, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	native, _, err := schema.Codec().NativeFromTextual(value)
	if err != nil {
		return nil, err
	}
	valueBytes, err := schema.Codec().BinaryFromNative(nil, native)
	if err != nil {
		return nil, err
	}

	var recordValue []byte
	recordValue = append(recordValue, byte(0))
	recordValue = append(recordValue, schemaIDBytes...)
	recordValue = append(recordValue, valueBytes...)
	fmt.Println(recordValue)
	fmt.Println(schemaIDBytes)
	return recordValue, nil
}
