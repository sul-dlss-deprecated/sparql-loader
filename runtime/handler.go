package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	honeybadger "github.com/honeybadger-io/honeybadger-go"
	"github.com/sul-dlss-labs/sparql-loader/sparql"
)

const urlEncoded = "application/x-www-form-urlencoded"

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

func (p *ProxyHandler) isValidContentType(contentType string) bool {
	return contentType == "application/sparql-update" || contentType == urlEncoded || contentType == "application/sparql-query"
}

// RequestHandler is the AWS Lambda proxy handler called by main.  It only handles POST requests
func (p *ProxyHandler) RequestHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// See https://www.w3.org/TR/sparql11-protocol/#query-via-post-urlencoded
	contentType := request.Headers["Content-Type"]
	if contentType == urlEncoded && !correctlyURIEncoded(request.Body) {
		return &events.APIGatewayProxyResponse{StatusCode: 422, Body: "[MalformedRequest] query string not properly escaped"}, nil
	} else if !p.isValidContentType(contentType) {
		body := fmt.Sprintf("[MalformedRequest] Invalid Content-Type: '%s'", contentType)
		return &events.APIGatewayProxyResponse{StatusCode: 422, Body: body}, nil
	}

	res, err := p.registry.Writer.Post(request.Body, contentType)
	if err != nil {
		honeybadger.Notify(err)
		return nil, err
	}

	if res.StatusCode == 400 {
		body := fmt.Sprintf("[BadRequest] There was a problem with the request:\n%s", res.Body)
		return &events.APIGatewayProxyResponse{StatusCode: 400, Body: body}, nil
	}

	start := time.Now()
	log.Printf("SPARQL parse begin: %s", start)

	message := p.formatMessage(request.Body, contentType)
	log.Printf("SPARQL parse elapsed time: %s", time.Since(start))
	if time.Since(start).Seconds() > 5 { // Log if the elapsed time is > 5 seconds. TODO: Make this configurable
		log.Printf("SPARQL Query: \n%s", request.Body)
	}

	if message != nil {
		start = time.Now()
		log.Printf("SNS publish begin: %s", start)
		err := p.registry.Publisher.Publish(string(message))

		if err != nil {
			honeybadger.Notify(err)
			return nil, err
		}
		log.Printf("SNS publish elapsed time: %s", time.Since(start))
	}
	return res, nil
}

func (p *ProxyHandler) formatMessage(body string, contentType string) []byte {
	if contentType == "application/sparql-query" || (contentType == urlEncoded && !strings.HasPrefix(body, "update=")) {
		return nil
	}
	sparqlQuery := sparql.NewQuery()

	var queryString string
	if contentType == urlEncoded {
		queryString, _ = url.QueryUnescape(strings.Replace(body, "update=", "", -1))
	} else {
		queryString = body
	}

	_ = sparqlQuery.Parse(queryString)

	subjects := []string{}
	for _, part := range sparqlQuery.Parts {
		subjects = append(subjects, uniqueSubjects(part.Graph)...) // The ... unpacks the returned []string
	}
	message, _ := json.Marshal(&snsMessage{Action: "touch", Entities: removeDuplicates(subjects)})
	return message
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

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] {
			// Do not add duplicate.
			continue
		}
		// Record this element as an encountered element.
		encountered[elements[v]] = true
		// Append to result slice.
		result = append(result, elements[v])
	}
	// Return the new slice.
	return result
}
