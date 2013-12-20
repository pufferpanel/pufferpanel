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
#	$1 = run_as_user
#	$2 = array
#		-> 0: backup_name
#		-> 1: node
#		-> 2: server_hash
#		-> 3: backup_token
#	$3 = backup_files
#	$4 = skip_files

#Remove Old backup Folder
#rm -rf /second/backups/$1/tmp

#Create Tempoary Backup Folder
mkdir -p /second/backups/$1/tmp

#Set Up Arrays

data=( $2 )
backup=( $3 )
skip=( $4 )

	#Backup Specificed Files
	cd /srv/servers/$1/server
	
		#Check for *
		if [ "$3" = "*" ]
			then

				rsync -azR * /second/backups/$1/tmp
				
			else
			
				backup=( $3 )
				for i in "${backup[@]}"
				do
					rsync -azR $i /second/backups/$1/tmp
				done
			
			fi

	
	#Remove Specific Files
	cd /second/backups/$1/tmp
	for i in "${skip[@]}"
	do
		rm -rf /second/backups/$1/tmp/$i
	done
	
	#Generate Compressed File
	cd /second/backups/$1/tmp 
	tar -zcf /second/backups/$1/${data[0]}.tar.gz *
	
	#Tell Server backup is Finished
	curl "http://localhost/backup_shellreturn.php?do=backup_done&server=${data[2]}&token=${data[3]}" > /dev/null 2>&1
	
	#Remove tmp
	rm -rf /second/backups/$1/tmp

exit 0