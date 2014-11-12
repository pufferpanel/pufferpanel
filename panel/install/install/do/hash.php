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
					<p>PufferPanel encrypts some sensitive information (not account passwords) with <code>AES-CBC-256</code> encryption prior to storing it in the database. In order to do this you must provide an encryption key that PufferPanel can use. Please execute the command below to create this.</p>
					<?php

                        if(isset($_POST['hash_do'])){

                            if(@fopen('/etc/HASHFILE', 'r')){

                                $fp = fopen('../../../../src/core/configuration.php', 'a+');
                                fwrite($fp, "

if(!defined('HASH'))
	define('HASH', '/etc/HASHFILE');");
                                fclose($fp);

                                exit('<meta http-equiv="refresh" content="0;url=account.php"/>');

                            }else{

                                echo '<div class="alert alert-danger">We were unable to access that file location.</div>';

                            }

                        }

						/* Make File */
						$keyset  = "abcdefghijklmABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*-=+";
						$randkey = "";

						for ($i=0; $i<48; $i++){
						    $randkey .= substr($keyset, rand(0, strlen($keyset)-1), 1);
						}

					?>
					<div class="well">
						<h5>Execute the Command Below on your Server:</h5>
						<p><code>sudo echo "<?php echo $randkey; ?>" > /etc/HASHFILE</code></p>
					</div>
					<hr />
					<form action="hash.php" method="post">
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
            <div class="col-8 nopad"><p>PufferPanel is licensed under a <a href="https://github.com/DaneEveritt/PufferPanel/blob/master/LICENSE">GPL-v3 License</a>.<br />Running <?php echo trim(file_get_contents('../../../../src/versions/current')).' ('.substr(trim(file_get_contents('../../../../.git/HEAD')), 0, 8).')'; ?> distributed by <a href="http://pufferpanel.com">PufferPanel Development</a>.</p></div>
		</div>
	</div>
</body>
</html>
