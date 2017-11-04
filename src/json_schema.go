package main

type JsonSchemaField struct {
	Name     string
	Type     string
	Optional bool
}

type JsonMessageSchema struct {
	MessageSchema
	UseProviderTime     bool
	UseProviderLocation bool
	HasTime             bool
	HasLocation         bool
	Fields              []JsonSchemaField
}
