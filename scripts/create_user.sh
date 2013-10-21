#!/bin/bash
#
# PufferPanel User Creation Script
# Parameters: 
#	$1 = server_name
#	$2 = sftp_password
#	$3 = soft_limit
#	$4 = hard_limit

#Create Home Directory
mkdir -p /srv/servers/$1/server

#Add User
useradd -m -d /srv/servers/$1 -s /usr/bin/rssh $1

#Set Password
echo -e "$2\n$2" | passwd $1

#Set Folder Permissions
#chown -R $1 /srv/servers/$1/server

#Set Disk Limits
setquota -u $1 $3 $4 0 0 -a

#End Setup Script
exit 0