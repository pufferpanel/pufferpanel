API Documentation
=================

Authentication
--------------
Requests to the PufferPanel API should ideally be made over a secure (HTTPS) connection to prevent man in the middle attacks, and protect sensitive client data.
Requests must also include an API key sent in the request using the ``X-Access-Token`` header.

Errors & Status Codes
^^^^^^^^^^^^^^^^^^^^^
``200 OK`` - Everything worked as expected.

``400 Bad Request`` - Missing a required parameter.

``401 Unauthorized`` - No valid API key provided.

``403 Forbidden`` - Request must be made using a secure connection.

``404 Not Found`` - The requested item doesn't exist.

``409 Conflict`` - Conflict occured in the request.

``500 Server Error`` - Something went wrong on the server.