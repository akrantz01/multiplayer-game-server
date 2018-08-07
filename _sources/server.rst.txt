Server
======
The part of the application that allows for communication between the various clients.

Usage
-----
::

    $ ./server --help
    usage: server [-h] [--host HOST] [--port PORT] [--debug] [--template TEMPLATE]

    optional arguments:
        -h, --help              show this help message and exit
        --host HOST             host to run the server on
        --port PORT             port to run the server on
        --debug                 run the server in debug mode
        --template TEMPLATE     template file to launch from

Default Value
-------------
* **Host** -> ``"127.0.0.1"``
* **Port** -> ``5000``
* **Debug** -> ``False``
* **Template** -> ``"index.html"``

Build from Source
-----------------
Build a binary from the master branch. Use the following commands::

    pip3 install -r requirements.txt
    pyinstaller --onefile server.py --hidden-import=dns --hidden-import=dns.dnssec --hidden-import=dns.e164 --hidden-import=dns.namedict --hidden-import=dns.tsigkeyring --hidden-import=dns.update --hidden-import=dns.version --hidden-import=dns.zone --hidden-import=engineio.async_eventlet

Find the executable file in the ``/dist`` folder. PyInstaller should generate the proper executable based on your system's OS.
