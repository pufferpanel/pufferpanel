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
					<a href="#" class="list-group-item list-group-item-heading"><strong>Account Actions</strong></a>
					<a href="../account.php" class="list-group-item">Settings</a>
					<a href="../servers.php" class="list-group-item">My Servers</a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Server Actions</strong></a>
					<a href="index.php" class="list-group-item">Overview</a>
					<a href="console.php" class="list-group-item active">Live Console</a>
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
					<textarea id="live_console" class="form-control console" readonly="readonly"></textarea>
				</div>
				<div class="col-12">
					<div class="alert alert-danger text_highlighted" style="display:none;margin: 15px 0 -5px 0;">You have selected text in the console. The console will not auto-update when this occurs. This is done to allow you to easily copy or select text in the console. To allow for automatic refreshing again simply un-select the text.</div>
				</div>
				<div class="col-6">
					<hr />
					<form action="#" method="post" id="console_command">
						<fieldset>
							<div class="input-group">
								<input type="text" class="form-control" name="command" id="ccmd" placeholder="Enter Console Command" />
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
					<button class="btn btn-primary btn-sm start" id="on">Start</button>
					<button class="btn btn-primary btn-sm poke" id="restart">Restart</button>
					<button class="btn btn-danger btn-sm poke" id="off">Stop</button>
					<div id="pw_resp" style="display:none;margin-top: 15px;"></div>
				</div>
			</div>
		</div>
		<div class="modal fade" id="pauseConsole" tabindex="-1" role="dialog" aria-labelledby="PauseConsole" aria-hidden="true">
			<div class="modal-dialog">
				<div class="modal-content">
					<div class="modal-header">
						<button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
						<h4 class="modal-title" id="PauseConsole">ScrollStop&trade; Console Window</h4>
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
						<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
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
		var socket = io.connect('http://<?php echo $core->server->nodeData('sftp_ip'); ?>:<?php echo $core->server->getData('server_port') + 1; ?>');
		$.ajaxSetup({
		        error: function(jqXHR, exception) {
		            if (jqXHR.status === 0) {
		                alert('Not connect.\n Verify Network.');
		            } else if (jqXHR.status == 404) {
		                alert('Requested page not found. [404]');
		            } else if (jqXHR.status == 500) {
		                alert('Internal Server Error [500].');
		            } else if (exception === 'parsererror') {
		                alert('Requested JSON parse failed.');
		            } else if (exception === 'timeout') {
		                alert('Time out error.');
		            } else if (exception === 'abort') {
		                alert('Ajax request aborted.');
		            } else {
		                alert('Uncaught Error.\n' + jqXHR.responseText);
		            }
		        }
		    });
		$('#live_console').scrollTop($('#live_console')[0].scrollHeight - $('#live_console').height());
		$("#console_command").submit(function(){
			$("#sending_command").html('<i class="fa fa-refresh fa-spin"></i>').addClass('disabled');
			var ccmd = $("#ccmd").val();
			$.ajax({
				type: "POST",
				url: 'http://<?php echo $core->server->nodeData('sftp_ip'); ?>:8003/gameservers/<?php echo $core->server->getData('gsd_id'); ?>/console',
				timeout: 5000,
				data: { command: ccmd },
				error: function(jqXHR, textStatus, errorThrown) {
				   	$("#sc_resp").html('An error was encountered with this AJAX request.').slideDown().delay(5000).slideUp();
				   	$("#sending_command").removeClass('disabled');
				   	$("#sending_command").html('&rarr;');
				   	$("#ccmd").val('');
				},
		  		success: function(data) {
		    		$("#sending_command").removeClass('disabled');
		    		$("#sending_command").html('&rarr;');
		    		$("#ccmd").val('');
		    			if(data != 'ok'){
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
						url: "http://<?php echo $core->server->nodeData('sftp_ip'); ?>:8003/gameservers/<?php echo $core->server->getData('gsd_id'); ?>/off",
						timeout: 5000,
						error: function(jqXHR, textStatus, errorThrown) {
						   	$("#pw_resp").attr('class', 'alert alert-danger').html('An error was encountered with this AJAX request.').slideDown().delay(5000).slideUp();
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
			    				$("#off").html('Stop');
			    				$("#restart").removeClass('disabled');
			    				$("#restart").html('Restart');
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
					   	$("#pw_resp").attr('class', 'alert alert-danger').html('An error was encountered with this AJAX request.').slideDown().delay(5000).slideUp();
					   	$("#on").removeClass('disabled');
					   	$("#on").html('Start');
					   	can_run = true;
					   	return false;
					},
			  		success: function(data) {
		    			if(data == "ok"){
		    				$("#pw_resp").attr('class', 'alert alert-success').html("Server has been started successfully.").slideDown().delay(5000).slideUp();
		    				can_run = true;
		    			}else{
		    				$("#pw_resp").attr('class', 'alert alert-danger').html(data).slideDown().delay(5000).slideUp();
		    				can_run = true;
		    			}
		    			$("#on").toggleClass('disabled');
		    			$("#on").html('Start');
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