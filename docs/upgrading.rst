Upgrading PufferPanel
=====================
Starting in version ``0.7.5 Beta`` we now support running an upgrader rather than having to re-install the entire panel. Please be aware that this upgrader is still very much beta software, and may still cause issues. Use at your own risk.

.. warning::

    Please be aware that upgrading from 0.7.5 Beta (and prior) to 0.7.6 Beta requires modification of your core GSD configuration files! **You must be running GSD version ``0.1.4`` on each node before upgrading! GSD MUST BE ON.**

Update MySQL User
-----------------
If you are running this upgrader and are updating from versions prior to ``0.7.5`` you will first need to update your MySQL user to have the correct permissions for the database. Please run the command below to update the user.

.. code-block:: mysql

  GRANT ALTER, DROP ON database.* TO 'user'@'localhost';

Update the Code Base
--------------------
Thanks to git updating the code is a very simple process. Simply run the commands below to update the code.

.. code-block:: sh

  [$]~ git fetch --tags
  [$]~ git checkout tag/<version>

If you don't know what version is the latest simply run the command below to list them all and find the most recent.

.. code-block:: sh

  [$]~ git tag -l

Running the Upgrader
--------------------
Running the upgrader is a very simple step and often requires no extra dependencies to be installed on your server.

First, we need to run composer again and check for any upgrades or new software.

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
