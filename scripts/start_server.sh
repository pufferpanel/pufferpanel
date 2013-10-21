#!/bin/bash
#
# PufferPanel Server Start Script
# Parameters: 
#	$1 = server_path
#	$2 = max_ram
#	$3 = run_as_user

invoke="java -Xmx$2M -Xincgc -XX:+CMSIncrementalPacing -XX:ParallelGCThreads=2 -XX:+AggressiveOpts -jar server.jar --nojline"
su - $3 -s /bin/bash -c "cd $1; screen -S minecraft_$3 -t bukkit_$3 -dm $invoke; exit;"

exit 0
