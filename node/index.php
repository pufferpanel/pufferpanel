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
					<a href="index.php" class="list-group-item active">Overview</a>
					<a href="console.php" class="list-group-item">Live Console</a>
					<a href="files/index.php" class="list-group-item">File Manager</a>
					<a href="backup.php" class="list-group-item">Backup Manager</a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Server Settings</strong></a>
					
					<a href="settings.php" class="list-group-item">Modpack Management</a>
					<a href="plugins/index.php" class="list-group-item">Server Plugins</a>
				</div>
			</div>
			<div class="col-9">
				<div class="col-12">
					<h3 class="nopad">Players Online</h3><hr />
					<div id="server_players">
						<p id="server_players_loading" style="margin: 1.25em;text-align: center;" class="text-muted"><i class="fa fa-cog fa-3x fa-spin"></i></p>
					</div>
				</div>
				<div class="col-12">
					<h3>Disk Space Used</h3><hr />
					<div id="server_stats">
						<p id="server_stats_loading" style="margin: 1.25em;text-align: center;" class="text-muted"><i class="fa fa-cog fa-3x fa-spin"></i></p>
					</div>
				</div>
				<div class="col-12">
					<h3>Server Information</h3><hr />
					<div id="server_info">
						<p id="server_info_loading" style="margin: 1.25em;text-align: center;" class="text-muted"><i class="fa fa-cog fa-3x fa-spin"></i></p>
					</div>
				</div>
			</div>
		</div>
		<div class="footer">
			<?php include('assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$(document).ready(function(){
			$.ajax({
				type: "POST",
				url: "core/ajax/overview/data.php",
				data: { command: 'info' },
			  		success: function(data) {
						$("#server_info_loading").slideUp("slow", function(){
							$("#server_info").hide();
							$("#server_info").html(data);
							$("#server_info").slideDown("slow");				
						});
			 		}
			});
			$.ajax({
				type: "POST",
				url: "core/ajax/overview/data.php",
				data: { command: 'players' },
			  		success: function(data) {
						$("#server_players_loading").slideUp("slow", function(){
							$("#server_players").hide();
							$("#server_players").html(data);
							$("#server_players").slideDown("slow");
							$("img[data-toggle='tooltip']").tooltip();			
						});
			 		}
			});
			$.ajax({
				type: "POST",
				url: "core/ajax/overview/data.php",
				data: { command: 'stats' },
			  		success: function(data) {
						$("#server_stats_loading").slideUp("slow", function(){
							$("#server_stats").hide();
							$("#server_stats").html(data);
							$("#server_stats").slideDown("slow");				
						});
			 		}
			});
		});
	</script>
</body>
</html>