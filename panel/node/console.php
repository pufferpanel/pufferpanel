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

$filesIncluded = true;

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
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.acc_actions'); ?></strong></a>
					<a href="../account.php" class="list-group-item"><?php echo $_l->tpl('sidebar.settings'); ?></a>
					<a href="../servers.php" class="list-group-item"><?php echo $_l->tpl('sidebar.servers'); ?></a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.server_acc'); ?></strong></a>
					<a href="index.php" class="list-group-item"><?php echo $_l->tpl('sidebar.overview'); ?></a>
					<a href="console.php" class="list-group-item active"><?php echo $_l->tpl('sidebar.console'); ?></a>
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
					<textarea id="live_console" class="form-control console" readonly="readonly"></textarea>
				</div>
				<div id="box"></div>
				<div class="col-12">
					<div class="alert alert-danger text_highlighted" style="display:none;margin: 15px 0 -5px 0;"><?php echo $_l->tpl('node.console.scrollstop_alert'); ?></div>
				</div>
				<div class="col-6">
					<hr />
					<form action="#" method="post" id="console_command">
						<fieldset>
							<div class="input-group">
								<input type="text" class="form-control" name="command" id="ccmd" placeholder="<?php echo $_l->tpl('node.console.command_help'); ?>" />
								<span class="input-group-btn">
									<button id="sending_command" class="btn btn-primary">&rarr;</button>
									<button class="btn btn-link" data-toggle="modal" data-target="#pauseConsole" id="pause_console"><small><i class="fa fa-pause"></i></small></button>
									
								</span>
							</div>
						</fieldset>
					</form>
					<div class="alert alert-danger" id="sc_resp" style="display:none;margin-top: 15px;"></div>
				</div>
				<div class="col-6" style="text-align:center;">
					<hr />
					<button class="btn btn-primary btn-sm start" id="on"><?php echo $_l->tpl('string.start'); ?></button>
					<button class="btn btn-primary btn-sm poke" id="restart"><?php echo $_l->tpl('string.restart'); ?></button>
					<button class="btn btn-danger btn-sm poke" id="off"><?php echo $_l->tpl('string.stop'); ?></button>
					<div id="pw_resp" style="display:none;margin-top: 15px;"></div>
				</div>
			</div>
		</div>
		<div class="modal fade" id="pauseConsole" tabindex="-1" role="dialog" aria-labelledby="PauseConsole" aria-hidden="true">
			<div class="modal-dialog">
				<div class="modal-content">
					<div class="modal-header">
						<button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
						<h4 class="modal-title" id="PauseConsole"><?php echo $_l->tpl('node.console.scrollstop_h1'); ?></h4>
					</div>
					<div class="modal-body">
						<div class="row">
							<div class="col-1"></div>
							<div class="col-10">
								<textarea id="paused_console" class="form-control console" readonly="readonly"></textarea>
							</div>
							<div class="col-1"></div>
						</div>
					</div>
					<div class="modal-footer">
						<button type="button" class="btn btn-default" data-dismiss="modal"><?php echo $_l->tpl('string.close'); ?></button>
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
		$('#live_console').scrollTop($('#live_console')[0].scrollHeight - $('#live_console').height());
		$("#console_command").submit(function(){
			$("#sending_command").html('<i class="fa fa-refresh fa-spin"></i>').addClass('disabled');
			var ccmd = $("#ccmd").val();
			$.ajax({
				type: "POST",
				headers: {"X-Access-Token": "<?php echo $core->server->nodeData('gsd_secret'); ?>"},
				url: 'http://<?php echo $core->server->nodeData('sftp_ip'); ?>:8003/gameservers/<?php echo $core->server->getData('gsd_id'); ?>/console',
				timeout: 5000,
				data: { command: ccmd },
				error: function(jqXHR, textStatus, errorThrown) {
				   	$("#sc_resp").html('<?php echo $_l->tpl('node.ajax.generic_error'); ?>').slideDown().delay(5000).slideUp();
				   	$("#sending_command").removeClass('disabled');
				   	$("#sending_command").html('&rarr;');
				   	$("#ccmd").val('');
				},
		  		success: function(data) {
		    		$("#sending_command").removeClass('disabled');
		    		$("#sending_command").html('&rarr;');
		    		$("#ccmd").val('');
		    			if(data !== true){
		    				$("#sc_resp").html(data).slideDown().delay(5000).slideUp();
		    			}
		 		 }
			});
			return false;
		});
		socket.on('console', function (data) {
			$("#live_console").val($("#live_console").val() + data.l);
			$('#live_console').scrollTop($('#live_console')[0].scrollHeight - $('#live_console').height());
		});
		$("#pause_console").click(function(){
			$("#paused_console").val();
			$("#paused_console").val($("#live_console").val());
		});
		var can_run = true;
		$(".poke").click(function(){
			var command = $(this).attr("id");
			if(command == 'off'){ uCommand = 'Stopping'; }else{ uCommand = 'Restarting';}
				if(can_run === true){
					can_run = false;
					$(this).append(' <i class="fa fa-refresh fa-spin"></i>');
					$(this).toggleClass('disabled');
					$.ajax({
						type: "GET",
						headers: {"X-Access-Token": "<?php echo $core->server->nodeData('gsd_secret'); ?>"},
						url: "http://<?php echo $core->server->nodeData('sftp_ip'); ?>:8003/gameservers/<?php echo $core->server->getData('gsd_id'); ?>/off",
						timeout: 5000,
						error: function(jqXHR, textStatus, errorThrown) {
						   	$("#pw_resp").attr('class', 'alert alert-danger').html('<?php echo $_l->tpl('node.ajax.generic_error'); ?>').slideDown().delay(5000).slideUp();
						   	$("#off").removeClass('disabled');
						   	$("#off").html('Stop');
						   	$("#restart").removeClass('disabled');
						   	$("#restart").html('Restart');
						   	can_run = true;
						   	return false;
						},
				  		success: function(data) {
			    			if(data == "ok"){
			    				$("#pw_resp").attr('class', 'alert alert-success').html("Server has been "+command+"ed successfully.").slideDown().delay(5000).slideUp();
			    				can_run = true;
			    				if(command == "restart"){
			    					start_server();
			    				}
			    				$("#off").removeClass('disabled');
			    				$("#off").html('<?php echo $_l->tpl('string.stop'); ?>');
			    				$("#restart").removeClass('disabled');
			    				$("#restart").html('<?php echo $_l->tpl('string.restart'); ?>');
			    				return false;
			    			}
				 		 }
					});
				}else{
					return false;
				}
		});
		$("#on").click(function(){
			start_server();
		});
		function start_server() {
			if(can_run === true){
				can_run = false;
				$("#on").append(' <i class="fa fa-refresh fa-spin"></i>');
				$("#on").toggleClass('disabled');
				$.ajax({
					type: "GET",
					url: "ajax/console/power.php",
					timeout: 5000,
					error: function(jqXHR, textStatus, errorThrown) {
					   	$("#pw_resp").attr('class', 'alert alert-danger').html('<?php echo $_l->tpl('node.ajax.generic_error'); ?>').slideDown().delay(5000).slideUp();
					   	$("#on").removeClass('disabled');
					   	$("#on").html('Start');
					   	can_run = true;
					   	return false;
					},
			  		success: function(data) {
		    			if(data == "ok"){
		    				$("#pw_resp").attr('class', 'alert alert-success').html("<?php echo $_l->tpl('node.console.ajax.server_started'); ?>").slideDown().delay(5000).slideUp();
		    				can_run = true;
		    			}else{
		    				$("#pw_resp").attr('class', 'alert alert-danger').html(data).slideDown().delay(5000).slideUp();
		    				can_run = true;
		    			}
		    			$("#on").toggleClass('disabled');
		    			$("#on").html('<?php echo $_l->tpl('string.start'); ?>');
			 		 }
				});
			}else{
				return false;
			}
		}
	});
	</script>
</body>
</html>