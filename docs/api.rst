PufferPanel API Documentation
=============================
This documentation is for the purpose of building the API, it is not currently functional.

Authentication
--------------
Requests to the PufferPanel API should ideally be made over a secure (HTTPS) connection to prevent man in the middle attacks, and protect sensitive client data.
Requests must also include an API key sent in the request using the ``X-Authorization`` header.

Errors & Status Codes
^^^^^^^^^^^^^^^^^^^^^
``200 OK`` - Everything worked as expected.

``400 Bad Request`` - Missing a required parameter.

``401 Unauthorized`` - No valid API key provided.

``402 Request Failed`` - Parameters were valid but request failed.

``404 Not Found`` - The requested item doesn't exist.

``500 Server Error`` - Something went wrong on the server.

Users
-----
``GET /users/[:uuid]``
^^^^^^^^^^^^^^^^^^^^^^
Returns information about the requested user ID.

.. code-block:: curl
  
  curl -X GET -i -H "X-Authorization: demo-1233-23445566-2343" http://example.com/api/users/1
  
.. code-block:: json

  HTTP/1.x 200 OK
  {
    "id": 1,
    "email": some@example.com
  }
    

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
