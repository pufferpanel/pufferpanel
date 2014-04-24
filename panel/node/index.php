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

	Page\components::redirect($core->settings->get('master_url').'index.php');
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
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.acc_actions'); ?></strong></a>
					<a href="../account.php" class="list-group-item"><?php echo $_l->tpl('sidebar.settings'); ?></a>
					<a href="../servers.php" class="list-group-item"><?php echo $_l->tpl('sidebar.servers'); ?></a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.server_acc'); ?></strong></a>
					<a href="index.php" class="list-group-item active"><?php echo $_l->tpl('sidebar.overview'); ?></a>
					<a href="console.php" class="list-group-item"><?php echo $_l->tpl('sidebar.console'); ?></a>
					<a href="files/index.php" class="list-group-item"><?php echo $_l->tpl('sidebar.files'); ?></a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.server_sett'); ?></strong></a>
					
					<a href="settings.php" class="list-group-item"><?php echo $_l->tpl('sidebar.manage'); ?></a>
					<a href="plugins/index.php" class="list-group-item"><?php echo $_l->tpl('sidebar.plugins'); ?></a>
				</div>
			</div>
			<div class="col-9">
				<div class="col-12">
					<h3 class="nopad"><?php echo $_l->tpl('node.overview.stats_h1'); ?></h3><hr />
					<div id="online_notice" class="alert alert-info"><i class="fa fa-spinner fa-spin"></i> <?php echo $_l->tpl('node.overview.collect_usage'); ?></div>
					<span id="toggle_on" style="display:none;">
						<h5><?php echo $_l->tpl('node.overview.cpu_h5'); ?></h5>
							<div class="progress">
								<div class="progress-bar" id="cpu_bar" style="width:100%;max-width:100%;"></div>
							</div>
						<h5><?php echo $_l->tpl('node.overview.memory_h5'); ?></h5>
							<div class="progress">
								<div class="progress-bar" id="memory_bar" style="width:100%; max-width:100%;"></div>
							</div>
					</span>
					<h5><?php echo $_l->tpl('node.overview.players_h5'); ?></h5>
						<div id="players_notice" class="alert alert-info"><i class="fa fa-spinner fa-spin"></i> <?php echo $_l->tpl('node.overview.collect_players'); ?></div>
						<span id="toggle_players" style="display:none;">
							<p class="text-muted"><?php echo $_l->tpl('node.overview.no_players'); ?></p>
						</span>
				</div>
				<div class="col-12">
					<h3><?php echo $_l->tpl('node.overview.disk_h1'); ?></h3><hr />
					<div id="server_stats">
						<div class="alert alert-info"><i class="fa fa-spinner fa-spin"></i> <?php echo $_l->tpl('node.overview.collect_disk'); ?></div>
					</div>
				</div>
				<div class="col-12">
					<h3><?php echo $_l->tpl('node.overview.information_h1'); ?></h3><hr />
					<table class="table table-striped table-bordered table-hover">
						<tbody>
							<tr>
								<td><strong><?php echo $_l->tpl('string.connection'); ?></strong></td>
								<td><?php echo $core->server->getData('server_ip').':'.$core->server->getData('server_port'); ?></td>
							</tr>
							<tr>
								<td><strong><?php echo $_l->tpl('string.node'); ?></strong></td>
								<td><?php echo $core->settings->nodeName($core->server->getData('node')); ?></td>
							</tr>
							<tr>
								<td><strong><?php echo $_l->tpl('string.mem_alloc'); ?></strong></td>
								<td><?php echo $core->server->getData('max_ram').' MB'; ?></td>
							</tr>
							<tr>
								<td><strong><?php echo $_l->tpl('string.disk_alloc'); ?></strong></td>
								<td><?php echo $core->server->getData('disk_space').' MB'; ?></td>
							</tr>
						</tbody>
					</table>
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
				if($("#online_notice").is(":visible")){
					$("#online_notice").hide();
					$("#toggle_on").show();
				}
			});
			socket.on('query', function (data) {
				if($("#players_notice").is(":visible")){
					$("#players_notice").hide();
					$("#toggle_players").show();
				}
				if(data.query.players.length !== 0){
					$("#toggle_players").html("");
					$.each(data.query.players, function(id, name) {
						$("#toggle_players").append('<img data-toggle="tooltip" src="http://i.fishbans.com/helm/'+name+'/32" title="'+name+'" style="padding: 0 2px 6px 0;"/>');
					});
				}else{
					$("#toggle_players").html('<p class="text-muted">No players are currently online.</p>');
				}
				$("img[data-toggle='tooltip']").tooltip();
			});
			$.ajax({
				type: "POST",
				url: "ajax/overview/data.php",
				data: { command: 'stats' },
			  		success: function(data) {
						$("#server_stats").html(data);
			 		}
			});
		});
	</script>
</body>
</html>