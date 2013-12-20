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
				    <h1>Hashing Information</h1>
                    <p>PufferPanel encrypts sensitive SFTP password information with AES-256 encryption prior to storing it in the database. In order to do this encryption you must provide an encryption key on all of the servers running PufferPanel (master and nodes). You can generate a key below that should be placed in a file that PHP can access, but is outside of the public web root. We suggest /etc/hashfile.txt. <strong>Do not move on to the next step until you have done so, and have entered the specific file location for this hash.</strong></p>
                    <div class="stripe-separator"><!--  --></div>
                    <pre><?php
/* Make File */
$keyset  = "abcdefghijklmABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*-=+";
$randkey = "";

for ($i=0; $i<64; $i++){
    $randkey .= substr($keyset, rand(0, strlen($keyset)-1), 1);
}

echo $randkey;
                        ?></pre>
                <br /><br />
                    <?php

                        if(isset($_POST['hash_do'])){
                            
                            if(fopen($_POST['hash'], 'r')){
                                
                                $fp = fopen('../../../core/framework/master_configuration.php', 'a+');
                                fwrite($fp, "
                                
if(!defined('HASH'))
	define('HASH', '".$_POST['hash']."');");
                                fclose($fp);
                                
                                exit('<meta http-equiv="refresh" content="0;url=account.php"/>');
                                
                            }else{
                             
                                echo '<div class="error-box round">We were unable to access that file location.</div>';
                                
                            }
                            
                        }

                    ?>
                    <form action="hash.php" method="post">
                        <p>
                            <label for="hash">Hash File Location</label>
                            <input type="text" name="hash" placeholder="/etc/hashfile.txt" class="round default-width-input" />
                        </p>
                        <input type="submit" name="hash_do" value="Next Step" class="round blue ic-right-arrow" />
                    </form><br />
                <a href="hash.php">Regenerate Key</a>
				</div>
            </div>
		</div>
	</div>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.3 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>