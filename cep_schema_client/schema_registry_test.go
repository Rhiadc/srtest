package cepsr_test

import (
	"testing"

	cepschemaregistry "github.com/rhiadc/srtest/cep_schema_client"
	"github.com/riferrei/srclient"
	"github.com/stretchr/testify/assert"
)

var (
	schemaRegistryURL = "http://localhost:8085"
	wrongSchemaError  = `cannot decode textual record "User": cannot decode textual map: cannot determine codec: "wrong_field"`
)

func TestValidSchemaJsonMessage(t *testing.T) {
	schemaRegistryClient := srclient.CreateSchemaRegistryClient(schemaRegistryURL)
	schema, err := schemaRegistryClient.GetLatestSchema("teste")
	user := map[string]interface{}{"id": 22, "name": "user_teste"}

	sr := cepschemaregistry.NewSchemaRegistryClient(schemaRegistryClient)
	val, err := sr.EncodeMessageIntoAvroAndInsertSchemaID(user, schema)

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 0, 0, 1, 44, 20, 117, 115, 101, 114, 95, 116, 101, 115, 116, 101}, val)
}

func TestInvalidFormats(t *testing.T) {
	schemaRegistryClient := srclient.CreateSchemaRegistryClient(schemaRegistryURL)
	schema, err := schemaRegistryClient.GetLatestSchema("teste")
	user := "teste"

	sr := cepschemaregistry.NewSchemaRegistryClient(schemaRegistryClient)
	val, err := sr.EncodeMessageIntoAvroAndInsertSchemaID(user, schema)

	assert.NotNil(t, err)
	assert.Nil(t, val)
}

func TestInvalidSchemaJsonMessage(t *testing.T) {
	schemaRegistryClient := srclient.CreateSchemaRegistryClient(schemaRegistryURL)
	schema, err := schemaRegistryClient.GetLatestSchema("teste")
	user := map[string]interface{}{"id": 22, "wrong_field": "wrong"}

	sr := cepschemaregistry.NewSchemaRegistryClient(schemaRegistryClient)
	_, err = sr.EncodeMessageIntoAvroAndInsertSchemaID(user, schema)

	assert.Equal(t, wrongSchemaError, err.Error())
}

func TestInvalidTopicForSchemaClient(t *testing.T) {
	schemaRegistryClient := srclient.CreateSchemaRegistryClient(schemaRegistryURL)
	_, err := schemaRegistryClient.GetLatestSchema("test")
	assert.Equal(t, err.Error(), `{"error_code":40401,"message":"Subject 'test' not found."}`)
}

func TestValidTopicForSchemaClient(t *testing.T) {
	schemaRegistryClient := srclient.CreateSchemaRegistryClient(schemaRegistryURL)
	_, err := schemaRegistryClient.GetLatestSchema("teste")
	assert.Nil(t, err)
}
