PufferPanel API Documentation
=============================
This documentation is for the purpose of building the API, it is not currently functional.

Users
-----
GET /[:key]/users/[:uuid]
^^^^^^^^^^^^^^^^^^^^^^^^
Returns information about the requested user ID.

POST /[:key]/users
^^^^^^^^^^^^^^^^^^
Creates a new user based on data sent in a JSON request.

PUT /[:key]/users/[:uuid]
^^^^^^^^^^^^^^^^^^^^^^^
Updates user information.

DELETE  /[:key]/users/[:hash]
^^^^^^^^^^^^^^^^^^^^^^^^^^^
Deletes a server given a specified ID.

Servers
-------
GET /[:key]/servers/[:hash]
^^^^^^^^^^^^^^^^^^^^^^^^
Returns information about the requested server.

POST /[:key]/servers
^^^^^^^^^^^^^^^^^^
Creates a new server based on data sent in a JSON request.

PUT /[:key]/servers/[:hash]
^^^^^^^^^^^^^^^^^^^^^^^
Updates server information.

DELETE  /[:key]/servers/[:hash]
^^^^^^^^^^^^^^^^^^^^^^^^^^^
Deletes a server given a specified hash.
