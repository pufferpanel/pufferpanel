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
	<meta charset="utf-8">
	<title>PufferPanel - Install</title>
	
	<!-- Stylesheets -->
	<link href='http://fonts.googleapis.com/css?family=Droid+Sans:400,700' rel='stylesheet'>
	<link rel="stylesheet" href="../../assets/css/style.css">
	
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
				    <h1>Install PufferPanel on your Server</h1>
                    <p>This script will guide you through the process for setting up PufferPanel on your server. Please ensure that you have installed all of the dependencies required or this install will fail.</p>
                    <p><strong><span style="color:red;">!!IMPORTANT!!</span> Please ensure that you have CHMOD'ed the /master/core/framework/ folder and this installer folder (and the /do/ directory) to 0777 for installation. The installer will reset the permission when it is done. </strong></p>
                    <div class="stripe-separator"></div>
                    	<p>
	                    <?php
	                    
	                    	/* Permissions Checking */
	                    	echo (substr(sprintf('%o', fileperms('../../core/framework/master_configuration.php.dist')), -4) == "0666") ? '<span style="color:green;">/core/framework/master_configuration.php.dist is correctly CHMOD\'d 0666</span><br />' : '<span style="color:red;"><strong>/core/framework/master_configuration.php.dist is improperly CHMOD\'d. It is '.substr(sprintf('%o', fileperms('../../core/framework/master_configuration.php.dist')), -4).' and should be 0666.</strong></span><br />';
	                    	
	                    	echo (substr(sprintf('%o', fileperms('../../core/framework')), -4) == "0777") ? '<span style="color:green;">/core/framework is correctly CHMOD\'d 0777</span><br />' : '<span style="color:red;"><strong>/core/framework is improperly CHMOD\'d. It is '.substr(sprintf('%o', fileperms('../../core/framework')), -4).' and should be 0777.</strong></span><br />';
	                    	
	                    	echo (substr(sprintf('%o', fileperms('../install')), -4) == "0777") ? '<span style="color:green;">/admin/install is correctly CHMOD\'d 0777</span><br />' : '<span style="color:red;"><strong>/admin/install is improperly CHMOD\'d. It is '.substr(sprintf('%o', fileperms('../install')), -4).' and should be 0777.</strong></span><br />';
	                    	
	                    	echo (substr(sprintf('%o', fileperms('do')), -4) == "0777") ? '<span style="color:green;">/admin/install/do is correctly CHMOD\'d 0777</span>' : '<span style="color:red;"><strong>/admin/install/do is improperly CHMOD\'d. It is '.substr(sprintf('%o', fileperms('do')), -4).' and should be 0777.</strong></span>';
	                    	
	                    ?>
                    	</p>
                    <div class="stripe-separator"><!--  --></div>
                    <pre>
                        <?php

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
                                
                                echo (extension_loaded($ext)) ? $ext." extension was loaded.\n" : $ext." extension was not loaded.\n";
                                
                            }
                            
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
                                
                                echo (function_exists($fct)) ? $fct." is enabled.\n" : $fct." is not enabled.\n";
                                
                            }

                        ?>
                    </pre>
                    <?php echo ($continue === true) ? '<a href="do/start.php">Start Install &rarr;</a>' : '<p>Please fix missing extensions and functions before continuing.</p>'; ?>
				</div>
            </div>
		</div>
	</div>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.3 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>