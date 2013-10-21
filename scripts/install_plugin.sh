#!/bin/bash
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