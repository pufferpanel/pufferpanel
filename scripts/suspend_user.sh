#!/bin/bash
#
# PufferPanel User Suspension Script
# Parameters: 
#	$1 = server_name
#	$2 = new_password

#Set Password
echo -e "$2\n$2" | passwd $1

#End Setup Script
exit 0