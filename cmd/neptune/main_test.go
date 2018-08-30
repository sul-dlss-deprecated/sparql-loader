package main

import (
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/matryer/is"
	"github.com/stretchr/testify/mock"
	"github.com/sul-dlss-labs/sparql-loader/runtime"
)

type MockSparqlWriter struct {
	mock.Mock
}

func (f *MockSparqlWriter) Post(query string) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

type MockSnsWriter struct {
	mock.Mock
}

var messages []string

func (f *MockSnsWriter) Publish(message string) error {
	messages = append(messages, message)
	return nil
}

func TestHandler_Unit(t *testing.T) {
	is := is.New(t)
	registry := runtime.NewRegistry(new(MockSparqlWriter), new(MockSnsWriter))
	handler := runtime.NewHandler(registry)
	var testCases = []struct {
		file     string
		out      int
		msgCount int
	}{
		{
			file:     "../../fixtures/example1.txt",
			out:      200,
			msgCount: 0, // This test is a SELECT, so no message should be published
		},
		{
			file:     "../../fixtures/example2.txt",
			out:      422,
			msgCount: 0, // No messages should be published on a bad query
		},
		{
			file:     "../../fixtures/example3.txt",
			out:      200,
			msgCount: 1, // A message should only be added on a successful INSERT
		},
	}

	for _, tt := range testCases {
		content, _ := ioutil.ReadFile(tt.file)
		actual, err := handler.RequestHandler(nil, events.APIGatewayProxyRequest{Body: string(content)})
		is.NoErr(err)
		is.Equal(tt.out, actual.StatusCode)
		is.Equal(tt.msgCount, len(messages))
	}

}
