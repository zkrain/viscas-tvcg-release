from flask import Flask, request
import json
from flask_cors import CORS
from fastdtw import fastdtw

app = Flask(__name__)
CORS(app)


@app.route('/')
def helloWorld():
    return 'Hello, World!'


@app.route('/fastdtw', methods=['POST', 'GET'])
def fastDtw():
    data = request.get_json(silent=True)
    x = data['x']
    y = data['y']
    distance, path = fastdtw(x, y)
    return json.dumps({'distance': distance})


if __name__ == '__main__':
    app.run()
