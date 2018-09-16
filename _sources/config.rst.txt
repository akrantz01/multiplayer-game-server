.. role:: yaml(code)
    :language: yaml

Configuration
=============
For configuring the server with its different parameters

Descriptions
------------
.. code-block:: yaml

    server:                                 # server config

      host: "127.0.0.1"                     # host
      port: "5000"                          # port

      debug-username: "debug"               # /debug username
      debug-password: "debug"               # /debug password

      upstream:                             # proxy settings
        active: yes                         # enable or disable
        override-root: yes                  # override / endpoint
        locations:                          # proxy from host to endpoint
          - url: "http://localhost:8080"    # service url
            endpoint: "/"                   # endpoint location

    globals:                                # define global variables

      s:                                    # section name

        value:                              # value name
          value: "abc"                      # actual value
          type: "string"                    # value type

        value2:                             # value name
          value: "abc2"                     # actual value
          type: "string"                    # value type

To view a config file that is being used, see config.yml in :doc:`the example <example/yml>`.

Special Cases
-------------
These are some cases where typical keys only being used once do not apply.

Server -> Upstream -> Locations:
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
    * Acts as an array of urls to proxy
    * Requires both :yaml:`url` and :yaml:`endpoint` to work

Globals:
^^^^^^^^
    * Acts as an array of values
    * Can have n amount of values
    * Each value must be named
    * Each value must have a section
