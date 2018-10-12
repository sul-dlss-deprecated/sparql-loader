package message

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// Client represents the endpoint used for AWS SNS Message publishing
type Client struct {
	endpoint string
	topic    string
	region   string
	ssl      bool

}

// NewClient returns a new Client instance
func NewClient(endpoint string, topic string, region string, ssl bool) *Client {
	return &Client{endpoint: endpoint, topic: topic, region: region, ssl: ssl}
}

// Publish exposes the SNS Publish function to the handler
func (n *Client) Publish(message string) error {
	snsConn := sns.New(session.New(), aws.NewConfig().
		WithDisableSSL(!n.ssl).
		WithEndpoint(n.endpoint).
		WithRegion(n.region))
	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: &n.topic,
	}
	_, err := snsConn.Publish(input)
	return err
}
