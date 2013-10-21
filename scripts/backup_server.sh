#!/bin/bash
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
	curl "http://${data[1]}.panel.pufferhost.com/pages/server/backup.php?do=backup_done&server=${data[2]}&token=${data[3]}" > /dev/null 2>&1
	
	#Remove tmp
	rm -rf /second/backups/$1/tmp

exit 0