from flask import Flask, render_template, session
from flask_socketio import SocketIO
from eventlet import monkey_patch

data = {}

app = Flask(__name__)
app.config["SECRET_KEY"] = ""
monkey_patch()
socketio = SocketIO(app)


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
    socketio.run(app)
