import requests


class NeptuneClient():
    def __init__(self, sparql_endpoint):
        self.sparql_endpoint = sparql_endpoint

    def post(self, request_body):
        response = requests.post(self.sparql_endpoint,
                                 data=request_body,
                                 headers={"Content-Type": "application/sparql-update"})
        return (response.text, response.status_code)
