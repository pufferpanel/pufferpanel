Installing Nodes
================
This documentation was written as we performed an installation of this program on Ubuntu 14.04 LTS. Please be aware that the command structure may vary depending on your *nix distro.

Dependencies
------------
* CURL
* OpenSSL
* Git
* ``gcc`` and ``make``
* Java (``apt-get install default-jre``)
* NodeJS

Installing NodeJS
^^^^^^^^^^^^^^^^^
Please follow the instructions below for adding NodeJS (thanks to `@gametainers/gsd <https://github.com/gametainers/gsd/>`_ for the documentation).

.. code-block:: sh

	sudo add-apt-repository ppa:chris-lea/node.js
	sudo apt-get update
	sudo apt-get install nodejs make g++ git

Create a Node
-------------
The first step in this process requires that you first add a node in the panel. Once you have done this you can continue with this process.

Initial Setup
-------------
Before we can begin installing a node, you must first create some directories that are needed by GSD.

.. code-block:: sh

	[$]~ mkdir -p /srv/gsd
	[$]~ mkdir -p /mnt/MC/CraftBukkit

You will need to upload a file called ``server.jar`` into the ``/mnt/MC/CraftBukkit`` directory. You can also upload any other miscelaneous files that you want to be added to each server upon creation.

Install Game Server Daemon
--------------------------

.. code-block:: sh

	[$]~ cd /srv
	[$]~ git clone https://github.com/gametainers/gsd.git

This will download all of the files necessary and place them into the correct directory. Please edit ``config.json.example`` to be ``config.json``.

Your ``config.json`` file should look similar to the code below. Please make sure to update the ```authurl`` value to link to your panel correctly. Update ``YOUR_NODE_TOKEN_HERE`` to be the code that is displayed in PufferPanel for your node, it is called ``gsd_secret`` in the panel.

.. code-block:: json

	{
	    "daemon": {
	        "listenport": 8003,
	        "consoleport": 8031
	    },
	    "tokens": [
			"YOUR_NODE_TOKEN_HERE"
		],
	    "interfaces":{
	      "rest":{},
	      "console":{},
	      "ftp":{
		    "authurl": "http://www.example.com/ajax/validate_ftp.php",
		    "port": 21,
		    "host": "127.0.0.1"
	      }
	    },
	    "servers": []
	}

After editing your config file, open ``gsd.js`` and modify it so that it looks like the code below.

.. code-block:: js

	var config = require('./config.json');

	require('./interfaces/console.js');
	require('./interfaces/rest.js');
	require('./interfaces/ftp.js');

	var servers = require('./services');

Once all of that is complete run the commands below to start GSD.

.. code-block:: sh

	[$]~ cd /srv/gsd
	[$]~ npm install
	[$]~ npm start

Congratulations! Your first node is configured.
