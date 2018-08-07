from argparse import ArgumentParser
from eventlet import monkey_patch
from flask import Flask, render_template, session
from flask_socketio import SocketIO
from os import urandom

parser = ArgumentParser()
parser.add_argument("--host", help="host to run the server on", default="127.0.0.1")
parser.add_argument("--port", help="port to run the server on", default=5000)
parser.add_argument("--debug", help="run the server in debug mode", action="store_true")
parser.add_argument("--template", help="template file to load from", default="index.html")
args = parser.parse_args()

app = Flask(__name__)
app.config["SECRET_KEY"] = urandom(16)
monkey_patch()
socketio = SocketIO(app)
data = {}


@app.route('/')
def index():
    return render_template(args.template)


@socketio.on("data")
def on_data(json):
    session["id"] = json.get("id")
    data[json.get("id")] = {"x": json.get("x"), "y": json.get("y"), "other": json.get("other")}


@socketio.on("get")
def on_get():
    return data


@socketio.on("disconnect")
def on_disconnect():
    try:
        data.pop(session.get("id"))
    except KeyError:
        pass


if __name__ == '__main__':
    socketio.run(app, host=args.host, port=args.port, debug=args.debug)
