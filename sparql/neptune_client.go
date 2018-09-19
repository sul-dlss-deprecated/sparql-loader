package sparql

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// NeptuneClient represents the endpoint used for AWS Neptune HTTP requests
type NeptuneClient struct {
	endpoint string
}

// NewNeptuneClient returns a new NeptuneClient instance
func NewNeptuneClient(endpoint string) *NeptuneClient {
	return &NeptuneClient{endpoint: endpoint}
}

// Post exposes the POST function to the handler
func (n *NeptuneClient) Post(query string, contentType string) (*events.APIGatewayProxyResponse, error) {
	res, err := n.HTTPProxy(query, contentType)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// HTTPProxy passes any query string directly to the NeptuneClient endpoint
func (n *NeptuneClient) HTTPProxy(query string, contentType string) (*events.APIGatewayProxyResponse, error) {

	proxyReq, _ := http.NewRequest("POST", n.endpoint, strings.NewReader(query))
	proxyReq.Header.Set("Content-Type", contentType)

	httpClient := http.Client{}

	resp, err := httpClient.Do(proxyReq)
	respBody, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{Body: string(respBody), StatusCode: resp.StatusCode}, nil
}
