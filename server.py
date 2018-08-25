from argparse import ArgumentParser
from configparser import ConfigParser, NoSectionError, NoOptionError
from eventlet import monkey_patch
from flask import Flask, render_template, session, request, jsonify
from flask_socketio import SocketIO
from os import urandom

app = Flask(__name__)
app.config["SECRET_KEY"] = urandom(16)
monkey_patch()
socketio = SocketIO(app)
data = {}

parser = ArgumentParser()
parser.add_argument("config", help="config file to load from. NOTE: the config file takes prescience over arguments")
args = parser.parse_args()
config = ConfigParser()
config.read(args.config)

if "server" not in config.sections():
    raise NoSectionError("Unable to find section 'server'")
host = config["server"].get("host")
port = config["server"].getint("port")
debug = config["server"].getboolean("debug")

globals = [s for s in config.sections() if s != "server"]
if not globals:
    data["globals"] = None
else:
    data["globals"] = {}

for s in globals:
    keys = [k for k in config[s].keys() if "-" not in k]
    if s != "unorganized":
        data["globals"][s] = {}
    for key in keys:
        try:
            t = eval(config[s][key + "-type"])
            if s != "unorganized":
                data["globals"][s][key] = t(config[s][key])
            else:
                data["globals"][key] = t(config[s][key])
        except KeyError:
            raise NoOptionError(key+"-type", "Argument type required for '{0}' in section '{1}'".format(key, s))


@app.route('/')
def index():
    return render_template("index.html")


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
    socketio.run(app, host=host, port=port, debug=debug)
