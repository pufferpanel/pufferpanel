#!/bin/bash
#
# PufferPanel Command Sender Script
# Parameters: 
#	$1 = user
#	$2 = command
# Usage:
#	./send_command.sh user "say This is my command!"


su - $1 -s /bin/bash -c "screen -q -S minecraft_$1 -p bukkit_$1 -X stuff \"$2$(echo -ne '\r')\" > /dev/null; exit;"

exit 0