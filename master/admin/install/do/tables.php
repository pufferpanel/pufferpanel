<?php
/*
    PufferPanel - A Minecraft Server Management Panel
    Copyright (c) 2013 Dane Everitt
 
    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.
 
    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.
 
    You should have received a copy of the GNU General Public License
    along with this program.  If not, see http://www.gnu.org/licenses/.
 */
if(file_exists('../install.lock'))
 	exit('Installer is Locked.');
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>PufferPanel - Install</title>
	
	<!-- Stylesheets -->
	<link href='http://fonts.googleapis.com/css?family=Droid+Sans:400,700' rel='stylesheet'>
	<link rel="stylesheet" href="../../../assets/css/style.css">
	
	<!-- Optimize for mobile devices -->
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
	
	<!-- jQuery & JS files -->
	<script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
</head>
<body>
	<div id="top-bar">
		<div class="page-full-width cf">
            &nbsp;
		</div>	
	</div>
	<div id="header-with-tabs">
		<div class="page-full-width cf">
		</div>
	</div>
	<div id="content">
		<div class="page-full-width cf">
            <div class="content-module">
				<div class="content-module-main">
				    <h1>Installing Tables</h1>
                    <p>Your database details were correct, the script is now installing the necessary tables.</p>
                    <div class="stripe-separator"></div>
                    <pre>
                    <?php
                    
                    	if(!file_exists('../../../core/framework/master_configuration.php'))
                    		echo '<div class="error-box round">The configuration file was not found.</div>';	
                    	else {
                    	
	                        include('../../../core/framework/framework.database.connect.php');
	                        $mysql = dbConn::getConnection();
	
	                        /*
	                         * CREATE TABLE `account_change`
	                         */
	                        $mysql->exec("CREATE TABLE `account_change` (
	                          `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	                          `user_id` int(11) DEFAULT NULL,
	                          `type` varchar(50) NOT NULL DEFAULT '',
	                          `content` text NOT NULL,
	                          `key` varchar(50) NOT NULL DEFAULT '',
	                          `time` int(15) NOT NULL,
	                          `verified` int(1) NOT NULL DEFAULT '0',
	                          PRIMARY KEY (`id`)
	                        ) ENGINE=InnoDB DEFAULT CHARSET=latin1");
	                        echo "\nTable `account_change` created.\n";
	
	                        /*
	                         * CREATE TABLE `acp_announcements`
	                         */
	                        $mysql->exec("CREATE TABLE `acp_announcements` (
	                          `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	                          `text` text NOT NULL,
	                          `enabled` int(1) NOT NULL DEFAULT '1',
	                          `priority` int(4) NOT NULL DEFAULT '0',
	                          PRIMARY KEY (`id`)
	                        ) ENGINE=InnoDB DEFAULT CHARSET=latin1");
	                        echo "Table `acp_announcements` created.\n";
	
	                        /*
	                         * CREATE TABLE `acp_email_templates`
	                         */
	                        $mysql->exec("CREATE TABLE `acp_email_templates` (
	                          `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	                          `tpl_name` char(30) DEFAULT NULL,
	                          `tpl_content` text,
	                          PRIMARY KEY (`id`)
	                        ) ENGINE=InnoDB DEFAULT CHARSET=latin1");
	                        echo "Table `acp_email_templates` created.\n";
	
	                            /* 
	                             * Add Email Templates
	                             */
	                             $mysql->exec("INSERT INTO `acp_email_templates` (`id`, `tpl_name`, `tpl_content`) VALUES
	                                (1,'login_failed','<html>\n	<head>\n		<title><%HOST_NAME%> - Account Login Failure Notification</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> - Account Login Failure Notification</h1></center>\n		<p>You are recieving this email as a part of our continuing efforts to improve our server security. <strong>An unsucessful login was made with your account on <%HOST_NAME%>.</strong></p>\n		<p><strong>IP Address:</strong> <%IP_ADDRESS%><br />\n			<strong>Hostname:</strong> <%GETHOSTBY_IP_ADDRESS%><br />\n			<strong>Time:</strong> <%DATE%></p>\n		<p>At this time your account is still safe and sound in our system. This email is simply to let you know that someone tried to login to your account and failed. You can change your notification preferences by <a href=\"<%MASTER_URL%>accounts.php\">clicking here</a>.</p>\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>'),
	                                (2,'login_success','<html>\n	<head>\n		<title><%HOST_NAME%> - Account Login Notification</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> - Account Login Notification</h1></center>\n		<p>You are recieving this email as a part of our continuing efforts to improve our server security. You are recieving this email as a sucessful login was made with your account on <%HOST_NAME%>.</strong></p>\n		<p><strong>IP Address:</strong> <%IP_ADDRESS%><br />\n			<strong>Hostname:</strong> <%GETHOSTBY_IP_ADDRESS%><br />\n			<strong>Time:</strong> <%DATE%></p>\n		<p>This email is intended to keep you aware of any possible malicious account activity. You can change your notification preferences by <a href=\"<%MASTER_URL%>accounts.php\">clicking here</a>.</p>\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>'),
	                                (3,'email_changed','<html>\n	<head>\n		<title><%HOST_NAME%> Email Change Notification</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> Email Change Notification</h1></center>\n		<p>Hello there! You are receiving this email because you requested to change your account email on <%HOST_NAME%>.</p>\n		<p>Please click the link below to confirm that you wish to use this email address on <%HOST_NAME%>. If you did not make this request, or do not wish to continue simply ignore this email and nothing will happen. <strong>This link will expire in 4 hours.</strong></p>\n		<p><a href=\"<%MASTER_URL%>account_actions.php?conf=email&key=<%EMAIL_KEY%>\"><%MASTER_URL%>account.php?conf=email&key=<%EMAIL_KEY%></a></p>\n		<p>This change was requested from <%IP_ADDRESS%> (<%GETHOSTBY_IP_ADDRESS%>) on <%DATE%>. Please do not hesitate to contact us if you belive something is wrong.\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>'),
	                                (4,'password_changed','<html>\n	<head>\n		<title><%HOST_NAME%> Password Change Notification</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> Password Change Notification</h1></center>\n		<p>Hello there! You are receiving this email because you recently changed your password on <%HOST_NAME%>.</p>\n		<p>This change was requested from <%IP_ADDRESS%> (<%GETHOTBY_IP_ADDRESS%>) on <%DATE%>. If you did not request this change then you should immediately check your computer for anything suspicious, and then login and change your password. You should also immediately contact support so that we can help to protect your account.\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>'),
	                                (5,'password_reset','<html>\n	<head>\n		<title><%HOST_NAME%> Lost Password Recovery</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> Lost Password Recovery</h1></center>\n		<p>Hello there! You are receiving this email because you requested a new password for your <%HOST_NAME%> account.</p>\n		<p>Please click the link below to confirm that you wish to change your password. If you did not make this request, or do not wish to continue simply ignore this email and nothing will happen. <strong>This link will expire in 4 hours.</strong></p>\n		<p><a href=\"<%MASTER_URL%>password.php?key=<%PKEY%>\"><%MASTER_URL%>password.php?key=<%PKEY%></a></p>\n		<p>This change was requested from <%IP_ADDRESS%> (<%GETHOSTBY_IP_ADDRESS%>) on <%DATE%>. Please do not hesitate to contact us if you belive something is wrong.\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>'),
	                                (7,'new_password','<html>\n	<head>\n		<title><%HOST_NAME%> - New Password</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> - New Password</h1></center>\n		<p>Hello there! You are receiving this email because you requested a new password for your <%HOST_NAME%> account.</p>\n		<p><strong>Login:</strong> <a href=\"<%MASTER_URL%>\"><%MASTER_URL%></a><br />\n			<strong>Email:</strong> <%EMAIL%><br />\n			<strong>Password:</strong> <%NEW_PASS%></p>\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>'),
	                                (8,'admin_newaccount','<html>\n	<head>\n		<title><%HOST_NAME%> - Account Created</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> - Account Created</h1></center>\n		<p>Hello there! This email is to inform you that an account has been created for you on <%HOST_NAME%>.</p>\n		<p><strong>Login:</strong> <a href=\"<%MASTER_URL%>\"><%MASTER_URL%></a><br />\n			<strong>Email:</strong> <%EMAIL%><br />\n			<strong>Password:</strong> <%PASS%></p>\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>'),
	                                (9,'admin_new_sftppass','<html>\n	<head>\n		<title><%HOST_NAME%> - SFTP Password Changed</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> - SFTP Password Changed </h1></center>\n		<p>Hello there! This email is to inform you that the SFTP password for <%SERVER%> has been changed by an administrator.</p>\n		<p><strong>New Password:</strong> <%PASS%><br />\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>'),
	                                (10,'admin_new_server','<html>\n	<head>\n		<title><%HOST_NAME%> - New Server Added</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> - New Server Added </h1></center>\n		<p>Hello there! This email is to inform you that a new server (<%NAME%>) has been created for you.</p>\n		<p><strong>Connect:</strong> <%CONNECT%><br />\n		<p><strong>SFTP Username:</strong> <%USER%><br />\n		<p><strong>SFTP Password:</strong> <%PASS%><br />\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>')");
	                            echo "Table `acp_email_templates` was populated with template data.\n";
	
	                        /*
	                         * CREATE TABLE `acp_settings`
	                         */
	                        $mysql->exec("CREATE TABLE `acp_settings` (
	                          `setting_ref` char(25) NOT NULL DEFAULT '',
	                          `setting_val` tinytext
	                        ) ENGINE=InnoDB DEFAULT CHARSET=latin1");
	                        echo "Table `acp_settings` created.\n";
	
	                        /*
	                         * CREATE TABLE `actions_log`
	                         */
	                        $mysql->exec("CREATE TABLE `actions_log` (
	                          `id` int(1) unsigned NOT NULL AUTO_INCREMENT,
	                          `priority` int(1) NOT NULL,
	                          `viewable` int(1) NOT NULL DEFAULT '0',
	                          `user` int(1) DEFAULT NULL,
	                          `time` int(1) NOT NULL,
	                          `ip` char(100) NOT NULL DEFAULT '',
	                          `url` tinytext NOT NULL,
	                          `action` char(100) NOT NULL DEFAULT '',
	                          `desc` text NOT NULL,
	                          PRIMARY KEY (`id`)
	                        ) ENGINE=InnoDB DEFAULT CHARSET=latin1");
	                        echo "Table `actions_log` created.\n";
	
	                        /*
	                         * CREATE TABLE `backup_datastore`
	                         */
	                        $mysql->exec("CREATE TABLE `backup_datastore` (
	                          `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	                          `server` text NOT NULL,
	                          `backup_pattern` longtext,
	                          PRIMARY KEY (`id`)
	                        ) ENGINE=InnoDB DEFAULT CHARSET=latin1");
	                        echo "Table `backup_datastore` created.\n";
	
	                        /*
	                         * CREATE TABLE `backups`
	                         */
	                        $mysql->exec("CREATE TABLE `backups` (
	                          `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	                          `server` varchar(100) NOT NULL DEFAULT '',
	                          `backup_token` varchar(50) NOT NULL DEFAULT '',
	                          `file_name` text NOT NULL,
	                          `email_done` int(1) NOT NULL DEFAULT '0',
	                          `user_email` text,
	                          `timestart` int(15) NOT NULL,
	                          `timeend` int(15) DEFAULT NULL,
	                          `complete` int(1) NOT NULL DEFAULT '0',
	                          `md5` text,
	                          `sha1` text,
	                          PRIMARY KEY (`id`)
	                        ) ENGINE=InnoDB DEFAULT CHARSET=latin1");
	                        echo "Table `backups` created.\n";
	
	                        /*
	                         * CREATE TABLE `nodes`
	                         */
	                        $mysql->exec("CREATE TABLE `nodes` (
	                          `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	                          `node` char(15) NOT NULL DEFAULT '',
	                          `node_link` tinytext NOT NULL,
	                          `node_ip` tinytext NOT NULL,
	                          `sftp_ip` tinytext NOT NULL,
	                          `server_dir` tinytext NOT NULL,
	                          `backup_dir` tinytext NOT NULL,
	                          `username` tinytext NOT NULL,
	                          `encryption_iv` tinytext NOT NULL,
	                          `password` tinytext NOT NULL,
	                          `ips` text NOT NULL,
	                          `ports` text NOT NULL,
	                          `dsdir` tinytext NOT NULL,
	                          PRIMARY KEY (`id`)
	                        ) ENGINE=InnoDB DEFAULT CHARSET=latin1");
	                        echo "Table `nodes` created.\n";
	                        
	                        /*
	                         * CREATE TABLE `jars`
	                         */
	                         $mysql->exec("CREATE TABLE `jars` (
	                          `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	                          `dir` tinytext NOT NULL,
	                          `name` tinytext NOT NULL DEFAULT '',
	                          `desc` text NOT NULL,
	                          `default` int(1) unsigned NOT NULL,
	                          `cat` int(11) unsigned NOT NULL,
	                          PRIMARY KEY (`id`)
	                        ) ENGINE=InnoDB DEFAULT CHARSET=latin1");
	                        echo "Table `jars` created.\n";
				//cat - category, for user permissions management, i.e. One user can user that file,
				//another cant.
	                        /*
	                         * CREATE TABLE `servers`
	                         */
	                        $mysql->exec("CREATE TABLE `servers` (
	                          `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	                          `whmcs_id` int(11) DEFAULT '0',
	                          `hash` char(42) NOT NULL DEFAULT '',
	                          `encryption_iv` tinytext NOT NULL,
	                          `node` int(11) NOT NULL,
	                          `name` varchar(200) NOT NULL DEFAULT '',
	                          `active` int(1) DEFAULT '1',
	                          `owner_id` int(11) NOT NULL,
	                          `max_ram` int(11) NOT NULL,
	                          `disk_space` int(11) NOT NULL,
	                          `date_added` int(15) NOT NULL,
	                          `server_ip` varchar(50) NOT NULL DEFAULT '',
	                          `server_port` int(11) NOT NULL,
	                          `ftp_host` tinytext NOT NULL,
	                          `ftp_user` tinytext NOT NULL,
	                          `ftp_pass` tinytext NOT NULL,
	                          `backup_file_limit` int(20) NOT NULL,
	                          `backup_disk_limit` int(20) NOT NULL,
	                          PRIMARY KEY (`id`)
	                        ) ENGINE=InnoDB DEFAULT CHARSET=latin1");
	                        echo "Table `servers` created.\n";
	
	                        /*
	                         * CREATE TABLE `users`
	                         */
	                        $mysql->exec("CREATE TABLE `users` (
	                          `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	                          `whmcs_id` int(11) DEFAULT NULL,
	                          `username` varchar(50) NOT NULL DEFAULT '',
	                          `email` tinytext NOT NULL,
	                          `password` tinytext NOT NULL,
	                          `register_time` int(15) NOT NULL,
	                          `position` varchar(50) DEFAULT '' COMMENT 'owner,admin,staff',
	                          `session_id` varchar(12) DEFAULT '',
	                          `session_ip` varchar(50) DEFAULT '',
	                          `session_expires` int(15) DEFAULT NULL,
	                          `root_admin` int(1) NOT NULL DEFAULT '0',
	                          `notify_login_s` int(1) DEFAULT '1',
	                          `notify_login_f` int(1) DEFAULT '1',
	                          PRIMARY KEY (`id`)
	                        ) ENGINE=InnoDB DEFAULT CHARSET=latin1");
	                        echo "Table `users` created.\n";
	
	                        /*
	                         * CREATE TABLE `whmcs_suspend_data`
	                         */
	                        $mysql->exec("CREATE TABLE `whmcs_suspend_data` (
	                          `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	                          `server_id` int(11) NOT NULL,
	                          `whmcs_server_id` int(11) NOT NULL,
	                          `old_password` varchar(40) NOT NULL DEFAULT '',
	                          `unsuspended` int(11) NOT NULL DEFAULT '0' COMMENT '0 = False, 1 = True',
	                          PRIMARY KEY (`id`)
	                        ) ENGINE=InnoDB DEFAULT CHARSET=latin1");
	                        echo "Table `whmcs_suspend_data` created.\n";
	                        
						}

                    ?>
                    </pre>
               	    <a href="settings.php">Next Step &rarr;</a>
				</div>
            </div>
		</div>
	</div>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.4.2 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>
