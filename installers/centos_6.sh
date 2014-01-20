#!/bin/bash
#
#   PufferPanel - A Minecraft Server Management Panel
#   Copyright (c) 2013 Dane Everitt
#
#   This program is free software: you can redistribute it and/or modify
#   it under the terms of the GNU General Public License as published by
#   the Free Software Foundation, either version 3 of the License, or
#   (at your option) any later version.
#
#   This program is distributed in the hope that it will be useful,
#   but WITHOUT ANY WARRANTY; without even the implied warranty of
#   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#   GNU General Public License for more details.
#
#   You should have received a copy of the GNU General Public License
#   along with this program.  If not, see http://www.gnu.org/licenses/.
#

echo "Beginning CentOS Automatic Installer for PufferPanel"
echo "USE AT YOUR OWN RISK!"

read -p "Do you want to install the EPEL and REMI repos? This is recommended to get the most up-to-date versions of software. [Y/n] " -n 1 -r
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    rpm -Uvh http://download.fedoraproject.org/pub/epel/6/i386/epel-release-6-8.noarch.rpm
    rpm -Uvh http://rpms.famillecollet.com/enterprise/remi-release-6.rpm
fi

yum update -y

read -p "Do you want to install Apache and MySQL? These are required, unless you are using a different web server or database solution. [Y/n] " -n 1 -r
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
echo "Installing Apache and MySQL"
yum install -y mysql-server httpd
fi

echo "Installing PHP Dependences (assuming apache)"
yum install -y php php-mysql php-pdo php-devel php-ssh2 php-mcrypt pear

echo "Installing Other Dependences"
yum install -y screen quota

echo "Done."
