API Documentation
=================
This documentation is for the purpose of building the API, it is not currently functional.

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

Users
-----
``GET /users``
^^^^^^^^^^^^^^
Returns a list of all users who have an account on the panel.

.. code-block:: curl

  curl -X GET -i -H "X-Access-Token: demo1111-2222-3333-4444-55556666" https://example.com/api/users

.. code-block:: json

  {
     "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx": {
       "id": 1,
       "username": "demoaccount",
       "email": "some@example.com",
       "admin": 1
     },
     "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyy": {
       "id": 2,
       "username": "demoaccount2",
       "email":"two@example.com",
       "admin": 0
     }
  }

``GET /users/[:uuid]``
^^^^^^^^^^^^^^^^^^^^^^
Returns information about the requested user.

.. code-block:: curl

  curl -X GET -i -H "X-Access-Token: demo1111-2222-3333-4444-55556666" https://example.com/api/users/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx

.. code-block:: json

  {
    "id": 1,
    "username": "demoaccount",
    "email": "some@example.com",
    "admin": 1,
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

.. code-block

  https/1.x 200 OK

Servers
-------
``GET /servers``
^^^^^^^^^^^^^^^^
Returns a list of all servers that are on the system.

.. code-block:: curl

  curl -X GET -i -H "X-Access-Token: demo1111-2222-3333-4444-55556666" https://example.com/api/servers

.. code-block:: json

    {
        "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa": {
            "id": 1,
            "owner": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx",
            "name": "MY_ADMIN_SERVER",
            "node": 1,
            "active": 1
        },
        "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbb": {
            "id": 2,
            "owner": "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyy",
            "name": "STEVES_SERVER",
            "node": 1,
            "active": 1
        }
    }

``GET /servers/[:hash]``
^^^^^^^^^^^^^^^^^^^^^^^^
Returns information about the requested server.

.. code-block:: curl

  curl -X GET -i -H "X-Access-Token: demo1111-2222-3333-4444-55556666" https://example.com/api/servers/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa

.. code-block:: json

    {
        "id": 1,
        "node": 1,
        "owner": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx",
        "name": "MY_ADMIN_SERVER",
        "server_jar": "server.jar",
        "active": 1,
        "ram": 512,
        "disk": 1024,
        "cpu": 30,
        "ip": "192.168.1.2",
        "port": 25565,
        "ftp_user": "mc-MY_DEMO_XyZab"
    }

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
``GET /nodes``
^^^^^^^^^^^^^^^^^^^^
Returns a list of all nodes that are on the system.

.. code-block:: curl

  curl -X GET -i -H "X-Access-Token: demo1111-2222-3333-4444-55556666" https://example.com/api/nodes

.. code-block:: json

    {
        "1": {
            "node": "My_First_Node",
            "fqdn": "ec2-255-255-255-0.us-west-2.compute.amazonaws.com",
            "ip": "255.255.255.0",
            "ports": {
                "255.255.255.1": {
                    "ports_free": 5
                },
                "255.255.255.2": {
                    "ports_free": 3
                },
                "255.255.255.3": {
                    "ports_free": 6
                }
            }
        },
        "2": {
            "node": "My_Second_Node",
            "fqdn": "192.168.1.1",
            "ip": "192.168.1.1",
            "ports": {
                "192.168.1.1": {
                    "ports_free": 5
                }
            }
        }
    }

``GET /nodes/[:id]``
^^^^^^^^^^^^^^^^^^^^
Returns information about the requested node.

.. code-block:: curl

  curl -X GET -i -H "X-Access-Token: demo1111-2222-3333-4444-55556666" https://example.com/api/nodes/1

.. code-block:: json

    {
        "id": 1,
        "node": "My_First_Node",
        "fqdn": "ec2-255-255-255-0.us-west-2.compute.amazonaws.com",
        "ip": "255.255.255.0",
        "gsd_listen": 8003,
        "gsd_console": 8031,
        "gsd_server_dir": "/home/",
        "ports": {
            "255.255.255.1": {
                "25565": 1,
                "25566": 1,
                "25567": 1
            },
            "255.255.255.2": {
                "25565": 1,
                "25566": 1,
                "25567": 1,
                "25568": 0,
                "25569": 1,
                "25570": 1
            },
            "255.255.255.3": {
                "25565": 1,
                "25566": 1,
                "25567": 1,
                "25568": 1,
                "25569": 1,
                "25570": 1
            }
        },
        "servers": [
            "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
        ]
    }

``POST /nodes``
^^^^^^^^^^^^^^^^^^
Creates a new node based on data sent in a JSON request.

.. code-block:: json

    {
        "node": "My_Third_Node",
        "ip": "10.0.1.1",
        "ips": "10.0.1.1|25565-25580\n10.0.1.2|25565,25570-25580,25590\n10.0.1.2|25565",
        "(OPTIONAL) fqdn": "example.com",
        "(OPTIONAL) gsd_listen": 8003,
        "(OPTIONAL) gsd_console": 8031,
        "(OPTIONAL) gsd_server_dir": "/home/",
    }

.. code-block:: curl

  curl -X POST -i \
    -H "X-Access-Token: demo1111-2222-3333-4444-55556666" \
    -H "Content-Type: application/json" \
    -d '{"node": "My_Third_Node","ip": "10.0.1.1","ips": "10.0.1.1|25565-25580\n10.0.1.2|25565,25570-25580,25590\n10.0.1.2|25565"}'
    https://example.com/api/nodes

.. code-block:: json

  {
    "id": 3,
    "node": "My_Third_Node"
  }


``PUT /nodes/[:id]``
^^^^^^^^^^^^^^^^^^^^
Updates node information.
