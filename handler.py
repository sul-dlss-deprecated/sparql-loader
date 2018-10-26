
import os
import logging
import json
import urllib.parse

from rdflib.plugins.sparql.parser import parseUpdate
from rdflib.plugins.sparql.algebra import translateUpdate

from sns_client import SnsClient
from neptune_client import NeptuneClient

def main(event, context):
    logger = logging.getLogger()
    logger.setLevel(logging.INFO)
    logger.info("EVENT = "+json.dumps(event))
    rialto_sparql_endpoint = os.getenv('RIALTO_SPARQL_ENDPOINT', "localhost:8080")
    rialto_sparql_path = os.getenv('RIALTO_SPARQL_PATH', "/bigdata/namespace/kb/sparql")
    rialto_sns_endpoint = os.getenv('RIALTO_SNS_ENDPOINT', "http://localhost:4575")
    rialto_topic_arn = os.getenv('RIALTO_TOPIC_ARN', "rialto")
    aws_region = os.getenv('AWS_REGION', "us-west-2")

    sns_client = SnsClient(rialto_sns_endpoint, rialto_topic_arn, aws_region)
    neptune_client = NeptuneClient(rialto_sparql_endpoint, rialto_sparql_path)

    response = neptune_client.post(event)
    logger.info("NEPTUNE_RESPONSE = "+str(response))

    if response['statusCode'] == 200:
        entities = getUniqueSubjects(getEntities(urllib.parse.unquote_plus(event['body']).replace('update=','')))
        message = "{'Action': 'touch', 'Entities': %s}" % entities
        sns_response = sns_client.publish(message)
        logger.info("SNS_RESPONSE = "+json.dumps(sns_response))

    return {
        'body' : str(response['body']),
        'statusCode' : response['statusCode']
	}

def getEntities(body):
    delimeter_count = body.count("}};") # determine if the sparql query is broken up by "}};"
    if delimeter_count in [0, 1]:
        return parseBody(body)

    subjects = []
    for chunk in body.split("}};"):
        if len(chunk.rstrip()) > 0:
            subjects += parseBody(chunk+"}};") # append the "}};" that was removed by split

    return subjects

def parseBody(body):
    subjects = []
    for block in translateUpdate(parseUpdate(body)):
        for key in block.keys():
            if key in ['delete', 'insert']:
                subjects += getSubjectsFromQuads(block[key]['quads'])
                subjects += getSubjectsFromTriples(block[key]['triples'])
            if key in ['quads']:
                subjects += getSubjectsFromQuads(block['quads'])
            if key in ['triples']:
                subjects += getSubjectsFromTriples(block['triples'])
    
    return subjects

def getSubjectsFromQuads(block):
    subjects = []
    for key in block.keys():
        for s, _p, _o in block[key]:
            subjects.append(s.toPython())
    
    return subjects

def getSubjectsFromTriples(block):
    subjects = []
    for s, _p, _o in block:
        subjects.append(s.toPython())

    return subjects

def getUniqueSubjects(subjectsList):
    unique_subjects = []
    for subject in subjectsList:
        if subject not in unique_subjects:
            unique_subjects.append(subject)
    
    return unique_subjects
