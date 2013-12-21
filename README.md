#About PufferPanel

Minecraft Web Management Panel for Hosts.

irc.esper.net #pufferpanel

This software is currently unstable and should not be used on a live environment.

#Installing PufferPanel

###Dependencies:

* Nginx
* MySQL (5.5+ or equivalent)
* PHP (5.3+) with the following modules:
    * ssh2
    * openssl
    * PDO
    * pdo_mysql
    * hash
    * curl
    * mcrypt
* screen
* gzip
* rssh
* setquota


####Automatic Install
PufferPanel comes with an automatic installer located in `/master/admin/install/`. After uploading all of the files, go to this directory and follow the instructions. This will get PufferPanel setup with a new account, and all of the basic settings and MySQL information.

After running the automatic installer you will need to manually add a new server. Please ensure that you have a node setup correctly before doing this. You will need to setup RSSH, as well as create a user who has rights to all of the server files and backup directory.

You will also need to upload the `scripts` folder to `/srv/scripts` and chmod all of the files using `chmod +x *`, and create a server directory in `/srv/servers`. Your backup directory can be wherever you want on the server, as right now it does not do off-site backups (in development). We suggest `/srv/backups` or a second hard-drive if possible.

You should then open `/srv/scripts/backup_server.sh` and edit line 75 changing `http://localhost/` to the URL for your node that the script is located on.

#####Setting up Node
You will need to copy the master_configuration.php file from the master directory into the same location in the node folder, renaming it node_configuration.php.

###Using RSSH
Make the group for RSSH Users
`[root@vpn ~]# groupadd rsshusers`

####Edit SSHD Config:
`/etc/ssh/sshd_config`

	Subsystem	sftp    internal-sftp
	Match Group rsshusers
	ChrootDirectory /srv/servers/%u
	X11Forwarding no
	AllowTcpForwarding no
	ForceCommand internal-sftp

####Grant SFTP Access to RSSH Users

The default action for rssh to lock down everything. To grant access sftp open the RSSH file:
`[root@vpn ~]# /etc/rssh.conf`

Append or uncomment following line:
`allowsftp`

####Setup Permissions

	[root@vpn ~]# service sshd restart
