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
if(file_exists('install.lock'))
	exit('Installer is Locked.');
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<link rel="stylesheet" href="../../assets/css/bootstrap.css">
	<title>PufferPanel Installer</title>
</head>
<body>
	<div class="container">
		<div class="alert alert-danger">
			<strong>WARNING:</strong> Do not run this version on a live environment! There are known security holes that we are working on getting patched. This is extremely beta software and this version is to get the features in place while we work on security enhancements.
		</div>
		<div class="navbar navbar-default">
			<div class="navbar-header">
				<a class="navbar-brand" href="#">Install PufferPanel</a>
			</div>
		</div>
		<div class="col-12">
			<div class="row">
				<div class="col-2"></div>
				<div class="col-8">
					<div class="alert alert-danger">When this installer finishes please manually change the permissions back to 0755 for the /src/framework folder, and delete this installer. Please set the permissions on the global_configuration.php file to 0444.</div>
					<p>This script will guide you through the process for setting up PufferPanel on your server. Please ensure that you have installed all of the dependencies required or this install will fail.</p>
					<p><?php
					
						/* Permissions Checking */
						echo (substr(sprintf('%o', fileperms('../../../src/framework/configuration.php.dist')), -4) == "0666") ? '<small class="text-success">/src/framework/configuration.php.dist is correctly CHMOD\'d 0666</small><br />' : '<small class="text-danger">/src/framework/configuration.php.dist is improperly CHMOD\'d. It is '.substr(sprintf('%o', fileperms('../../../src/framework/configuration.php.dist')), -4).' and should be 0666.</small><br />';
						
						echo (substr(sprintf('%o', fileperms('../../../src/framework')), -4) == "0777") ? '<small class="text-success">/src/framework is correctly CHMOD\'d 0777</small><br />' : '<small class="text-danger">/src/framework is improperly CHMOD\'d. It is '.substr(sprintf('%o', fileperms('../../../src/framework')), -4).' and should be 0777.</small><br />';
						
						echo (substr(sprintf('%o', fileperms('../install')), -4) == "0777") ? '<small class="text-success">/panel/admin/install is correctly CHMOD\'d 0777</small><br />' : '<small class="text-danger">/panel/admin/install is improperly CHMOD\'d. It is '.substr(sprintf('%o', fileperms('../install')), -4).' and should be 0777.</small><br />';
						
						echo (substr(sprintf('%o', fileperms('do')), -4) == "0777") ? '<small class="text-success">/panel/admin/install/do is correctly CHMOD\'d 0777</small>' : '<small class="text-danger">/panel/admin/install/do is improperly CHMOD\'d. It is '.substr(sprintf('%o', fileperms('do')), -4).' and should be 0777.</small>';
						
					?></p>
					<hr />
					<p><?php

						$continue = true;
                        /* Check for Required Dependencies */
                        $list = array(
                            'curl',
                            'hash',
                            'openssl',
                            'mcrypt',
                            'PDO',
                            'pdo_mysql',
                            'ssh2'
                        );
                        echo "\n";
                        foreach($list as $ext) {
                            
                            echo (extension_loaded($ext)) ? '<small class="text-success">The php-'.$ext.' extension was loaded.</small><br />' : '<small class="text-danger"><strong>The php-'.$ext.' extension was not loaded.</strong></small><br />';
                            
                            if(!extension_loaded($ext))
                            	$continue = false;
                            
                        }
					?></p>
					<hr />
					<p><?php
					
                        $functions = array(
                        	'fopen',
                        	'fclose',
                        	'fwrite',
                        	'session_start',
                        	'socket_set_option',
                        	'socket_send',
                        	'socket_connect',
                        	'socket_create',
                        	'stream_set_timeout',
                        	'fsockopen',
                        	'crypt',
                        	'hash'
                        );
                        
                        foreach($functions as $fct) {
                            
                            echo (function_exists($fct)) ? '<small class="text-success">'.$fct.'() is enabled.</small><br />' : '<small class="text-danger"><strong>'.$fct.'() is not enabled.</strong></small><br />';
                            
                            if(!function_exists($fct))
                            	$continue = false;
                            
                        }

                    ?></p>
                    <hr />
                    <?php echo ($continue === true) ? '<a href="do/start.php">Continue &rarr;</a>' : '<p>Please fix missing extensions and functions before continuing.</p>'; ?>
				</div>
				<div class="col-2"></div>
			</div>
		</div>
		<div class="footer">
			<div class="col-8 nopad"><p>PufferPanel is licensed under a <a href="https://github.com/DaneEveritt/PufferPanel/blob/master/LICENSE">GPL-v3 License</a>.<br />Running Version 0.7.0 Alpha RC2 distributed by <a href="http://kelp.in">Kelpin' Systems</a>.</p></div>
		</div>
	</div>
</body>
</html>
