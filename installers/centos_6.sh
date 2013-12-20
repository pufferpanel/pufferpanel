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
echo "!!NOTICE!! This script will NOT install MySQL or NGINX."
echo "USE AT YOUR OWN RISK!"

rpm -Uvh http://download.fedoraproject.org/pub/epel/6/i386/epel-release-6-8.noarch.rpm
rpm -Uvh http://rpms.famillecollet.com/enterprise/remi-release-6.rpm

yum update -y

echo "Installing PHP Dependences"
yum install -y php-fpm php-pdo php-devel pecl

pecl install ssh2
pecl install mcrypt

echo "Installing Other Dependences"
yum install -y screen gzip quota

echo "Done."