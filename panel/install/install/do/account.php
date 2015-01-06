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
namespace PufferPanel\Core;
use \ORM as ORM;
use \PDO as PDO;
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
				<a class="navbar-brand" href="#">Install PufferPanel - Your Account</a>
			</div>
		</div>
		<div class="col-12">
			<div class="row">
				<div class="col-2"></div>
				<div class="col-8">
					<p>You've reached the final step of the process. Please fill out the information below to create your admin account. After finishing this step you will be redirected to the login page where you will be able to access the Admin Control Panel to add nodes, users, and servers. Thank you for installing PufferPanel on your server. Please contact us on IRC <code>(irc.esper.net/#pufferpanel)</code> if you encounter any problems or have questions, comments, or concerns.</p>
					?>
					<form action="account" method="post">
					    <div class="form-group">
					    	<label for="email" class="control-label">Email</label>
					    	<div>
					    		<input type="text" class="form-control" name="email" autocomplete="off" />
					    	</div>
					    </div>
					    <div class="form-group">
					    	<label for="username" class="control-label">Username</label>
					    	<div>
					    		<input type="text" class="form-control" name="username" autocomplete="off" />
					    	</div>
					    </div>
					    <div class="form-group">
					    	<label for="password" class="control-label">Password</label>
					    	<div>
					    		<input type="password" class="form-control" name="password" autocomplete="off" />
					    	</div>
					    </div>
					    <div class="form-group">
					    	<div>
					    		<input type="submit" class="btn btn-primary" name="do_account" value="Finish &rarr;" />
					    	</div>
					    </div>
					</form>
				</div>
				<div class="col-2"></div>
			</div>
		</div>
		<div class="footer">
			<div class="col-8 nopad"><p>PufferPanel is licensed under a <a href="https://github.com/PufferPanel/PufferPanel/blob/master/LICENSE">GPL-v3 License</a>.<br />Running <?php echo trim(file_get_contents('../../../../src/versions/current')).' ('.substr(trim(file_get_contents('../../../../.git/HEAD')), 0, 8).')'; ?> distributed by <a href="http://pufferpanel.com">PufferPanel Development</a>.</p></div>
		</div>
	</div>
</body>
</html>