API: Nodes
==========

``GET`` /nodes
--------------
Returns a list of all nodes that are on the system.

Request
^^^^^^^
.. code-block:: curl

	curl -X GET -i -H "X-Access-Token: ABCDEFGH-1234-5678-0000-abcdefgh" https://example.com/api/nodes

Response
^^^^^^^^
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

``GET`` /nodes/[:id]
--------------------
Returns information about the requested node.

Parameters
^^^^^^^^^^
+--------+------------+-----------+------------------------------------------------------------+
| Method | Parameter  | Required  | Description                                                |
+========+============+===========+============================================================+
| GET    | id         | yes       | The id of the node that you wish to return data about.     |
+--------+------------+-----------+------------------------------------------------------------+

Request
^^^^^^^
.. code-block:: curl

	curl -X GET -i -H "X-Access-Token: ABCDEFGH-1234-5678-0000-abcdefgh" https://example.com/api/nodes/1

Response
^^^^^^^^
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

``POST`` /nodes
---------------
Creates a new node based on data sent in a JSON request.

Parameters
^^^^^^^^^^
+----------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Parameter      | Optional | Description                                                                                                                                                                 |
+================+==========+=============================================================================================================================================================================+
| node           |          | The name of the node you are creating.                                                                                                                                      |
+----------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| ip             |          | The IP for the node. Can be local, but it is suggested to use a public IP to prevent any connection issues.                                                                 |
+----------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| ips            |          | A list of IPs to add to the server. They should be sent with a newline character between each set. Please see the example in the Admin CP for how to string these together. |
+----------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| fqdn           | yes      | The Fully Qualified Domain Name for the node you are adding. If not specified or invalid defaults to the server IP.                                                         |
+----------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| gsd_listen     | yes      | The port that GSD will be listening on. (Default: ``8003``)                                                                                                                 |
+----------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| gsd_console    | yes      | The port that the GSD console will be listening on. (Default: ``8031``)                                                                                                     |
+----------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| gsd_server_dir | yes      | The folder where you would like servers to be created. (Default: ``/home/``)                                                                                                |
+----------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------+

Request
^^^^^^^
.. code-block:: curl

	curl -X POST -i \
		-H "X-Access-Token: ABCDEFGH-1234-5678-0000-abcdefgh" \
		-H "Content-Type: application/json" \
		-d '{"node": "My_Third_Node","ip": "10.0.1.1","ips": "10.0.1.1|25565-25580\n10.0.1.2|25565,25570-25580,25590\n10.0.1.2|25565"}'
		https://example.com/api/nodes

Response
^^^^^^^^
+------------+------------------------------------+
| Parameter  | Description                        |
+============+====================================+
| id         | The ID of the newly created node.  |
+------------+------------------------------------+
| node       | The name of the newly created node.|
+------------+------------------------------------+
.. code-block:: json

	{
		"id": 3,
		"node": "My_Third_Node"
	}


``PUT /nodes/[:id]``
^^^^^^^^^^^^^^^^^^^^
Updates node information.