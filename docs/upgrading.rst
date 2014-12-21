Upgrading PufferPanel
=====================
Starting in version ``0.7.5 Beta`` we now support running an upgrader rather than having to re-install the entire panel. Please be aware that this upgrader is still very much beta software, and may still cause issues. Use at your own risk.

.. warning::

    Please be aware that upgrading from 0.7.5 Beta (and prior) to 0.7.6 Beta requires modification of your core GSD configuration files! **You must be running GSD version 0.1.4 on each node before upgrading! GSD MUST BE ON.**

Update GSD Code Base
--------------------

.. code-block:: sh

    [$]~ git reset --hard HEAD
    [$]~ git fetch --tags
    [$]~ git checkout tags/<version>
    [$]~ npm update && npm start

Update MySQL User
-----------------
If you are running this upgrader and are updating from versions prior to ``0.7.5`` you will first need to update your MySQL user to have the correct permissions for the database. Please run the command below to update the user.

.. code-block:: mysql

  GRANT ALTER, DROP ON database.* TO 'user'@'localhost';

Update the PufferPanel Code Base
--------------------------------
Thanks to git updating the code is a very simple process. Simply run the commands below to update the code.

.. warning::

    Running *git reset --hard HEAD* will delete any changes you made to PufferPanel files (except for your config file). In most cases this should not pose a problem, but any customizations will be overwitten. You should back these up to a secure location before running this.

.. code-block:: sh

  [$]~ git reset --hard HEAD
  [$]~ git fetch --tags
  [$]~ git checkout tags/<version>

If you don't know what version is the latest simply run the command below to list them all and find the most recent.

.. code-block:: sh

  [$]~ git tag -l

Handle File Uploads
-------------------

.. note::

    You only need to do this if you are upgrading from 0.7.5 to 0.7.6.

Beginning in ``0.7.6`` file uploads are possible. Because of this you will need to modify Apache or Nginx to allow for file uploads of up to 100MB. PufferPanel sets a hard limit of ``100MB`` per file which is handled through the code for PHP (so you do not need to update your ``php.ini`` file).

Configuring Apache Uploads
^^^^^^^^^^^^^^^^^^^^^^^^^^
You should not need to change anything in your Apache configuration for file uploads to work.

Configuring Nginx Uploads
^^^^^^^^^^^^^^^^^^^^^^^^^
Add the lines below to your ``nginx.conf`` file in the ``http`` block and then restart nginx using ``service nginx restart``.

.. code-block:: text

    http {

    [...]

    client_max_body_size 100m;
    client_body_timeout 120s;

    }

Running the Upgrader
--------------------
Running the upgrader is a very simple step and often requires no extra dependencies to be installed on your server.

First we need to update some new folders to have the correct permissions.

.. code-block:: sh

    chmod -R 0777 src/cache

Then we need to run composer again and check for any upgrades or new software.

.. code-block:: sh

  php composer.phar self-update
  php composer.phar update

After that, you should navigate to your PufferPanel install in your browser, and go to ``http://example.com/install/upgrade/index.php``. After doing that, select the version that you are upgrading from, and click start. The upgrader will update all of the tables necessary, and let you know when it finished.

Finishing
---------
When finished run the command below to remove the install and upgrader.

.. code-block:: sh

  [$]~ cd /var/www/example.com
  [$]~ rm -rf panel/install
