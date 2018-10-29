import http.client


class NeptuneClient():
    def __init__(self, sparql_endpoint, sparql_path):
        self.sparql_endpoint = sparql_endpoint
        self.sparql_path = sparql_path

    def post(self, request_body):
        http_conn = http.client.HTTPConnection(self.sparql_endpoint)
        http_conn.request('POST',
                          self.sparql_path,
                          body=request_body,
                          headers={"Content-Type": "application/sparql-update"})
        response = http_conn.getresponse()
        return (response.read(), response.status)
