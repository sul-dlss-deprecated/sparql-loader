package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sul-dlss-labs/sparql-loader/message"
	"github.com/sul-dlss-labs/sparql-loader/runtime"
	"github.com/sul-dlss-labs/sparql-loader/sparql"
)

func main() {
	// Establish the clients and register the Lambda handler
	neptuneClient := sparql.NewNeptuneClient(os.Getenv("RIALTO_SPARQL_ENDPOINT"))
	snsClient := message.NewClient(os.Getenv("RIALTO_SNS_ENDPOINT"),
		os.Getenv("RIALTO_TOPIC_ARN"),
		os.Getenv("AWS_REGION"))
	registry := runtime.NewRegistry(neptuneClient, snsClient)
	handler := runtime.NewHandler(registry)

	lambda.Start(handler.RequestHandler)
}
