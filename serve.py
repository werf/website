import http.server
import socketserver

PORT = 3001

class MyHttpRequestHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.path = self.path.replace('applications_guide_ru/','applications_guide_ru/_site/')
        self.path = self.path.replace('applications_guide_en/','applications_guide_en/_site/')
        return http.server.SimpleHTTPRequestHandler.do_GET(self)

# Create an object of the above class
handler = MyHttpRequestHandler

with socketserver.TCPServer(("", PORT), handler) as httpd:
    print("Server started at localhost:" + str(PORT))
    httpd.serve_forever()