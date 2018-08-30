package sparql

import "github.com/aws/aws-lambda-go/events"

// Writer sends an HTTP Post Request to a SPARQL endpoint
type Writer interface {
	Post(query string) (*events.APIGatewayProxyResponse, error)
}
