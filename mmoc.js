var MMOC = (function() {
    const reqd = (name) => { throw new Error("Expected argument '" + name + "'") };

    let _id = "";
    let _x = 0;
    let _y = 0;
    let _other = {};
    let _data = {};
    let _connected = false;

    class MMOC {
        constructor(add_depends = false) {
            if (add_depends) {
                let script = document.createElement("script");
                script.src = "https://cdnjs.cloudflare.com/ajax/libs/socket.io/2.1.1/socket.io.js";
                document.head.appendChild(script);
            }
        }

        init(id_len = 8, wsurl = "//" + document.domain + ":" + location.port) {
            this.socket = io.connect(wsurl);

            this.socket.on("connect", function() {
                let possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

                for (let i = 0; i < id_len; i++) {
                    _id += possible.charAt(Math.floor(Math.random() * 16));
                }
                _connected = true;
            });
        }

        sendData() {
            this.socket.emit("data", {id: _id, x: _x, y: _y, other: _other});
        }

        getDataFromServer() {
            this.socket.emit("get", function (data) {
                _data = data;
            });
        }

        getData() {
            return _data;
        }

        changeX(by = reqd("by")) {
            _x += by;
        }

        changeY(by = reqd("by")) {
            _y += by;
        }

        setOther(key = reqd("key"), value = reqd("value")) {
            _other[key] = value;
        }

        isconnected() {
            return new Promise(function(resolve, reject) {
                if (_connected) {
                    resolve();
                } else {
                    reject();
                }
            });
        }
    }

    return MMOC;
}());
