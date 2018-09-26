# Rialto SparqlLoader

[![CircleCI](https://circleci.com/gh/sul-dlss-labs/sparql-loader.svg?style=svg)](https://circleci.com/gh/sul-dlss-labs/sparql-loader)

## Calling a lambda over http through API Gateway

```
curl --http1.1 --data-urlencode update@basic_insert.txt -X 'X-API-Key: <API KEY>' https://d176x3sh52.execute-api.us-west-2.amazonaws.com/development/rialto-sparql-loader
```

Where `basic_insert.txt` is a text file with a sparql query like:
```
PREFIX dc: <http://purl.org/dc/elements/1.1/> INSERT DATA { <http://example/book3> dc:title    'A new book' ; dc:creator  'A.N.Other' . }
```

## Testing

### Unit testing

```shell
go test -v ./... -short
```

### Integration testing

1. Start localstack.
```
SERVICES=sns localstack start
```

2. Start Blazegraph
```
export JAVA_HOME="$(/usr/libexec/java_home -v 1.8)"
java -server -Xmx4g -jar blazegraph.jar
```

3. Run the test
```shell
go test -v ./...
```

or 
```shell
go test -v ./test
```
*To only run the integration test.*


**NOTE:** We do not upload a lambda into localstack or setup API Gateway as it currently does not support the body pass through that we are using.

