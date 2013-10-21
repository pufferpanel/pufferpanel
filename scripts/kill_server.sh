#!/bin/bash
#
# PufferPanel Server PID Kill Script
# Parameters: 
#	$1 = user

#su - $1 -s /bin/bash -c "screen -S minecraft_$1 -X quit"

#Get Java Process
bukkitPID=`ps -U $1 | grep java`

#Kill It, with Fire!
kill -9 ${bukkitPID:0:5}

exit 0