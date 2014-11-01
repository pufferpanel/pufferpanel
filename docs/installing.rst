Installing PufferPanel
======================
Please be aware that PufferPanel has not been tested with Windows based environments. The documentation on this wiki is written for unix systems, and has been tested on Ubuntu 14.04 LTS.

Dependencies
------------
* Git
* PHP (version 5.5 or higher with the following modules: ``php-pear``, ``php-cli``, ``php-curl``, ``php-mcrypt``)
* Apache or Nginx
* MySQL
  * Version 5.5 or higher is required. MariaDB should work as well.

Downloading
-----------
You will need to be comfortable with SSH in order to install this panel. PufferPanel can be automatically downloaded and placed in the correct folder by executing the commands below.

.. code-block:: sh
  [$]~ cd /var/www/yoursite.com
  [$]~ git clone https://github.com/DaneEveritt/PufferPanel.git
  [$]~ cd PufferPanel && git checkout tags/<version>
  [$]~ cp -R ./ ../
  [$]~ cd ../ && rm -rf PufferPanel/

Replace ``<version>`` with whichever version of the panel you wish to install. We highly recommend using the latest release when possible. As of the last update, the latest release is ``0.7.4.1-beta``. You must checkout a branch in order to use the panel, running the ``master`` repo is hghly unstable and should not be used.

Configure and Prepare the Panel
-------------------------------
Edit your Apache or Nginx settings to make ``panel`` the home folder. If you do not want to have ``panel/`` as the home folder, please consider either using a subdomain for the panel, or make completely sure that none of the other included files or folders can be accessed through a browser.

After setting up Apache/Nginx you will need to run composer which will install the dependencies for the panel. This is a very simple step, but is often the reason most installs fail. To download composer run the command below in the PufferPanel directory (this directory should have ``src``, ``app``, and ``panel`` folders in it).

.. code-block:: sh
  [$]~ curl -sS https://getcomposer.org/installer | php
  [$]~ php composer.phar install

After running composer we need to setup the other folders for the installer so that we don't have any errors occur. Run the commands below to do this.

.. code-block:: sh
  [$]~ chmod -R 0777 panel/admin/install
  [$]~ chmod -R 0777 src/framework
  [$]~ chmod 0666 src/framework/configuration.php.dist

Configuring MySQL
-----------------
You will need to add a non-root MySQL user for the panel to operate. Please run the command below in your MySQL terminal, replacing ``database`` and ``user`` with the name of the database you are using for puffer panel, and the name of the user who will have access to it respectively. You may need to restart MySQL after running the following commands.

.. code-block:: sh
  [$]~ mysql -u root -p
  mysql> CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';
  mysql> GRANT CREATE, INSERT, SELECT, UPDATE, DELETE, DROP ON database.* TO 'user'@'localhost';
  mysql> FLUSH PRIVILEGES;

Running the Installer
---------------------
Point your browser to the panel, and go to ``/admin/install`` following the instructions from there. This will set up the MySQL database, general settings, hashing information, and the root administrator account.

Cleaning Up
-----------
Delete the ``/panel/admin/install`` directory and fix the file permissions in other directories.

.. code-block:: sh
  [$]~ rm -rf panel/admin/install
  [$]~ chmod -R 0755 src/framework
  [$]~ chmod 0444 src/framework/configuration.php

After doing that, we need to change some permissions for our MySQL user for security purposes. You should restart MySQL after running the following command.

.. code-block:: sh
  [$]~ mysql -u root -p
  mysql> REVOKE DROP ON database.* FROM 'user'@'localhost';
  mysql> FLUSH PRIVILEGES;

Congratulations! You should have PufferPanel running smoothly at this point. From here, you should move on to `setting up your first node <installing_nodes.md>`_.