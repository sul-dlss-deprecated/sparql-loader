
import os
import json
import urllib.parse

from rdflib.plugins.sparql.parser import parseUpdate
from rdflib.plugins.sparql.algebra import translateUpdate

from sns_client import SnsClient
from neptune_client import NeptuneClient


def main(event, _):
    rialto_sparql_endpoint = os.getenv('RIALTO_SPARQL_ENDPOINT', "http://localhost:8080/bigdata/namespace/kb/sparql")
    rialto_sns_endpoint = os.getenv('RIALTO_SNS_ENDPOINT', "http://localhost:4575")
    rialto_topic_arn = os.getenv('RIALTO_TOPIC_ARN', "rialto")
    aws_region = os.getenv('AWS_REGION', "us-west-2")

    sns_client = SnsClient(rialto_sns_endpoint, rialto_topic_arn, aws_region)
    neptune_client = NeptuneClient(rialto_sparql_endpoint)

    response, status_code = neptune_client.post(event['body'])

    if status_code == 200:
        if "update=" in event['body'] or event['Content-Type'] in ["application/sparql-update"]:
            entities = get_unique_subjects(
                            get_entities(
                                urllib.parse.unquote_plus(
                                    event['body']).replace('update=', '')))
            if entities:
                message = {"Action": "touch", "Entities": entities}
                _ = sns_client.publish(json.dumps(message))  # currently not using the neptune response

    return {
        'body': str(response),
        'statusCode': status_code
    }


def get_entities(body):
    delimeter_count = body.count("}};")  # determine if the sparql query is broken up by "}};"
    if delimeter_count in [0, 1]:
        return parse_body(body)

    subjects = []
    for chunk in body.split("}};"):
        if len(chunk.rstrip()) > 0:
            subjects += parse_body(chunk+"}};")  # append the "}};" that was removed by split

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
    unique_subjects = []
    for subject in subjectsList:
        if subject not in unique_subjects:
            unique_subjects.append(subject)

    return unique_subjects
