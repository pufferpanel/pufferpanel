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

$filesIncluded = true;

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
					<a href="console.php" class="list-group-item active">Live Console</a>
					<a href="settings.php" class="list-group-item">Server Settings</a>
					<a href="plugins.php" class="list-group-item">Server Plugins</a>
					<a href="files/index.php" class="list-group-item">File Manager</a>
					<a href="backup.php" class="list-group-item">Backup Manager</a>
				</div>
			</div>
			<div class="col-9">
				<div class="col-12">
					<textarea id="live_console" class="form-control console" readonly="readonly"><?php echo $core->framework->files->last_lines($core->framework->server->nodeData('server_dir').$core->framework->server->getData('ftp_user').'/server/logs/latest.log', 250); ?></textarea>
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
								</span>
							</div>
						</fieldset>
					</form>
					<div class="alert alert-danger" id="sc_resp" style="display:none;margin-top: 15px;"></div>
				</div>
				<div class="col-6" style="text-align:center;">
					<hr />
					<button class="btn btn-primary btn-sm poke" id="start">Start</button>
					<button class="btn btn-primary btn-sm poke" id="stop">Stop</button>
					<button class="btn btn-danger btn-sm poke" id="kill">Kill</button>
					<div class="alert alert-info" id="pw_resp" style="display:none;margin-top: 15px;"></div>
				</div>
			</div>
		</div>
		<div class="footer">
			<?php include('assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
	$(document).ready(function(){
		$('#live_console').scrollTop($('#live_console')[0].scrollHeight - $('#live_console').height());
		$("#console_command").submit(function(){
			$("#sending_command").html('<i class="fa fa-refresh fa-spin"></i>').addClass('disabled');
			var ccmd = $("#ccmd").val();
			$.ajax({
				type: "POST",
				url: 'core/ajax/console/send.php',
				data: { command: ccmd },
			  		success: function(data) {
			    		$("#sending_command").removeClass('disabled');
			    		$("#sending_command").html('&rarr;');
			    		$("#ccmd").val('');
			    			if(data != ''){
			    				$("#sc_resp").html(data).slideDown().delay(5000).slideUp();
			    			}
			    		updateConsole();
			 		 }
			});
			return false;
		});
		var isScroll;
		$("#live_console").scroll($.debounce(100, true, function(){ isScroll = true;}));
		$("#live_console").scroll($.debounce(100, function(){ isScroll = false;}));
		function updateConsole() {
			var b = true;
			var curloc = 0;
			if(isScroll !== true){
				$.ajax({
					type: "GET",
					url: 'core/ajax/console/update.php',
				  		success: function(data) {
				    		if(isTextSelected($('#live_console')[0]) === false){
								if(isBottom() !== true){
									b = false;
									curloc = $('#live_console').scrollTop();
								}
								$("#live_console").html(data);
								if(b === true){
					    			$('#live_console').scrollTop($('#live_console')[0].scrollHeight - $('#live_console').height());
					    		}else{
					    			$('#live_console').scrollTop(curloc);
					    		}
					    	}else{ /*Do Nothing*/ }
				 		 }
				});
			}
		}		
		function isTextSelected(input){
		   var startPos = input.selectionStart;
		   var endPos = input.selectionEnd;
		   var doc = document.selection;
		   if(doc && doc.createRange().text.length != 0){
		      return true;
		   }else if (!doc && input.value.substring(startPos,endPos).length != 0){
		      return true;
		   }
		   return false;
		}
		function isBottom() {
			if(($('#live_console').scrollTop() + $('#live_console').innerHeight()) >= $('#live_console')[0].scrollHeight){
				return true;
			}
		}
		setInterval(function(){
			if(isTextSelected($('#live_console')[0]) === false){
				$(".text_highlighted").slideUp();
				updateConsole();
			}else{
				$(".text_highlighted").slideDown();
			}
			updateConsole();
		}, 1000);
		var can_run = true;
		$(".poke").click(function(){
			var command = $(this).attr("id");
			if(command == 'stop'){ uCommand = 'Stopping'; }else if(command == 'start'){ uCommand = 'Starting';}else{ uCommand = 'Killing';}
				if(can_run === true){
					can_run = false;
					$(this).append(' <i class="fa fa-refresh fa-spin"></i>');
					$(this).toggleClass('disabled');
					$.ajax({
						type: "POST",
						url: "core/ajax/console/power.php",
						data: { process: "power", command: command },
					  		success: function(data) {
				    			if(data == "Server Started."){
				    				$("#"+command).toggleClass('disabled');
				    				$("#"+command).html('Start');
				    				$("#pw_resp").html("Server has been started successfully.").slideDown().delay(5000).slideUp();
				    				can_run = true;
				    				return false;
				    			}else if(data == "Server Stopped."){
				    				$("#"+command).toggleClass('disabled');
				    				$("#"+command).html('Stop');
				    				$("#pw_resp").html("Server has been stopped successfully.").slideDown().delay(5000).slideUp();
				    				can_run = true;
				    				return false;
				    			}else if(data == "Server Killed."){
				    				$("#"+command).toggleClass('disabled');
				    				$("#"+command).html('Kill');
				    				$("#pw_resp").html("The server java process has been killed. Please check your data for possible corruption.").slideDown().delay(5000).slideUp();
				    				can_run = true;
				    				return false;
				    			}else{
				    				$("#pw_resp").html(data);				    				
				    				$("#stop").removeClass('disabled');
				    				$("#stop").html('Stop');
				    				$("#start").removeClass('disabled');
				    				$("#start").html('Start');
				    				$("#kill").removeClass('disabled');
				    				$("#kill").html('Kill');
				    				can_run = true;
				    			}
					 		 }
					});
				}else{
					return false;
				}
		});
	});
	</script>
</body>
</html>