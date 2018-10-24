
import os
import rdflib

from sns_client import SnsClient
from neptune_client import NeptuneClient

def main(event, context):
    # honeybadger_api_key = os.getenv('HONEYBADGER_API_KEY', "12345")
    rialto_sparql_endpoint = os.getenv('RIALTO_SPARQL_ENDPOINT', "localhost:8080")
    rialto_sparql_path = os.getenv('RIALTO_SPARQL_PATH', "/bigdata/namespace/kb/sparql")
    rialto_sns_endpoint = os.getenv('RIALTO_SNS_ENDPOINT', "http://localhost:4575")
    rialto_topic_arn = os.getenv('RIALTO_TOPIC_ARN', "rialto")
    aws_region = os.getenv('AWS_REGION', "us-west-2")
    sns_client = SnsClient(rialto_sns_endpoint, rialto_topic_arn, aws_region)
    neptune_client = NeptuneClient(rialto_sparql_endpoint, rialto_sparql_path)

    response = neptune_client.post(event)

    print("RESPONSE")
    print(response)

    if response['statusCode'] == 200:
        sns_response = sns_client.publish(formatMessage(event['body']))
        return sns_response

    return ""

def formatMessage(body):
    graph = rdflib.Graph()
    graph.update(body)
    subjects = []

    for s, _p, _o in graph:
        subjects.append(s.toPython())

    return "{'Action': 'touch', 'Entities': subjects}"
