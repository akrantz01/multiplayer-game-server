# Massively Multiplayer Online Game Server
### NodeJS Edition
A simple massively multiplayer online game server re-written in NodeJS for my Advanced Animation teacher so he does not need to learn [Golang](https://golang.org).

## Running
By default the server will run on [localhost:8080](http://localhost:8080), but can be configured to run on other hosts and ports with the environment variables `MMOS_HOST` and `MMOS_PORT`.
```
# Run on defaults
node index.js

# Run on different host
export MMOS_HOST=0.0.0.0
node index.js

# Run on different port
export MMOS_PORT=80
node index.js

# Run on different host and port
export MMOS_HOST=0.0.0.0
export MMOS_PORT=80
node index.js
```
**NOTE**: if you are on Windows, use `set` instead of `export`
