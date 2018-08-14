package sparql

// Sparql struct is the return struct for parsing
type Sparql struct {
	Verb  string
	Graph []Triple
	Body  string
}

// Query is a struct to hold a full SPARQL query
type Query struct {
	Parts      []Sparql
	Prefixes   map[string]string
	NamedGraph *string
}

// NewQuery returns an empty Sparql struct to the calling method
func NewQuery() Query {
	return Query{Prefixes: make(map[string]string)}
}
