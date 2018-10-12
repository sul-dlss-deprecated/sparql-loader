package test

import (
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/matryer/is"
	"github.com/sul-dlss-labs/sparql-loader/message"
	"github.com/sul-dlss-labs/sparql-loader/runtime"
	"github.com/sul-dlss-labs/sparql-loader/sparql"
)

func TestHandlerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	is := is.New(t)
	sparqlClient := sparql.NewNeptuneClient("http://localhost:9999/blazegraph/namespace/kb/sparql")
	snsClient := message.NewClient("http://localhost:4575", "arn:aws:sns:us-east-1:123456789012:rialto", "localstack", false)
	registry := runtime.NewRegistry(sparqlClient, snsClient)
	handler := runtime.NewHandler(registry)
	var testCases = []struct {
		file        string
		out         int
		contentType string
	}{
		{
			file:        "../fixtures/select_triples.txt",
			contentType: "application/x-www-form-urlencoded",
			out:         200,
		},
		{
			file:        "../fixtures/decoded_query.txt",
			contentType: "application/x-www-form-urlencoded",
			out:         422,
		},
		{
			file:        "../fixtures/insert.txt",
			contentType: "application/x-www-form-urlencoded",
			out:         200,
		},
		{
			file:        "../fixtures/bad_insert.txt",
			contentType: "application/x-www-form-urlencoded",
			out:         400,
		},
	}

	for _, tt := range testCases {
		content, _ := ioutil.ReadFile(tt.file)
		actual, err := handler.RequestHandler(nil, events.APIGatewayProxyRequest{Body: string(content), Headers: map[string]string{"Content-Type": tt.contentType}})
		is.NoErr(err)
		is.Equal(tt.out, actual.StatusCode)
	}

}
