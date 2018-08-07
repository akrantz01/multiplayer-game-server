JavaScript
==========
The JavaScript behind the example

main.js
-------
Primary game code

.. code-block:: javascript

    let player, positions;
    let multiplayer = new MMOC(true);
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

        multiplayer.getDataFromServer();
        let data = multiplayer.getData();
        for (let key in data) {
            fill(data[key]["other"]["color"]);
            ellipse(data[key]["x"], data[key]["y"], 20);
        }
    }

mmoc.js
-------
MMOS client

.. code-block:: javascript

    var MMOC = (function() {
        const reqd = (name) => { throw new Error("Expected argument '" + name + "'") };

        let _id = "";
        let _x = 0;
        let _y = 0;
        let _other = {};
        let _data = {};

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
        }

        return MMOC;
    }());
