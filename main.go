package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/sul-dlss-labs/sparql"
)

// Handler is the Lambda function handler
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (string, error) {
	proxyReq, err := http.NewRequest("POST", os.Getenv("RIALTO_SPARQL_ENDPOINT"), strings.NewReader(request.Body))
	proxyReq.Header = make(http.Header)

	proxyReq.Header.Add("Content-type", "application/x-www-form-urlencoded")
	proxyReq.Header.Set("Content-Length", strconv.Itoa(len(request.Body)))

	httpClient := http.Client{}
	resp, err := httpClient.Do(proxyReq)
	respBody, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return string(respBody), err
	}

	if strings.HasPrefix(request.Body, "update=") {
		sparqlQuery := sparql.NewSparql()
		err = sparqlQuery.Parse(strings.NewReader(request.Body))

		subjects := uniqueSubjects(sparqlQuery.Triples)
		err = sendMessage(strings.Join(subjects, "\", \""), request.Body)
		if err != nil {
			return "Error sending SNS message", err
		}
	}
	resp.Body.Close()
	return string(respBody), nil
}

func main() {
	lambda.Start(Handler)
}

func uniqueSubjects(in []sparql.Triple) []string {
	u := make([]string, 0, len(in))
	m := make(map[string]bool)

	for _, val := range in {
		if _, ok := m[val.Subject]; !ok {
			m[val.Subject] = true
			u = append(u, val.Subject)
		}
	}

	return u
}

func sendMessage(subjects string, document string) error {
	message := fmt.Sprintf("{\"action\": \"touch\", \"entities\": [\"%s\"], \"body\": \"%s\"}", subjects, document)
	topicArn := os.Getenv("RIALTO_TOPIC_ARN")
	endpoint := os.Getenv("RIALTO_SNS_ENDPOINT")
	snsConn := sns.New(session.New(), aws.NewConfig().
		WithDisableSSL(false).
		WithEndpoint(endpoint))
	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: &topicArn,
	}
	_, err := snsConn.Publish(input)
	return err
}
