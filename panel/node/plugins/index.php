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
session_start();
require_once('../../core/framework/framework.core.php');

$filesIncluded = true;

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false){

	Page\components::redirect($core->settings->get('master_url').'index.php?login');
	exit();
}
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('../../assets/include/header.php'); ?>
	<title>PufferPanel - Manage Your Server</title>
</head>
<body>
	<div class="container">
		<?php include('../../assets/include/navbar.php'); ?>
		<div class="row">
			<div class="col-3">
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Account Actions</strong></a>
					<a href="../../account.php" class="list-group-item">Settings</a>
					<a href="../../servers.php" class="list-group-item">My Servers</a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Server Actions</strong></a>
					<a href="../index.php" class="list-group-item">Overview</a>
					<a href="../console.php" class="list-group-item">Live Console</a>
					<a href="../files/index.php" class="list-group-item">File Manager</a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Server Settings</strong></a>
					
					<a href="../settings.php" class="list-group-item">Server Management</a>
					<a href="index.php" class="list-group-item active">Server Plugins</a>
				</div>
			</div>
			<div class="col-9">
				<div class="panel panel-primary">
					<div class="panel-heading">
						<h3 class="panel-title">Plugin Installer</h3>
					</div>
					<div class="panel-body">
						<p>Welcome to your built-in plugin installer. This installer gathers results and data from <a href="http://bukget.org" target="_blank">Bukget</a> and displays them here in an easy to use manner. The automatic installer can be buggy at times, especially for complex plugins. We recommend downloading plugins to your computer first if possible.</p>
					</div>
				</div>
				<form action="search.php" method="get">
					<legend>Search Plugin Repository</legend>
					<fieldset>
						<div class="form-group">
							<label for="email" class="control-label">Plugin Name</label>
							<div>
								<input type="text" class="form-control" name="term" placeholder="e.g. worldedit" />
							</div>
						</div>
						<div class="form-group">
							<div>
								<input type="submit" class="btn btn-primary" value="Search" />
								<input type="reset" class="btn btn-default" value="Clear" />
							</div>
						</div>
					</fieldset>
				</form>
			</div>
		</div>
		<div class="footer">
			<?php include('../../assets/include/footer.php'); ?>
		</div>
	</div>
</body>
</html>