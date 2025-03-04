package Ollama

type Setup struct {
	Model string
	Url   string
}

type RequestPayload struct {
	Model    string
	Messages []map[string]Message
	Stream   bool
	Format   Schema
}

type Message struct {
	Role    string
	Content string
}

type Schema struct {
	Type       SchemaType
	Properties map[string]interface{}
	Required   []string
}

type SchemaType string

const (
	Object SchemaType = "string"
)
