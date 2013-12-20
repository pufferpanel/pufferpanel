PufferPanel
===========

Minecraft Web Management Panel for Hosts.

irc.esper.net #pufferpanel

This software is currently unstable and should not be used on a live environment.

#Installing PufferPanel

###Dependencies:

* nginx or Apache
* MySQL (5.5+) or MariaDB (10.0+)
* PHP (5.3+)
    * SSH2
    * OpenSSL
    * PDO
    * mcrypt
* screen
* gzip
* rssh    

###MySQL Setup
Create a database named `pufferpanel` and edit the connection details for the database in the following files:

`/master/core/framework/php/master_configuration.php.dist`

`/node/core/framework/php/node_configuration.php.dist`

Make sure to remove the `.dist` from the ends of the file name.

After that, upload the `db_structure.sql` file into `pufferpanel`.

### Making a User Account
To create your account you will need to manually edit the `pufferpanel.users` table. To do this, run the following command on the MySQL table:

	INSERT INTO `users` (`id`, `whmcs_id`, `username`, `email`, `password`, `register_time`, `position`, `session_id`, `session_ip`, `session_expires`, `root_admin`, `notify_login_s`, `notify_login_f`) VALUES(NULL, NULL, 'YOUR_USERNAME', 'YOUR_EMAIL', 'YOUR_PASSWORD', 0, 'owner', NULL, NULL, NULL, 1, 0, 0);
    
Replace `YOUR_USERNAME` with the username you want to use for your RSSH user (discussed below). Replace `YOUR_EMAIL` with the email you want to use to login to PufferPanel. In order to get an encrypted password for `YOUR_PASSWORD` follow the directions below.

####Creating an Encrypted Password
Open the file `/master/encrypt.php` and edit the values it says to edit. You will need to create an AES encryption key. Use the [GRC](https://www.grc.com/passwords.htm) site and copy the ASCII line into a file somewhere **outside of the web accessable part of your server**. Name this file whatever you want to name it. Edit line 4 of the file to point to the file. Edit line 13 to have a new slat for your passwords. You will  not be able to change this later.

After this, edit `/master/core/framework/php/framework.auth.php` and edit line 18 to match the salt you selected. Ignore the one already in there. Afterwards, edit the same file in the node folder (`/node/core/framework/php/framework.auth.php`) doing the same thing.

To generate your password go to `http://pufferpanel/master/encrypt.php` and add the following GET variables to the end: `openssl` and `ripemd`. The `openssl` variable should be the password you are going to use for your user that can access via SFTP (discussed below). The `ripemd` variable should eb the password you want to use for your PufferPanel account. When you do this, it will output your encrypted passwords. Copy the `ripemd` password into the database for `YOUR_PASSWORD`. IF you re-run the file the `openssl` password will change each time.

Example Request: `http://pufferpanel/master/encrypt.php?openssl=sftppassword&ripemd=mypassword`

Ignore the OpenSSL password for now, it can be generated in the Admin CP when you make a server.

###Using RSSH
You will need to do this manually for now until an automatic script is made to interface with the Admin CP server creation.

Make the group for RSSH Users
`[root@vpn ~]# groupadd rsshusers`

Add a user for a server. (No SSH Access, SFTP Only)
`[root@vpn ~]# useradd -d /srv/servers/dane -s /usr/bin/rssh -G rsshusers dane`

Change the user's password
`[root@vpn ~]# passwd dane`

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
	[root@vpn ~]# mkdir /srv/servers/dane/server
	[root@vpn ~]# chown root.root /srv/servers/dane
	[root@vpn ~]# chmod 755 /srv/servers/dane
	[root@vpn ~]# chown dane.rsshusers /srv/servers/dane/server
