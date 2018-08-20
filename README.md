# Rialto SparqlLoader

[![CircleCI](https://circleci.com/gh/sul-dlss-labs/sparql-loader.svg?style=svg)](https://circleci.com/gh/sul-dlss-labs/sparql-loader)

## Running a lambda on localstack with API Gateway

Note: This is the initial setup of a lambda to proxy to the AWS Neptune SPARQL endpoint. Some of the steps below may need editing.

## Calling a lambda over http through API Gateway

```
curl -i -d@./basic_insert.txt -X POST https://<API_GATEWAY_ENDPOINT>/<DEPLOYMENT}>
```

Where `basic_insert.txt` is a text file with a sparql query like:
```
PREFIX dc: <http://purl.org/dc/elements/1.1/> INSERT DATA { <http://example/book3> dc:title    'A new book' ; dc:creator  'A.N.Other' . }
```

## Create and upload the lambda
1. Create a zip
```
GOOS=linux go build -o main
zip lambda.zip main
```

2. Start localstack. If you're on a Mac, ensure you are running the docker daemon.
```
SERVICES=lambda,apigateway LAMBDA_EXECUTOR=docker localstack start
```

3. Upload zip and create a function definition
```
AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 aws \
--endpoint-url http://localhost:4574 lambda create-function \
--function-name proxy \
--runtime go1.x \
--role r1 \
--handler main \
--zip-file fileb://lambda.zip
```

4. Verify the LAMBDA_ARN
```
awslocal lambda list-functions
```

5. Create a REST API
```
awslocal apigateway create-rest-api \
    --region us-east-1 \
    --name sparqlProxy

OUTPUT:
{
    "name": "sparqlProxy", 
    "id": "987918802A-Z", 
    "createdDate": 1530302743
}
```

6. Get the PARENT_RESOURCE_ID
```
awslocal apigateway get-resources --rest-api-id ${API_ID_FROM_ABOVE} --query 'items[?path==`/`].id' --output text

OUTPUT:

527A-Z825141
```

7. Create the resource
```
awslocal apigateway create-resource \
    --region us-east-1 \
    --rest-api-id {THE_API_ID_FROM_ABOVE} \
    --parent-id {THE_PARENT_RESOURCE_ID} \
    --path-part "sparqlProxy"

OUTPUT:
{
    "resourceMethods": {
        "GET": {}
    }, 
    "pathPart": "sparqlProxy", 
    "parentId": "527A-Z825141", 
    "path": "/sparqlProxy", 
    "id": "253851A-Z274"
}
```

8. Put a method
```
awslocal apigateway put-method \
    --region us-east-1 \
    --rest-api-id {API_ID} \
    --resource-id {RESOURCE_ID} \
    --http-method POST \
    --authorization-type "NONE" \

OUTPUT:
{
    "httpMethod": "GET", 
    "authorizationType": "NONE"
}
```

9. Put the integration
```
awslocal apigateway put-integration \
    --region us-east-1 \
    --rest-api-id {API_ID} \
    --resource-id {RESOURCE_ID} \
    --http-method POST \
    --type AWS_PROXY \
    --integration-http-method POST \
    --uri arn:aws:apigateway:${REGION}:lambda:path/2015-03-31/functions/arn:aws:lambda:localstack:000000000000:function:proxy/invocations \
    --passthrough-behavior WHEN_NO_MATCH \

OUTPUT:
{
    "httpMethod": "GET", 
    "integrationResponses": {
        "200": {
            "responseTemplates": {
                "application/json": null
            }, 
            "statusCode": 200
        }
    }, 
    "type": "AWS_PROXY", 
    "uri": "arn:aws:apigateway::lambda:path/2015-03-31/functions/arn:aws:lambda:localstack:000000000000:function:proxy/invocations"
}
```

10. Create the deployment
```
awslocal apigateway create-deployment \
    --region us-east-1 \
    --rest-api-id {REST_API_ID} \
    --stage-name dev \

OUTPUT:
{
    "description": "", 
    "id": "1706090266", 
    "createdDate": 1530304218
}
```

11. The endpoint

```
http://localhost:4567/restapis/987918802A-Z/dev/_user_request_/sparqlProxy
```
