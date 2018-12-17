let MMOC = (function() {
    const reqd = (name) => { throw new Error("Expected argument '" + name + "'") };

    let _id = "";
    let _x = 0;
    let _y = 0;
    let _z = 0;
    let _orientation = 0;
    let _other = {};
    let _data = {};
    let _broadcasts = {};
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
                let d = JSON.parse(event.data);
                if (d.Type === 4) _broadcasts[d.ID] = d;
                _data = d;
            };
        }

        sendPlayerData() {
            this.ws.send(JSON.stringify({
                type: 1,
                id: _id,
                other: _other,
                coordinates: {
                    x: _x,
                    y: _y,
                    z: _z
                },
                orientation: _orientation
            }));
        }

        sendObjectData(object=reqd('object')) {
            this.ws.send(JSON.stringify({
                type: 2,
                id: object.id,
                other: object.other,
                coordinates: {
                    x: object.x,
                    y: object.y,
                    z: object.z
                }
            }));
        }

        removeObject(object=reqd('object')) {
            this.ws.send(JSON.stringify({
                type: 3,
                id: object.id
            }));
        }

        broadcast(id=reqd('id'), x=reqd('x'), y=reqd('y'), z=reqd('z')) {
            this.ws.send(JSON.stringify({
                type: 4,
                id: id,
                coordinates: {
                    x: x,
                    y: y,
                    z: z
                }
            }))
        }

        getData() {
            return _data;
        }

        getPlayers() {
            return _data["Users"];
        }

        getGlobals() {
            return _data["Globals"];
        }

        getObjects() {
            return _data["Objects"];
        }

        getBroadcasts() {
            setTimeout(() => {_broadcasts = {}}, 1);
            return _broadcasts;
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

        setPosition(x=reqd('x'), y=reqd('y'), z=reqd('z')) {
            _x = x;
            _y = y;
            _z = z;
        }

        getID() {
            return _id;
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

class MovingObject {
    constructor(mesh, p, r) {
        this.id = (function () {
            function s4() {
                return Math.floor((1 + Math.random()) * 0x1000).toString(16).substring(1);
            }
            return s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4() + s4() + s4();
        })();
        this.mesh = mesh;

        this.x = 0;
        this.y = 0;
        this.z = 0;

        this.other = {};

        this.p = p;
        this.r = r;
    }

    render() {
        this.p();
        this.r();
    }

    setOther(key, value) {
        this.other[key] = value;
    }

    getOther(key) {
        return this.other[key];
    }
}
