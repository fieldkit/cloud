package main

type JsonSchemaField struct {
	Name     string
	Type     string
	Optional bool
}

type JsonMessageSchema struct {
	MessageSchema
	HasTime         bool
	UseProviderTime bool
	HasLocation     bool
	Fields          []JsonSchemaField
}
