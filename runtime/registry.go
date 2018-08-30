package runtime

import (
	"github.com/sul-dlss-labs/sparql-loader/message"
	"github.com/sul-dlss-labs/sparql-loader/sparql"
)

// Registry exposes the writer to the configured sparql endpoint
type Registry struct {
	Writer    sparql.Writer
	Publisher message.Publisher
}

// NewRegistry creates a new instance of Registry
func NewRegistry(writer sparql.Writer, publisher message.Publisher) *Registry {
	return &Registry{Writer: writer, Publisher: publisher}
}
