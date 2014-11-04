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

Single Sign On
--------------
``POST /sso``
^^^^^^^^^^^^^
Allows you to use PufferPanel as a single sign on system. Posting an email and password returns either a HTTP error, or a JSON string containing the users login token.
In order to allow a seamless login to the panel you will need to set a cookie on your end named ``pp_auth_token`` with a value of the token returned.

Setting the ``sso`` value to be ``false`` allows for you to simply check if the email and password combination is valid, it does not return any JSON data, only a HTTP status code.

.. code-block:: json

  {
    "email": "some@example.com",
    "password": "somepassword",
    "sso": true
  }

.. code-block:: curl

  curl -X POST -i \
    -H "X-Authorization: demo1111-2222-3333-4444-55556666" \
    -H "Content-Type: application/json" \
    -d '{"email":"some@example.com","password":"somepassword","sso":true}'
    http://example.com/api/sso
  
.. code-block:: json

  {
    "token": "XyZ"
  }

Users
-----
``GET /users``
^^^^^^^^^^^^^^
Returns a list of all users who have an account on the panel.

.. code-block:: curl

  curl -X GET -i -H "X-Authorization: demo1111-2222-3333-4444-55556666" http://example.com/api/users
  
.. code-block:: json

  {
    "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx": {
      "id": 1,
      "username": "demoaccount",
      "email": "some@example.com"
    },
    "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyy": {
      "id": 2,
      "username": "demoaccount2",
      "email":"two@example.com"
    }
  }

``GET /users/[:uuid]``
^^^^^^^^^^^^^^^^^^^^^^
Returns information about the requested user.

.. code-block:: curl
  
  curl -X GET -i -H "X-Authorization: demo1111-2222-3333-4444-55556666" http://example.com/api/users/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx
  
.. code-block:: json

  {
    "id": 1,
    "username": "demoaccount",
    "email": "some@example.com",
    "servers": [
      "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa",
      "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbb"
    ]
  }
    

``POST /users``
^^^^^^^^^^^^^^^
Creates a new user based on data sent in a JSON request.

``PUT /users/[:uuid]``
^^^^^^^^^^^^^^^^^^^^^^
Updates user information.

``DELETE  /users/[:uuid]``
^^^^^^^^^^^^^^^^^^^^^^^^^^
Deletes a user given a specified ID.

.. code-block:: curl

  curl -X DELETE -i -H "X-Authorization: demo1111-2222-3333-4444-55556666" http://example.com/api/users/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx
  
.. code-block

  HTTP/1.x 200 OK

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
