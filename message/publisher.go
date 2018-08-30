package message

// Publisher sends an HTTP Post Request to a SPARQL endpoint
type Publisher interface {
	Publish(message string) error
}
