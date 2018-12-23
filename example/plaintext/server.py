# Python 2.7.x

from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer

def main():
    host, port = "localhost", 5001
    httpd = HTTPServer((host, port), PingHandler)
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        httpd.socket.close()

class PingHandler(BaseHTTPRequestHandler):
    """
    A request handler that always returns "pong".
    """

    def do_GET(self):
        if self.path == "/ping":
            self.send_response(200)
            self.send_header("Content-Type", "text/plain; charset=utf-8")
            self.end_headers()
            self.wfile.write("pong")
            self.wfile.close()
        else:
            self.send_response(404)
            self.send_header("Content-Type", "text/plain; charset=utf-8")
            self.end_headers()
            self.wfile.write("404 page not found")
            self.wfile.close()

if __name__ == "__main__":
    main()
