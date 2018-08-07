# Massively Multiplayer Online Game Server
A simple massively multiplayer online game server written in Python 3.6 and designed for use in JavaScript ES6.

## Get the Server
### Option 1
Download it from the releases tab. [Quick link](https://github.com/akrantz01/mmos/releases/latest) to the most recent release.
### Option 2
Build a binary from the master branch. Use the following commands:
```text
pip3 install -r requirements.txt
pyinstaller --onefile server.py --hidden-import=dns --hidden-import=dns.dnssec --hidden-import=dns.e164 --hidden-import=dns.namedict --hidden-import=dns.tsigkeyring --hidden-import=dns.update --hidden-import=dns.version --hidden-import=dns.zone --hidden-import=engineio.async_eventlet
```
Find the executable file in the `/dist` folder. PyInstaller should generate the proper executable based on your system's OS. 

## Simple Example
You can find an example in the [/example](/example) directory.
