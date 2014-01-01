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
# PufferPanel User Creation Script
# Parameters: 
#	$1 = server_name
#	$2 = sftp_password
#	$3 = soft_limit
#	$4 = hard_limit
# $5 = default server's folder

#Add User
useradd -d /srv/servers/$1 -s /usr/bin/rssh -G rsshusers $1

#Set Password
echo -e "$2\n$2" | passwd $1

#Set Folder Permissions
mkdir /srv/servers/$1/server
chown root.root /srv/servers/$1
chmod 755 /srv/servers/$1
chown $1.rsshusers /srv/servers/$1/server

#Copy files of the default server to user's folder
cp -r $5/* /srv/servers/$1/server

#Set Disk Limits
setquota -u $1 $3 $4 0 0 -a

#End Setup Script
exit 0
