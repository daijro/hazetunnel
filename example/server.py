"""
This is an example server to test the capabilities of hazetunnel.

Example command for testing HTML injection,
Assuming hazetunnel is running on 8080, and this server on 5000:

curl --proxy http://localhost:8080 \
    --insecure http://localhost:5000/html \
    -H "x-mitm-payload: alert('Hello world');" \
    -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"

Example command for testing JavaScript injection:

curl --proxy http://localhost:8080
    --insecure http://localhost:5000/js \
    -H "x-mitm-payload: alert('Hello world');" \
    -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"
"""

from base64 import b64encode

from flask import Flask, Response

app = Flask(__name__)

# JavaScript script to be served
js_script = "console.log('Original JavaScript executed.');"
# Encode the JavaScript script in base64
encoded_js_script = b64encode(js_script.encode()).decode('utf-8')

# HTML content with an embedded base64 JavaScript blob
html_content = f"""
<!DOCTYPE html>
<html>
<head>
    <title>Testing Page</title>
</head>
<body>
    <h1>Base64 JavaScript Testing Page</h1>
    <p>This page includes an embedded base64 encoded JavaScript for testing.</p>
    <!-- Embedding the JavaScript directly using a data URI scheme -->
    <script src="data:application/javascript;base64,{encoded_js_script}"></script>
</body>
</html>
"""


@app.route('/html')
def home():
    return html_content


@app.route('/js')
def serve_js():
    return Response(js_script, mimetype='application/javascript')


if __name__ == '__main__':
    app.run(debug=True)
