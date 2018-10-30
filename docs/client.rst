Client
======
The part of the application that allows for connection to the server

Usage
-----
.. js:class:: mmoc.MMOC()

    The MMOC object is the basis for connecting to the server. It acts as the centerpoint for all multiplayer actions, such as sending data and getting data.

    .. js:method:: init(id_len=8, wsurl="//" + document.domain + ":" + location.port)

        Initialize the player's ID and connection to the main server

        :param number id_len: Length of the id to be generated
        :param string wsurl: WebSockets URL to connect to

    .. js:method:: sendPlayerData()

        Send the player's coordinates, ID and other data to the server

    .. js:method:: sendObjectData(object)

        Send a projectile's or any moving object's data that is not a player

        :param MovingObject object: The object that is having its data being sent

    .. js:method:: removeObject(object)

        Remove an object from the server

        :param MovingObject object: The object to be removed

    .. js:method:: getPlayers()

        Get every player as an object

        :return: each player's coordinates, id and other data
        :rtype: object

    .. js:method:: getGlobals()

        Get all of the server globals

        :return: all of the server's global variables set in config.yml
        :rtype: object

    .. js:method:: getObjects()

        Get all of the moving objects

        :return: all of the moving objects
        :rtype: object

    .. js:method:: changeX(by)

        Change the player's X value by a specified value

        :param number by: value to change X by

    .. js:method:: changeY(by)

        Change the player's Y value by a specified value

        :param number by: value to change Y by

    .. js:method:: changeZ(by)

        Change the player's Z value by a specified value

        :param number by: value to change Z by

    .. js:method:: changeOrientation(by)

        Change the player's orientation by a specified value

        :param nubmer by: value to change the orientation by

    .. js:method:: setOther(key, value)

        Set a piece of data to be broadcasted to every other player. The data will be located at the key *key*

        :param string key: key to store the data at
        :param any value: value to be stored

    .. js:method:: isconnected()

        Check if the client has connected to the server

        :return: promise that will resolve to a boolean value
        :rtype: Promise



.. js:class:: mmoc.MovingObject(mesh, p, r)

    The MovingObject class is for creating new objects that need to be moved and synced between clients

    :param mesh object: the mesh/data for the object
    :param p function: a function for getting the position of the mesh. Must modify this.x, this.y, and this.z
    :param r function: a function for rendering the mesh

    .. js:method:: render()

        Render the mesh

    .. js:method:: setOther(key, value)

        Set a property for the object

        :param key string: the key for accessing the object
        :param value object: a json serializable object

    .. js:method:: getOther(key)

        Get a property at the key 'key'

        :param key string: the key for accessing the value
        :return: the value at the specified key
        :rtype: object

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
