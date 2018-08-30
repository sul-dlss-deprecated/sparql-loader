package runtime

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	honeybadger "github.com/honeybadger-io/honeybadger-go"
	"github.com/sul-dlss-labs/sparql-loader/sparql"
)

// Handler is an interface that is called by main to allow handler dependency injection
type Handler interface {
	RequestHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
}

// ProxyHandler allows the registry to be set outside of the normal handler operation
type ProxyHandler struct {
	registry Registry
}

type snsMessage struct {
	Action   string
	Entities []string
}

// NewHandler creates an new ProxyHandler instance
func NewHandler(registry *Registry) *ProxyHandler {
	return &ProxyHandler{registry: *registry}
}

// RequestHandler is the AWS Lambda proxy handler called by main
func (p *ProxyHandler) RequestHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if !correctlyURIEncoded(request.Body) {
		return &events.APIGatewayProxyResponse{StatusCode: 422, Body: "[MalformedRequest] query string not properly escaped"}, nil
	}

	res, err := p.registry.Writer.Post(request.Body)
	if err != nil {
		honeybadger.Notify(err)
		return nil, err
	}

	if res.StatusCode == 400 {
		// log.Printf("There was a problem with the request (%v) %s", resp.StatusCode, respBody)
		return &events.APIGatewayProxyResponse{StatusCode: 400, Body: "[BadRequest] There was a problem with the request"}, nil
	}

	message := p.formatMessage(request.Body)
	if message != nil {
		err := p.registry.Publisher.Publish(string(message))

		if err != nil {
			honeybadger.Notify(err)
			return nil, err
		}
	}
	return res, nil
}

func (p *ProxyHandler) formatMessage(body string) []byte {
	if strings.HasPrefix(body, "update=") {
		sparqlQuery := sparql.NewQuery()

		queryString, _ := url.QueryUnescape(strings.Replace(body, "update=", "", -1))
		_ = sparqlQuery.Parse(queryString)

		for _, part := range sparqlQuery.Parts {
			message, _ := json.Marshal(&snsMessage{Action: "touch", Entities: uniqueSubjects(part.Graph)})
			return message
		}
	}
	return nil
}

// Returns true if the provided string is correctly URI encoded
func correctlyURIEncoded(bodyIn string) bool {
	unescaped, _ := url.QueryUnescape(bodyIn)
	if bodyIn == unescaped {
		return false
	}
	return true
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
