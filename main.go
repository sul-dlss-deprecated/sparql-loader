package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/sul-dlss-labs/sparql-loader/sparql"
)

// Handler is the Lambda function handler
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (string, error) {
	// MAP of queries that will trigger a message to SNS
	knownQueries := map[string]bool{
		"INSERT":      true,
		"INSERT DATA": true,
		"DELETE":      true,
		"DELETE DATA": true,
	}

	log.Printf("REQUST BODY: %s", request.Body)

	proxyReq, _ := http.NewRequest("POST", os.Getenv("RIALTO_SPARQL_ENDPOINT"), encodeBody(request.Body))
	proxyReq.Header = make(http.Header)

	proxyReq.Header.Add("Content-type", "application/x-www-form-urlencoded")
	proxyReq.Header.Add("Accept", "application/json")
	proxyReq.Header.Set("Content-Length", strconv.Itoa(len(request.Body)))

	httpClient := http.Client{}

	resp, err := httpClient.Do(proxyReq)
	respBody, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return string(respBody), err
	}

	if resp.StatusCode == 400 {
		log.Printf("There was a problem with the request (%v) %s", resp.StatusCode, respBody)
		return fmt.Sprintf("[BadRequest] %s", respBody), nil
	}

	log.Printf("NEPTUNE [SUCCESS]: %s", string(respBody))

	if strings.HasPrefix(request.Body, "update=") {
		sparqlQuery := sparql.NewQuery()
		_ = sparqlQuery.Parse(strings.NewReader(strings.Replace(request.Body, "update=", "", -1)))

		for _, part := range sparqlQuery.Parts {
			if knownQueries[strings.ToUpper(part.Verb)] {
				err = sendMessage("touch", uniqueSubjects(part.Graph), request.Body)
				if err != nil {
					return "Error sending SNS message", err
				}
			}
		}
	}
	resp.Body.Close()
	return string(respBody), nil
}

func main() {
	lambda.Start(Handler)
}

func encodeBody(bodyIn string) *strings.Reader {
	// The body that we get from API gateway is decoded, so we need to re-encode it
	// However, we need to first remove the query/update directive before encoding.
	isUnescaped, _ := url.QueryUnescape(bodyIn)
	if bodyIn == isUnescaped {
		i := strings.Index(bodyIn, "=")
		queryPrefix := bodyIn[:i]
		bodyIn = bodyIn[i+1:]
		bodyIn = fmt.Sprintf("%s=%s", queryPrefix, url.QueryEscape(bodyIn))
	}
	return strings.NewReader(bodyIn)
}

func uniqueSubjects(in []sparql.Triple) []string {
	u := make([]string, 0, len(in))
	m := make(map[string]bool)

	for _, val := range in {
		val.Subject = strings.Replace(val.Subject, "<", "", -1)
		val.Subject = strings.Replace(val.Subject, ">", "", -1)
		if _, ok := m[val.Subject]; !ok {
			m[val.Subject] = true
			u = append(u, val.Subject)
		}
	}

	return u
}

func sendMessage(action string, subjects []string, document string) error {
	message := fmt.Sprintf("{\"Action\": \"%s\", \"Entities\": [\"%s\"], \"Body\": \"\"}", action, strings.Join(subjects, "\", \""))
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
