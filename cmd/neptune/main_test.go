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

func (f *MockSparqlWriter) Post(query string, contentType string) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

type MockSparqlWriterWithError struct {
	mock.Mock
}

func (f *MockSparqlWriterWithError) Post(query string, contentType string) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{StatusCode: 400}, nil
}

type MockSnsWriter struct {
	mock.Mock
}

var messages []string

func (f *MockSnsWriter) Publish(message string) error {
	messages = append(messages, message)
	return nil
}

func TestHandlerUnit(t *testing.T) {
	is := is.New(t)
	registry := runtime.NewRegistry(new(MockSparqlWriter), new(MockSnsWriter))
	handler := runtime.NewHandler(registry)
	var testCases = []struct {
		file        string
		contentType string
		out         int
		msgCount    int
	}{
		{
			file:        "../../fixtures/select_triples.txt",
			contentType: "application/x-www-form-urlencoded",
			out:         200,
			msgCount:    0, // This test is a SELECT, so no message should be published
		},
		{
			file:        "../../fixtures/select_triples.txt",
			contentType: "application/sparql-query",
			out:         200,
			msgCount:    0,
		},
		{
			file:        "../../fixtures/decoded_query.txt",
			contentType: "application/x-www-form-urlencoded",
			out:         422,
			msgCount:    0, // No messages should be published on a bad query
		},
		{
			file:        "../../fixtures/decoded_query.txt",
			contentType: "application/sparql-update",
			out:         200,
			msgCount:    1,
		},

		{
			file:        "../../fixtures/insert.txt",
			contentType: "application/x-www-form-urlencoded",
			out:         200,
			msgCount:    2, // A message should only be added on a successful INSERT
		},
	}

	for _, tt := range testCases {
		content, _ := ioutil.ReadFile(tt.file)
		actual, err := handler.RequestHandler(nil, events.APIGatewayProxyRequest{Body: string(content), Headers: map[string]string{"Content-Type": tt.contentType}})
		is.NoErr(err)
		is.Equal(tt.out, actual.StatusCode)
		is.Equal(tt.msgCount, len(messages))
	}

}

func TestHandlerWithBadQuery(t *testing.T) {
	is := is.New(t)
	registry := runtime.NewRegistry(new(MockSparqlWriterWithError), new(MockSnsWriter))
	handler := runtime.NewHandler(registry)
	content, _ := ioutil.ReadFile("../../fixtures/bad_insert.txt")
	actual, err := handler.RequestHandler(nil, events.APIGatewayProxyRequest{Body: string(content), Headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"}})
	is.NoErr(err)
	is.Equal(400, actual.StatusCode)
	is.Equal(2, len(messages))
}
