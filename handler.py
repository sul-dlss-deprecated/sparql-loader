
import os
import logging
import json
import urllib.parse
import time

from rdflib.plugins.sparql.parser import parseUpdate
from rdflib.plugins.sparql.algebra import translateUpdate

from sns_client import SnsClient
from neptune_client import NeptuneClient

def main(event, context):
    # Setup the logger at the INFO level while we continue to profile
    logger = logging.getLogger()
    logger.setLevel(logging.INFO)

    rialto_sparql_endpoint = os.getenv('RIALTO_SPARQL_ENDPOINT', "localhost:8080")
    rialto_sparql_path = os.getenv('RIALTO_SPARQL_PATH', "/bigdata/namespace/kb/sparql")
    rialto_sns_endpoint = os.getenv('RIALTO_SNS_ENDPOINT', "http://localhost:4575")
    rialto_topic_arn = os.getenv('RIALTO_TOPIC_ARN', "rialto")
    aws_region = os.getenv('AWS_REGION', "us-west-2")

    sns_client = SnsClient(rialto_sns_endpoint, rialto_topic_arn, aws_region)
    neptune_client = NeptuneClient(rialto_sparql_endpoint, rialto_sparql_path)

    start_time = time.time()
    logger.info("NEPTUNE START: " + time.asctime( time.localtime(start_time)))
    response = neptune_client.post(event)
    logger.info("NEPTUNE ELAPSED: %f" % (time.time() - start_time))

    if response['statusCode'] == 200:
        start_time = time.time()
        logger.info("SPARQL PARSE START: " + time.asctime( time.localtime(start_time)))
        entities = getUniqueSubjects(getEntities(urllib.parse.unquote_plus(event['body']).replace('update=','')))
        logger.info("SPARQL PARSE ELAPSED: %f" % (time.time() - start_time))

        message = "{'Action': 'touch', 'Entities': %s}" % entities

        start_time = time.time()
        logger.info("SNS START: " + time.asctime( time.localtime(start_time)))
        sns_response = sns_client.publish(message)
        logger.info("SNS ELAPSED: %f" % (time.time() - start_time))

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
