import http.client

class NeptuneClient():
    def __init__(self, sparql_endpoint, sparql_path):
        self.sparql_endpoint = sparql_endpoint
        self.sparql_path = sparql_path

    def post(self,event):
        # return {'body': "TEST", 'statusCode': 200}
        http_conn = http.client.HTTPConnection(self.sparql_endpoint)
        http_conn.request('POST', self.sparql_path, body=event['body'], headers={"Content-Type": event['content_type']})
        response = http_conn.getresponse()
        return {'body': response.read(), 'statusCode': response.status}
