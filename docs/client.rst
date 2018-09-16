Client
======
The part of the application that allows for connection to the server

Usage
-----
.. js:class:: mmmoc.MMOC()

    The MMOC object is the basis for connecting to the server. It acts as the centerpoint for all multiplayer actions, such as sending data and getting data.

    .. js:method:: init(id_len=8, wsurl="//" + document.domain + ":" + location.port)

        Initialize the player's ID and connection to the main server

        :param number id_len: Length of the id to be generated
        :param string wsurl: WebSockets URL to connect to

    .. js:method:: sendData()

        Send the player's coordinates, ID and other data to the server

    .. js:method:: getData()

        Get every players' as an object

        :return: each player's coordinates, id and other data
        :rtype: object

    .. js:method:: changeX(by)

        Change the player's X value by a specified value

        :param number by: value to change X by

    .. js:method:: changeY(by)

        Change the player's Y value by a specified value

        :param number by: value to change Y by

    .. js:method:: setOther(key, value)

        Set a piece of data to be broadcasted to every other player. The data will be located at the key *key*

        :param string key: key to store the data at
        :param any value: value to be stored

    .. js:method:: isconnected()

        Check if the client has connected to the server

        :return: promise that will resolve to a boolean value
        :rtype: Promise

Example
-------
.. code-block:: javascript

    let mp = new MMOC();
    mp.init();

    while (true):
        if(...) {
            mp.changeX(1);
        }
        if(...) {
            mp.changeX(-1);
        }
        if(...) {
            mp.changeY(1);
        }
        if(...) {
            mp.changeY(-1);
        }

        let d = mp.getData();
        drawToScreenFunction(d);

To view a more flushed out example see main.js in :doc:`the example <example/js>`.
