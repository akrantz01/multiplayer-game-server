let MMOC = (function() {
    const reqd = (name) => { throw new Error("Expected argument '" + name + "'") };

    let _id = "";
    let _x = 0;
    let _y = 0;
    let _z = 0;
    let _orientation = 0;
    let _other = {};
    let _data = {};
    let _connected = false;

    class MMOC {
        init(id_len=8, wsurl="//" + document.domain + ":" + location.port + "/ws") {
            if (location.protocol === "https") wsurl = "wss:" + wsurl;
            else wsurl = "ws:" + wsurl;
            this.ws = new WebSocket(wsurl);

            this.ws.onopen = function (event) {
                let possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

                for (let i = 0; i < id_len; i++) {
                    _id += possible.charAt(Math.floor(Math.random() * possible.length));
                }
                _connected = true;
            };

            this.ws.onmessage = function (event) {
                _data = JSON.parse(event.data);
            };

            setInterval(() => {
                this.ws.send(JSON.stringify({
                    type: 2
                }));
            }, 15);
        }

        sendData() {
            this.ws.send(JSON.stringify({
                type: 1,
                id: _id,
                other: _other,
                coordinates: {
                    x: _x,
                    y: _y,
                    z: _z
                }
                orientation: _orientation
            }));
        }

        getData() {
            return _data;
        }

        changeX(by=reqd("by")) {
            _x += by;
        }

        changeY(by=reqd("by")) {
            _y += by;
        }

        changeZ(by=reqd("by")) {
            _z += by;
        }

        changeOrientation(by=reqd("by")) {
            _orientation += by;
        }

        setOther(key=reqd("key"), value=reqd("value")) {
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
})();