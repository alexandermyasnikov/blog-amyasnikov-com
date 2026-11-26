from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse, parse_qs
import json

class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        parsed_url = urlparse(self.path)
        query_params = parse_qs(parsed_url.query)

        response = {
            "method": "GET",
            "path": parsed_url.path,
            "query_params": query_params,
            "headers": dict(self.headers)
        }

        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(response).encode('utf-8'))

def run_server(port=8000):
    server_address = ('127.0.0.1', port)
    httpd = HTTPServer(server_address, RequestHandler)
    print(f'Start server on {port}')
    httpd.serve_forever()

if __name__ == '__main__':
    run_server()