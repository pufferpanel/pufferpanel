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
	<link rel="stylesheet" href="../../../assets/css/bootstrap.css">
	<title>PufferPanel Installer</title>
</head>
<body>
	<div class="container">
		<div class="alert alert-danger">
			<strong>WARNING:</strong> Do not run this version on a live environment! There are known security holes that we are working on getting patched. This is extremely beta software and this version is to get the features in place while we work on security enhancements.
		</div>
		<div class="navbar navbar-default">
			<div class="navbar-header">
				<a class="navbar-brand" href="#">Install PufferPanel - Encryption</a>
			</div>
		</div>
		<div class="col-12">
			<div class="row">
				<div class="col-2"></div>
				<div class="col-8">
					<div class="alert alert-danger">Do not move on to the next step until you have done so, and have entered the specific file location for this hash.</div>
					<p>PufferPanel encrypts sensitive SFTP password information with <code>AES-CBC-256</code> encryption prior to storing it in the database. In order to do this encryption you must provide an encryption key on all of the servers running PufferPanel. You can generate a key below that should be placed in a file that PHP can access, but is outside of the public web root. We suggest <code>/etc/hashfile.txt</code>.</p>
					<?php

                        if(isset($_POST['hash_do'])){
                            
                            if(fopen($_POST['hash'], 'r')){
                                
                                $fp = fopen('../../../core/framework/configuration.php', 'a+');
                                fwrite($fp, "
                                
if(!defined('HASH'))
	define('HASH', '".$_POST['hash']."');");
                                fclose($fp);
                                
                                exit('<meta http-equiv="refresh" content="0;url=account.php"/>');
                                
                            }else{
                             
                                echo '<div class="alert alert-danger">We were unable to access that file location.</div>';
                                
                            }
                            
                        }

                    ?>
					<h5><strong>Here is a randomly generated encryption key you can use:</strong></h5>
					<code>
					<?php
						/* Make File */
						$keyset  = "abcdefghijklmABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*-=+";
						$randkey = "";
						
						for ($i=0; $i<64; $i++){
						    $randkey .= substr($keyset, rand(0, strlen($keyset)-1), 1);
						}
						
						echo $randkey;
					?>
					</code>
					<hr />
					<form action="hash.php" method="post">
					    <div class="form-group">
					    	<label for="hash" class="control-label">Hash File Location</label>
					    	<div>
					    		<input type="text" class="form-control" name="hash" placeholder="/etc/hashfile.txt" autocomplete="off" />
					    	</div>
					    </div>
					    <div class="form-group">
					    	<div>
					    		<input type="submit" class="btn btn-primary" name="hash_do" value="Continue &rarr;" />
					    	</div>
					    </div>
					</form>
				</div>
				<div class="col-2"></div>
			</div>
		</div>
		<div class="footer">
			<div class="col-8 nopad"><p>PufferPanel is licensed under a <a href="https://github.com/DaneEveritt/PufferPanel/blob/master/LICENSE">GPL-v3 License</a>.<br />Running Version 0.6.0.1 Beta distributed by <a href="http://kelp.in">Kelpin' Systems</a>.</p></div>
		</div>
	</div>
</body>
</html>