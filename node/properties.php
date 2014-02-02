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
					<a href="files/index.php" class="list-group-item">File Manager</a>
					<a href="backup.php" class="list-group-item">Backup Manager</a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Server Settings</strong></a>
					<a href="properties.php" class="list-group-item active">Server Properties</a>
					<a href="settings.php" class="list-group-item">Modpack Management</a>
					<a href="plugins/index.php" class="list-group-item">Server Plugins</a>
				</div>
			</div>
			<div class="col-9">
				<h3 class="nopad">Server Settings</h3><hr />
				<form action="core/ajax/settings/update_settings.php" method="POST">
					<div class="row well">
						<div class="form-group">
							<label for="motd" class="control-label">Server MOTD</label>
							<div>
								<input type="text" name="motd" class="form-control" />
							</div>
						</div>
						<div class="form-group">
							<label for="level-name" class="control-label">Level Name</label>
							<div>
								<input type="text" name="level-name" class="form-control" />
							</div>
						</div>
					</div>
					<div class="row well">
						<div class="form-group col-4 nopad">
							<label for="spawn-protection" class="control-label">Spawn Protection Distance</label>
							<div class="input-group">
								<input type="number" name="spawn-protection" min="1" max="64" class="form-control"/>
								<span class="input-group-addon">Blocks</span>
							</div>
						</div>
						<div class="form-group col-4 nopad">
							<label for="view-distance" class="control-label">View Distance</label>
							<div class="input-group">
								<input type="number" name="view-distance" min="1" max="15" class="form-control"/>
								<span class="input-group-addon">Chunks</span>
							</div>
						</div>
						<div class="form-group col-4 nopad">
							<label for="max-players" class="control-label">Max Players</label>
							<div>
								<input type="number" name="max-players" min="1" class="form-control"/>
							</div>
						</div>
					</div>
					<div class="row well">
						<div class="form-group col-4 nopad">
							<label for="generate-structures" class="control-label">Generate Structures</label>
							<div>
								<select class="form-control" name="generate-structures">
									<option value="true">True</option>
									<option value="false">False</option>
								</select>
							</div>
						</div>
						<div class="form-group col-4 nopad">
							<label for="spawn-monsters" class="control-label">Spawn Monsters</label>
							<div>
								<select class="form-control" name="spawn-monsters">
									<option value="true">True</option>
									<option value="false">False</option>
								</select>
							</div>
						</div>
						<div class="form-group col-4 nopad">
							<label for="allow-nether" class="control-label">Allow Nether</label>
							<div>
								<select class="form-control" name="allow-nether">
									<option value="true">True</option>
									<option value="false">False</option>
								</select>
							</div>
						</div>
						<div class="form-group col-4 nopad">
							<label for="pvp" class="control-label">Enable PvP</label>
							<div>
								<select class="form-control" name="pvp">
									<option value="true">True</option>
									<option value="false">False</option>
								</select>
							</div>
						</div>
						<div class="form-group col-4 nopad">
							<label for="spawn-animals" class="control-label">Spawn Animals</label>
							<div>
								<select class="form-control" name="spawn-animals">
									<option value="true">True</option>
									<option value="false">False</option>
								</select>
							</div>
						</div>
						<div class="form-group col-4 nopad">
							<label for="spawn-npcs" class="control-label">Spawn NPCs</label>
							<div>
								<select class="form-control" name="spawn-npcs">
									<option value="true">True</option>
									<option value="false">False</option>
								</select>
							</div>
						</div>
						<div class="form-group col-4 nopad">
							<label for="white-list" class="control-label">Enable Whitelist</label>
							<div>
								<select class="form-control" name="white-list">
									<option value="true">True</option>
									<option value="false">False</option>
								</select>
							</div>
						</div>
						<div class="form-group col-4 nopad">
							<label for="allow-flight" class="control-label">Allow Flight</label>
							<div>
								<select class="form-control" name="allow-flight">
									<option value="true">True</option>
									<option value="false">False</option>
								</select>
							</div>
						</div>
						<div class="form-group col-4 nopad">
							<label for="announce-player-achievements" class="control-label">Announce Player Achievements</label>
							<div>
								<select class="form-control" name="announce-player-achievements">
									<option value="true">True</option>
									<option value="false">False</option>
								</select>
							</div>
						</div>
					</div>
					<div class="row well">
						<input type="submit" name="update" value="Update Settings" class="btn btn-primary btn-sm" />
					</div>
				</form>
			</div>
		</div>
		<div class="footer">
			<?php include('assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
	</script>
</body>
</html>