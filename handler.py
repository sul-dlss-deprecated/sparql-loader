import json
import logging
import os
import sys
import time
import urllib.parse

from rdflib.plugins.sparql.parser import parseUpdate
from rdflib.plugins.sparql.algebra import translateUpdate
from pyparsing import ParseException

from sns_client import SnsClient
from neptune_client import NeptuneClient

# Set a constant for the statement delimiter used for parsing
STATEMENT_DELIMITER = "}};"
URL_ENCODED = "application/x-www-form-urlencoded"
SPARQL_UPDATE = "application/sparql-update"
VALID_CONTENT_TYPES = [
    "application/sparql-query",
    SPARQL_UPDATE,
    URL_ENCODED
]

# Setup the logger at the INFO level while we continue to profile
logger = logging.getLogger()
logger.setLevel(logging.INFO)

# Set the recursion limit to catch long sparql statements
sys.setrecursionlimit(10000)


def main(event, _):
    rialto_sparql_endpoint = os.getenv('RIALTO_SPARQL_ENDPOINT', "http://localhost:8080/bigdata/namespace/kb/sparql")
    rialto_sns_endpoint = os.getenv('RIALTO_SNS_ENDPOINT', "http://localhost:4575")
    rialto_topic_arn = os.getenv('RIALTO_TOPIC_ARN', "rialto")
    aws_region = os.getenv('AWS_REGION', "us-west-2")

    request_body = event['body']
    request_content_type = event['headers']['Content-Type']

    # Verify that the body is properly encoded and has a supported content type
    # BEFORE setting up AWS resources
    clean_request_content_type = clean_content_type(request_content_type)
    verify_query = is_malformed_query(request_body, clean_request_content_type)
    if verify_query is not None:
        logger.warning("Received bad query: {}".format(verify_query['body']))
        return verify_query

    sns_client = SnsClient(rialto_sns_endpoint, rialto_topic_arn, aws_region)
    neptune_client = NeptuneClient(rialto_sparql_endpoint)

    start_time = time.time()
    response, status_code = neptune_client.post(request_body, request_content_type)
    logger.info("NEPTUNE ELAPSED: %f" % (time.time() - start_time))

    if status_code == 200:
        entities = []
        start_time = time.time()
        if "update=" in request_body and clean_request_content_type == URL_ENCODED:
            entities = get_unique_subjects(
                            get_entities(
                                urllib.parse.unquote_plus(
                                    request_body).replace('update=', '')))

        if clean_request_content_type == SPARQL_UPDATE:
            entities = get_unique_subjects(
                            get_entities(request_body))

        if entities:
            message = {"Action": "touch", "Entities": entities}
            _ = sns_client.publish(json.dumps(message))  # currently not using the neptune response

        logger.info("SPARQL PARSE ELAPSED: %f" % (time.time() - start_time))
    else:
        logger.error("NEPTUNE RETURNED %s: %s" % (status_code, response))

    return {
        'body': str(response),
        'statusCode': status_code
    }


def get_entities(body):
    delimiter_count = body.count(STATEMENT_DELIMITER)  # determine if the sparql query is broken up by "}};"
    if delimiter_count in [0, 1]:
        return parse_body(body)

    subjects = []
    for chunk in body.split(STATEMENT_DELIMITER):
        if len(chunk.rstrip()) > 0:
            subjects += parse_body(chunk+STATEMENT_DELIMITER)  # append the "}};" that was removed by split

    return subjects


def parse_body(body):
    subjects = []
    try:
        for block in translateUpdate(parseUpdate(body)):
            for key in block.keys():
                if key in ['delete', 'insert']:
                    subjects += get_subjects_from_quads(block[key]['quads'])
                    subjects += get_subjects_from_triples(block[key]['triples'])
                if key in ['quads']:
                    subjects += get_subjects_from_quads(block['quads'])
                if key in ['triples']:
                    subjects += get_subjects_from_triples(block['triples'])
    except (RecursionError, ParseException):
        # Swallow a parse error, since the sparql made it to Neptune
        logger.error("SPARQL ERROR PARSING: %s", body)

    return subjects


def get_subjects_from_quads(block):
    subjects = []
    for key in block.keys():
        for s, _p, _o in block[key]:
            subjects.append(s.toPython())

    return subjects


def get_subjects_from_triples(block):
    subjects = []
    for s, _p, _o in block:
        subjects.append(s.toPython())

    return subjects


def get_unique_subjects(subjectsList):
    # return a sorted list of unique subjects. Not sure if this is idiomatic however
    return list(set(subjectsList))


# returns None on happy path, otherwise returns error structure to be passed through
def is_malformed_query(body, content_type):
    if content_type == URL_ENCODED:
        if not correctly_uri_encoded(body):
            return {'body': "[MalformedRequest] query string not properly escaped",
                    'statusCode': 422}
    elif content_type not in VALID_CONTENT_TYPES:
        return {'body': "[MalformedRequest] Invalid Content-Type: '%s'" % content_type,
                'statusCode': 422}

    return None


# Returns true if the provided string is correctly URI encoded
def correctly_uri_encoded(body):
    unescaped = urllib.parse.unquote_plus(body)
    if body == unescaped:
        return False

    return True


# Cleans content type, e.g., to remove charset=utf-8
def clean_content_type(content_type):
    if not content_type:
        return content_type
    return content_type.split(';')[0]
