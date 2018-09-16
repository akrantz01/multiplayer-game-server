YAML
====
The YAML configuration file behind the example

config.yml
----------
Main configuration file

.. code-block:: yaml

    server:
      host: "127.0.0.1"
      port: "5000"
      debug-username: "debug"
      debug-password: "debug"
      upstream:
        active: no
        override-root: no
        locations:
