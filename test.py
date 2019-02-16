# test.py test for image-uploader
import http.client

host = 'localhost:8080'

conn = http.client.HTTPConnection(host)
conn.request("GET", "/")

response = conn.getresponse()
print("{} {}\n{}\n".format(response.status, response.reason, response.read()))
conn.close()

