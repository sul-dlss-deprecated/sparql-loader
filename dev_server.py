from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
import handler
import logging
import os

log = logging.getLogger(__name__)

# A server wrapper for the lambda that can be used in a local dev environment.
# This is necessary because localstack does not support API gateway to lambda
# that uses a post body.


class SparqlLoaderRequestHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        log.info("Handling post")
        # Read the body.
        try:
            log.debug("Reading request")
            content_length = int(self.headers['Content-Length'])
            body = self.rfile.read(content_length).decode('utf-8')
        finally:
            self.rfile.close()
        log.debug("Processing request")
        response = handler.main({'body': body, 'headers': {'Content-Type': self.headers['Content-Type']}}, None)
        log.info("Response status code is %s", response['statusCode'])
        self.send_response(response['statusCode'])
        self.end_headers()
        if response['body']:
            self.wfile.write(bytes(response['body'], "utf8"))


def run(server_class=ThreadingHTTPServer, handler_class=SparqlLoaderRequestHandler):
    log.info("Starting server")
    server_address = ('0.0.0.0', 8080)
    httpd = server_class(server_address, handler_class)
    httpd.serve_forever()


if __name__ == '__main__':
    log_level = logging.DEBUG if os.getenv('DEBUG', 'false').lower() == 'true' else logging.INFO
    logging.basicConfig(level=log_level)
    run()
