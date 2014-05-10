#!/bin/sh

#############################################################
#                                                           #
#                       PufferPanel                         #
#                                                           #
#               Ubuntu installer by Chomiciak               #
#                                                           #
# Use on the same license as Panel's (and on your own risk!)#
#                                                           #
#############################################################
#                                                           #
#                 This is Master installer                  #
#                                                           #
#############################################################

# Info...
echo "You must configure some software... I will help you a little. Wait 7 seconds..."
# Wait...
read -t 7

# Update!
apt-get update

# Apache2 install
apt-get install apache2

# Php5 install
apt-get install php5

# Mysql install
apt-get install mysql-server

# Apache-mysql
apt-get install libapache2-mod-auth-mysql

# Php-mysql
apt-get install php5-mysql

# SSH2-php
apt-get install libssh2-1-dev libssh2-php

# pdo-php
pecl install pdo

# mcrypt-php
pecl install mcrypt

# Apache restart
/etc/init.d/apache2 restart
