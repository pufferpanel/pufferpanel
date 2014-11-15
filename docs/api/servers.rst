API: Servers
============

``GET`` /servers
----------------
Returns a list of all servers that are on the system.

Request
^^^^^^^
.. code-block:: curl

	curl -X GET -i -H "X-Access-Token: ABCDEFGH-1234-5678-0000-abcdefgh" https://example.com/api/servers

Response
^^^^^^^^
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

``GET`` /servers/[:hash]
------------------------
Returns information about the requested server.

Parameters
^^^^^^^^^^
+--------+------------+-----------+------------------------------------------------------------+
| Method | Parameter  | Required  | Description                                                |
+========+============+===========+============================================================+
| GET    | hash       | yes       | The hash of the server that you wish to return data about. |
+--------+------------+-----------+------------------------------------------------------------+

Request
^^^^^^^
.. code-block:: curl

	curl -X GET -i -H "X-Access-Token: ABCDEFGH-1234-5678-0000-abcdefgh" https://example.com/api/servers/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa

Response
^^^^^^^^
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

``POST`` /servers
-----------------

``PUT`` /servers/[:hash]
------------------------

``DELETE``  /servers/[:hash]
----------------------------