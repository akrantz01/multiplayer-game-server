JavaScript
==========
The JavaScript behind the example

main.js
-------
Primary game code

.. code-block:: javascript

    let player, positions;
    let multiplayer = new MMOC();
    multiplayer.setOther("color", "#ffffff");

    function update(jscolor) {
        multiplayer.setOther("color", "#" + jscolor);
    }

    function setup() {
        multiplayer.init();

        multiplayer.changeX(Math.floor(windowWidth/2));
        multiplayer.changeY(Math.floor(windowHeight/2));

        document.getElementById("modal").style.display = "block";
        document.getElementById("observer").onclick = function () {
            document.getElementById("modal").style.display = "none";
            player = false;
        };
        document.getElementById("player").onclick = function () {
            document.getElementById("modal").style.display = "none";
            player = true;
        };

        let canvas = createCanvas(windowWidth, windowHeight);
        canvas.style("display", "block");
    }

    function draw() {
        background(255);
        if (player) {
            multiplayer.sendData();
            if (keyIsDown(87)) {
                multiplayer.changeY(-1);
            }
            if (keyIsDown(83)) {
                multiplayer.changeY(1);
            }
            if (keyIsDown(65)) {
                multiplayer.changeX(-1);
            }
            if (keyIsDown(68)) {
                multiplayer.changeX(1);
            }
        }

        let data = multiplayer.getData();
        for (let key in data["Users"]) {
            fill(data["Users"][key]["Other"]["color"]);
            ellipse(data["Users"][key]["X"], data["Users"][key]["Y"], 20);
        }
    }

mmoc.js
-------
MMOS client

.. code-block:: javascript

    let MMOC = (function() {
        const reqd = (name) => { throw new Error("Expected argument '" + name + "'") };

        let _id = "";
        let _x = 0;
        let _y = 0;
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
                        y: _y
                    }
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
