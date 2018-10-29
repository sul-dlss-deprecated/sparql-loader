# Rialto SparqlLoader

[![CircleCI](https://circleci.com/gh/sul-dlss-labs/sparql-loader.svg?style=svg)](https://circleci.com/gh/sul-dlss-labs/sparql-loader)

## Calling a lambda over http through API Gateway

```
curl --http1.1 --data-urlencode update@basic_insert.txt -H 'X-API-Key: <API KEY>' https://d176x3sh52.execute-api.us-west-2.amazonaws.com/development/rialto-sparql-loader
```

Where `basic_insert.txt` is a text file with a sparql query like:
```
PREFIX dc: <http://purl.org/dc/elements/1.1/> INSERT DATA { <http://example/book3> dc:title    'A new book' ; dc:creator  'A.N.Other' . }
```

## Testing


### Install dependencies

```shell

### Unit testing

```shell
pytest -vv -k unit
```

### Integration testing

1. Start localstack and blazegraph via docker.
```
docker-compose up
```

3. Run the test
```shell
pytest -vv
```

### Building an AWS Lambda deployment package

Per the [AWS Documentation](https://docs.aws.amazon.com/lambda/latest/dg/lambda-python-how-to-create-deployment-package.html), a deployment package is made from the `virtualenv` installed dependencies.

1. Create a [python virtual environment](https://docs.python-guide.org/dev/virtualenvs/#lower-level-virtualenv)
1. Activate your virtual environment

```shell
source env/bin/activate
```

1. Install dependencies

```shell
pip install -r requirements.txt
```

1. Create zip file

```shell
zip sparql-loader.zip handler.py sns_client.py neptune_client.py
```

1. Copy dependencies into zip file

```shell
cd env/lib/python3.6/site-packages/
zip -r ../../../../sparql-loader.zip honeybadger isodate psutil rdflib rdflib_sparql requests
```

Note: We are packaging the minimum level of dependencies to try to keep our deployment package small.

