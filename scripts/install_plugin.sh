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
# PufferPanel Plugin Install Script
# Parameters: 
#	$1 = user
#	$2 = dl_link
#	$3 = server_plugin_path
#	$4 = plugin folder name
#	$5 = plugin file name

su - $1 -s /bin/bash -c "mkdir -p $3 && cd $3 && rm -rf $5 && wget $2 && exit"

file=$(echo "$3/$5"|awk -F . '{print $NF}')

#Unzip and remove zip, otherwise ignore it
if [ "$file" == "zip" ];
	then
		su - $1 -s /bin/bash -c "cd $3 && rm -rf $4 && mkdir $4 && unzip -od $4 \"$5\" && mv $4/*.jar ./ && rm -rf \"$5\""
fi

exit