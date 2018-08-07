Client
======
The part of the application that allows for connection to the server

Usage
-----
.. js:class:: mmmoc.MMOC(add_depends=false)

    The MMOC object is the basis for connecting to the server. It acts as the centerpoint for all multiplayer actions, such as sending data and getting data.

    :param bool add_depends: Automatically write dependencies to the document, default to false

    .. js:method:: init(id_len=8, wsurl="//" + document.domain + ":" + location.port)

        Initialize the player's ID and connection to the main server

        :param number id_len: Length of the id to be generated
        :param string wsurl: WebSockets URL to connect to

    .. js:method:: sendData()

        Send the player's coordinates, ID and other data to the server

    .. js:method:: getDataFromServer()

        Get every players' data from the server

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

Example
-------
.. code-block:: javascript

    let mp = new MMOC(true);
    mp.init();

    while true:
        mp.sendData();
        mp.getDataFromServer();

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
