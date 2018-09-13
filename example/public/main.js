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
