# ************************************************************
# Sequel Pro SQL dump
# Version 4096
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.5.29)
# Database: pufferpanel
# Generation Time: 2013-10-23 04:05:05 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table account_change
# ------------------------------------------------------------

DROP TABLE IF EXISTS `account_change`;

CREATE TABLE `account_change` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `type` varchar(50) NOT NULL DEFAULT '',
  `content` text NOT NULL,
  `key` varchar(50) NOT NULL DEFAULT '',
  `time` int(15) NOT NULL,
  `verified` int(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table acp_announcements
# ------------------------------------------------------------

DROP TABLE IF EXISTS `acp_announcements`;

CREATE TABLE `acp_announcements` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `text` text NOT NULL,
  `enabled` int(1) NOT NULL DEFAULT '1',
  `priority` int(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table acp_email_templates
# ------------------------------------------------------------

DROP TABLE IF EXISTS `acp_email_templates`;

CREATE TABLE `acp_email_templates` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `tpl_name` char(30) DEFAULT NULL,
  `tpl_content` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `acp_email_templates` WRITE;
/*!40000 ALTER TABLE `acp_email_templates` DISABLE KEYS */;

INSERT INTO `acp_email_templates` (`id`, `tpl_name`, `tpl_content`)
VALUES
	(1,'login_failed','<html>\n	<head>\n		<title><%HOST_NAME%> - Account Login Failure Notification</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> - Account Login Failure Notification</h1></center>\n		<p>You are recieving this email as a part of our continuing efforts to improve our server security. <strong>An unsucessful login was made with your account on <%HOST_NAME%>.</strong></p>\n		<p><strong>IP Address:</strong> <%IP_ADDRESS%><br />\n			<strong>Hostname:</strong> <%GETHOSTBY_IP_ADDRESS%><br />\n			<strong>Time:</strong> <%DATE%></p>\n		<p>At this time your account is still safe and sound in our system. This email is simply to let you know that someone tried to login to your account and failed. You can change your notification preferences by <a href=\"<%MASTER_URL%>accounts.php\">clicking here</a>.</p>\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>'),
	(2,'login_success','<html>\n	<head>\n		<title><%HOST_NAME%> - Account Login Notification</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> - Account Login Notification</h1></center>\n		<p>You are recieving this email as a part of our continuing efforts to improve our server security. You are recieving this email as a sucessful login was made with your account on <%HOST_NAME%>.</strong></p>\n		<p><strong>IP Address:</strong> <%IP_ADDRESS%><br />\n			<strong>Hostname:</strong> <%GETHOSTBY_IP_ADDRESS%><br />\n			<strong>Time:</strong> <%DATE%></p>\n		<p>This email is intended to keep you aware of any possible malicious account activity. You can change your notification preferences by <a href=\"<%MASTER_URL%>accounts.php\">clicking here</a>.</p>\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>'),
	(3,'email_changed','<html>\n	<head>\n		<title><%HOST_NAME%> Email Change Notification</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> Email Change Notification</h1></center>\n		<p>Hello there! You are receiving this email because you requested to change your account email on <%HOST_NAME%>.</p>\n		<p>Please click the link below to confirm that you wish to use this email address on <%HOST_NAME%>. If you did not make this request, or do not wish to continue simply ignore this email and nothing will happen. <strong>This link will expire in 4 hours.</strong></p>\n		<p><a href=\"<%MASTER_URL%>account_actions.php?conf=email&key=<%EMAIL_KEY%>\"><%MASTER_URL%>account.php?conf=email&key=<%EMAIL_KEY%></a></p>\n		<p>This change was requested from <%IP_ADDRESS%> (<%GETHOSTBY_IP_ADDRESS%>) on <%DATE%>. Please do not hesitate to contact us if you belive something is wrong.\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>'),
	(4,'password_changed','<html>\n	<head>\n		<title><%HOST_NAME%> Password Change Notification</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> Password Change Notification</h1></center>\n		<p>Hello there! You are receiving this email because you recently changed your password on <%HOST_NAME%>.</p>\n		<p>This change was requested from <%IP_ADDRESS%> (<%GETHOTBY_IP_ADDRESS%>) on <%DATE%>. If you did not request this change then you should immediately check your computer for anything suspicious, and then login and change your password. You should also immediately contact support so that we can help to protect your account.\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>'),
	(5,'password_reset','<html>\n	<head>\n		<title><%HOST_NAME%> Lost Password Recovery</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> Lost Password Recovery</h1></center>\n		<p>Hello there! You are receiving this email because you requested a new password for your <%HOST_NAME%> account.</p>\n		<p>Please click the link below to confirm that you wish to change your password. If you did not make this request, or do not wish to continue simply ignore this email and nothing will happen. <strong>This link will expire in 4 hours.</strong></p>\n		<p><a href=\"<%MASTER_URL%>password.php?key=<%PKEY%>\"><%MASTER_URL%>password.php?key=<%PKEY%></a></p>\n		<p>This change was requested from <%IP_ADDRESS%> (<%GETHOSTBY_IP_ADDRESS%>) on <%DATE%>. Please do not hesitate to contact us if you belive something is wrong.\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>'),
	(7,'new_password','<html>\n	<head>\n		<title><%HOST_NAME%> - New Password</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> - New Password</h1></center>\n		<p>Hello there! You are receiving this email because you requested a new password for your <%HOST_NAME%> account.</p>\n		<p><strong>Login:</strong> <a href=\"<%MASTER_URL%>\"><%MASTER_URL%></a><br />\n			<strong>Email:</strong> <%EMAIL%><br />\n			<strong>Password:</strong> <%NEW_PASS%></p>\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>'),
	(8,'admin_newaccount','<html>\n	<head>\n		<title><%HOST_NAME%> - Account Created</title>\n	</head>\n	<body>\n		<center><h1><%HOST_NAME%> - Account Created</h1></center>\n		<p>Hello there! This email is to inform you that an account has been created for you on <%HOST_NAME%>.</p>\n		<p><strong>Login:</strong> <a href=\"<%MASTER_URL%>\"><%MASTER_URL%></a><br />\n			<strong>Email:</strong> <%EMAIL%><br />\n			<strong>Password:</strong> <%PASS%></p>\n		<p>Thanks!<br /><%HOST_NAME%></p>\n	</body>\n</html>');

/*!40000 ALTER TABLE `acp_email_templates` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table acp_settings
# ------------------------------------------------------------

DROP TABLE IF EXISTS `acp_settings`;

CREATE TABLE `acp_settings` (
  `setting_ref` char(25) NOT NULL DEFAULT '',
  `setting_val` tinytext
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `acp_settings` WRITE;
/*!40000 ALTER TABLE `acp_settings` DISABLE KEYS */;

INSERT INTO `acp_settings` (`setting_ref`, `setting_val`)
VALUES
	('company_name','PufferPanel Server Panel'),
	('node_url','http://node2.example.com/node/'),
	('master_url','http://example.com/panel/'),
	('cookie_website',NULL),
	('postmark_api_key','NULL'),
	('mandrill_api_key','NULL'),
	('mailgun_api_key','NULL'),
	('sendmail_email','webmaster@localhost'),
	('main_website','http://example.com/panel/'),
	('sendmail_method','php'),
	('captcha_pub','6LdSzuYSAAAAAHkmq8LlvmhM-ybTfV8PaTgyBDII'),
	('captcha_priv','6LdSzuYSAAAAAISSAYIJrFGGGJHi5a_V3hGRvIAz'),
	('assets_url','http://example.com/panel/assets/'),
	('use_api','0'),
	('api_key','NULL'),
	('api_allowed_ips','*'),
	('api_module_controls_all','0');

/*!40000 ALTER TABLE `acp_settings` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table backup_datastore
# ------------------------------------------------------------

DROP TABLE IF EXISTS `backup_datastore`;

CREATE TABLE `backup_datastore` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `server` text NOT NULL,
  `backup_pattern` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table backups
# ------------------------------------------------------------

DROP TABLE IF EXISTS `backups`;

CREATE TABLE `backups` (
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
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table nodes
# ------------------------------------------------------------

DROP TABLE IF EXISTS `nodes`;

CREATE TABLE `nodes` (
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
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table servers
# ------------------------------------------------------------

DROP TABLE IF EXISTS `servers`;

CREATE TABLE `servers` (
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
  `path` text NOT NULL,
  `date_added` int(15) NOT NULL,
  `server_ip` varchar(50) NOT NULL DEFAULT '',
  `server_port` int(11) NOT NULL,
  `ftp_host` tinytext NOT NULL,
  `ftp_user` tinytext NOT NULL,
  `ftp_pass` tinytext NOT NULL,
  `backup_file_limit` int(20) NOT NULL,
  `backup_disk_limit` int(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
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
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table whmcs_suspend_data
# ------------------------------------------------------------

DROP TABLE IF EXISTS `whmcs_suspend_data`;

CREATE TABLE `whmcs_suspend_data` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `server_id` int(11) NOT NULL,
  `whmcs_server_id` int(11) NOT NULL,
  `old_password` varchar(40) NOT NULL DEFAULT '',
  `unsuspended` int(11) NOT NULL DEFAULT '0' COMMENT '0 = False, 1 = True',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
