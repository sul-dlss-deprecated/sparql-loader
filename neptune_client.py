import requests


class NeptuneClient():
    def __init__(self, sparql_endpoint):
        self.sparql_endpoint = sparql_endpoint

    def post(self, request_body, request_content_type):
        response = requests.post(self.sparql_endpoint,
                                 data=request_body.encode('utf-8'),
                                 headers={"Content-Type": request_content_type})
        return response.text, response.status_code
