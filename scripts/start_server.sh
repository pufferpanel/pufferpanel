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
# PufferPanel Server Start Script
# Parameters: 
#	$1 = server_path
#	$2 = max_ram
#	$3 = run_as_user

invoke="java -Xmx$2M -Xincgc -XX:+CMSIncrementalPacing -XX:ParallelGCThreads=2 -XX:+AggressiveOpts -jar server.jar --nojline"
su - $3 -s /bin/bash -c "cd $1; screen -S minecraft_$3 -t bukkit_$3 -dm $invoke; exit;"

exit 0
