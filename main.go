package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is the Lambda function handler
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (string, error) {

	form := url.Values{}
	form.Add("update", request.Body)
	proxyReq, err := http.NewRequest("POST", os.Getenv("RIALTO_SPARQL_ENDPOINT"), strings.NewReader(form.Encode()))
	proxyReq.Header = make(http.Header)

	proxyReq.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))

	httpClient := http.Client{}
	resp, err := httpClient.Do(proxyReq)
	respBody, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return string(respBody), err
	}

	/*
		sp.sendMessage(body)
	*/
	resp.Body.Close()
	return string(respBody), nil
}

func main() {
	lambda.Start(Handler)
}

/*
func (sp *sparqlProxy) sendMessage(document string) error {
	message := fmt.Sprintf("{\"action\": \"sparql_update\", \"document\": %s}", document)
	return sp.messageService.Publish(message)
}
*/
