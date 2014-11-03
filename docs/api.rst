PufferPanel API Documentation
=============================
This documentation is for the purpose of building the API, it is not currently functional.

Authentication
--------------
Requests to the PufferPanel API should ideally be made over a secure (HTTPS) connection to prevent man in the middle attacks, and protect sensitive client data.
Requests must also include an API key sent in the request using the ``X-Authorization`` header.

Users
-----
``GET /users/[:uuid]``
^^^^^^^^^^^^^^^^^^^^^^
Returns information about the requested user ID.

``POST /users``
^^^^^^^^^^^^^^^
Creates a new user based on data sent in a JSON request.

``PUT /users/[:uuid]``
^^^^^^^^^^^^^^^^^^^^^^
Updates user information.

``DELETE  /users/[:hash]``
^^^^^^^^^^^^^^^^^^^^^^^^^^
Deletes a server given a specified ID.

Servers
-------
``GET /servers/[:hash]``
^^^^^^^^^^^^^^^^^^^^^^^^
Returns information about the requested server.

``POST /servers``
^^^^^^^^^^^^^^^^^
Creates a new server based on data sent in a JSON request.

``PUT /servers/[:hash]``
^^^^^^^^^^^^^^^^^^^^^^^^
Updates server information.

``DELETE  /servers/[:hash]``
^^^^^^^^^^^^^^^^^^^^^^^^^^^^
Deletes a server given a specified hash.

Nodes
-----
``GET /nodes/[:id]``
^^^^^^^^^^^^^^^^^^^^
Returns information about the requested node.

``POST /nodes``
^^^^^^^^^^^^^^^^^^
Creates a new node based on data sent in a JSON request.

``PUT /nodes/[:id]``
^^^^^^^^^^^^^^^^^^^^
Updates node information.
