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
require_once('core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), $core->framework->auth->getCookie('pp_server_hash')) === false){

	$core->framework->page->redirect($core->framework->settings->get('master_url').'index.php');
	exit();
}

?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('assets/include/header.php'); ?>
	<title>PufferPanel - Manage Your Server</title>
</head>
<body>
	<div class="container">
		<?php include('assets/include/navbar.php'); ?>
		<div class="row">
			<div class="col-3">
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Account Actions</strong></a>
					<a href="<?php echo $core->framework->settings->get('master_url'); ?>account.php" class="list-group-item">Settings</a>
					<a href="<?php echo $core->framework->settings->get('master_url'); ?>servers.php" class="list-group-item">My Servers</a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Server Actions</strong></a>
					<a href="index.php" class="list-group-item">Overview</a>
					<a href="console.php" class="list-group-item">Live Console</a>
					<a href="settings.php" class="list-group-item active">Server Settings</a>
					<a href="plugins/index.php" class="list-group-item">Server Plugins</a>
					<a href="files/index.php" class="list-group-item">File Manager</a>
					<a href="backup.php" class="list-group-item">Backup Manager</a>
				</div>
			</div>
			<div class="col-9">
				<h3 class="nopad">Server Modpack Management</h3><hr />
				<form action="" method="post">
					<div class="form-group col-8 nopad">
						<div>
							<select class="form-control">
								<option>Some Modpack</option>
							</select>
						</div>
					</div>
					<div class="form-group col-4 nopad-right">
						<div>
							<input type="submit" value="Install Modpack" class="btn btn-primary" />
						</div>
					</div>
				</form>
			</div>
		</div>
		<div class="footer">
			<?php include('assets/include/footer.php'); ?>
		</div>
	</div>
</body>
</html>