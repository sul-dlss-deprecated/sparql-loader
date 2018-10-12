FROM golang:latest

WORKDIR /go/src/github.com/sul-dlss-labs/sparql-loader

ENV AWS_ACCESS_KEY_ID 999999
ENV AWS_SECRET_ACCESS_KEY 1231
ENV NO_SSL true
ENV RIALTO_TOPIC_ARN arn:aws:sns:us-east-1:123456789012:data-update
ENV AWS_REGION us-east-1

EXPOSE 8080

# Code must be linked in to /go/src/github.com/sul-dlss-labs/sparql-loader
# Following env variables are required: RIALTO_SPARQL_ENDPOINT, RIALTO_SPARQL_ENDPOINT

CMD ["go", "run", "cmd/server/main.go"]
