
import json
import logging
import os
import time
import urllib.parse

from rdflib.plugins.sparql.parser import parseUpdate
from rdflib.plugins.sparql.algebra import translateUpdate

from sns_client import SnsClient
from neptune_client import NeptuneClient

# Set a constant for the statement delimiter used for parsing
STATEMENT_DELIMITER = "}};"


def main(event, _):
    # Setup the logger at the INFO level while we continue to profile
    logger = logging.getLogger()
    logger.setLevel(logging.INFO)

    rialto_sparql_endpoint = os.getenv('RIALTO_SPARQL_ENDPOINT', "http://localhost:8080/bigdata/namespace/kb/sparql")
    rialto_sns_endpoint = os.getenv('RIALTO_SNS_ENDPOINT', "http://localhost:4575")
    rialto_topic_arn = os.getenv('RIALTO_TOPIC_ARN', "rialto")
    aws_region = os.getenv('AWS_REGION', "us-west-2")

    request_content_type = event['headers']['Content-Type']  # capture this value for use throughout

    sns_client = SnsClient(rialto_sns_endpoint, rialto_topic_arn, aws_region)
    neptune_client = NeptuneClient(rialto_sparql_endpoint)

    start_time = time.time()
    logger.info("NEPTUNE START: " + time.asctime(time.localtime(start_time)))
    response, status_code = neptune_client.post(event['body'], request_content_type)
    logger.info("NEPTUNE ELAPSED: %f" % (time.time() - start_time))

    if status_code == 200:
        if "update=" in event['body'] or request_content_type == "application/sparql-update":
            start_time = time.time()
            logger.info("SPARQL PARSE START: " + time.asctime(time.localtime(start_time)))
            entities = get_unique_subjects(
                            get_entities(
                                urllib.parse.unquote_plus(
                                    event['body']).replace('update=', '')))
            logger.info("SPARQL PARSE ELAPSED: %f" % (time.time() - start_time))

            if entities:
                message = {"Action": "touch", "Entities": entities}
                _ = sns_client.publish(json.dumps(message))  # currently not using the neptune response

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
    for block in translateUpdate(parseUpdate(body)):
        for key in block.keys():
            if key in ['delete', 'insert']:
                subjects += get_subjects_from_quads(block[key]['quads'])
                subjects += get_subjects_from_triples(block[key]['triples'])
            if key in ['quads']:
                subjects += get_subjects_from_quads(block['quads'])
            if key in ['triples']:
                subjects += get_subjects_from_triples(block['triples'])

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
