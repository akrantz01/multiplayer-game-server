Server
======
The part of the application that allows for communication between the various clients.

Usage
-----
::

    $ ./server --help
    Usage of ./server:
      -config string
            Alternative yaml configuration file (default "config.yml")

Default Values
--------------
* **Config** -> ``"config.yml"``

Build from Source
-----------------
Build a binary from the master branch. Use the following commands::

    go get
    go build -o server

Find the executable file called ``server`` in the folder that you ran the command in. Currently this has only been tested on Ubuntu 18.04, but feel free to try to make it work on other platforms.
