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
require_once('../core/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false){

	$core->page->redirect($core->settings->get('master_url').'index.php');
	exit();
}
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('../assets/include/header.php'); ?>
	<title>PufferPanel - Manage Your Server</title>
</head>
<body>
	<div class="container">
		<?php include('../assets/include/navbar.php'); ?>
		<div class="row">
			<div class="col-3">
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Account Actions</strong></a>
					<a href="../account.php" class="list-group-item">Settings</a>
					<a href="../servers.php" class="list-group-item">My Servers</a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Server Actions</strong></a>
					<a href="index.php" class="list-group-item active">Overview</a>
					<a href="console.php" class="list-group-item">Live Console</a>
					<a href="files/index.php" class="list-group-item">File Manager</a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Server Settings</strong></a>
					
					<a href="settings.php" class="list-group-item">Server Management</a>
					<a href="plugins/index.php" class="list-group-item">Server Plugins</a>
				</div>
			</div>
			<div class="col-9">
				<div class="col-12">
					<h3 class="nopad">Stats Overview</h3><hr />
					<h5>CPU Usage</h5>
						<div class="progress">
							<div class="progress-bar" id="cpu_bar" style="width:100%;max-width:100%;">Gathering...</div>
						</div>
					<h5>Memory Usage</h5>
						<div class="progress">
							<div class="progress-bar" id="memory_bar" style="width:100%; max-width:100%;">Gathering...</div>
						</div>
					<h5>Players Online</h5>
						<span id="player_list">
							<p class="text-muted">No players are currently online.</p>
						</span>
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
			<?php include('../assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$(window).load(function(){
			
			var socket = io.connect('http://<?php echo $core->server->nodeData('sftp_ip'); ?>:8031/<?php echo $core->server->getData('gsd_id'); ?>');
			socket.on('process', function (data) {
				$("#cpu_bar").css('width', data.process.cpu + '%');
				$("#cpu_bar").html(data.process.cpu + '%');
				$("#memory_bar").css('width', (data.process.memory / (1024 * 1024)).toFixed(0) + '%');
				$("#memory_bar").html((data.process.memory / (1024 * 1024)).toFixed(0) + 'MB / <?php echo $core->server->getData('max_ram'); ?>MB');
			});
			socket.on('query', function (data) {
				$("#live_console").val($("#live_console").val() + data.process.memory);
				$('#live_console').scrollTop($('#live_console')[0].scrollHeight - $('#live_console').height());
			});
			$.ajax({
				type: "POST",
				url: "ajax/overview/data.php",
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
				url: "ajax/overview/data.php",
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