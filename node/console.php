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

/*
 * Are we on the correct node?
 */
//$url = parse_url($_SERVER["SERVER_NAME"], PHP_URL_PATH);
//$parts = explode('.', $url);
//
//	if($parts[0] != $core->framework->server->getData('node')){
//		$core->framework->page->redirect($core->framework->settings->get('master_url').'index.php');
//	}
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title><?php echo $core->framework->settings->get('company_name'); ?> - Server Console</title>
	
	<!-- Stylesheets -->
	<link href='http://fonts.googleapis.com/css?family=Droid+Sans:400,700' rel='stylesheet'>
	<link rel="stylesheet" href="<?php echo $core->framework->settings->get('master_url'); ?>assets/css/style.css">
	
	<!-- Optimize for mobile devices -->
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
	
	<!-- jQuery & JS files -->
	<script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
	
	<script type="text/javascript" src="<?php echo $core->framework->settings->get('master_url'); ?>assets/javascript/jquery.ba-throttle-debounce.min.js"></script>
</head>
<body>
	<div id="top-bar">
		<div class="page-full-width cf">
			<ul id="nav" class="fl">
				<li><a href="#" class="round button dark"><i class="fa fa-user"></i>&nbsp;&nbsp; <strong><?php echo $core->framework->user->getData('username'); ?></strong></a></li>
				<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>servers.php" class="round button dark"><i class="fa fa-home"></i></a></li>
				<li><a class="round button dark"><i class="fa fa-hdd"></i>&nbsp;&nbsp; <?php echo $core->framework->server->getData('name'); ?></a></li>
			</ul>
			<ul id="nav" class="fr">
				<?php if($core->framework->user->getData('root_admin') == 1){ echo '<li><a href="'.$core->framework->settings->get('master_url').'admin/index.php" class="round button dark"><i class="fa fa-bar-chart-o"></i>&nbsp;&nbsp; Admin CP</a></li>'; } ?>
				<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>logout.php" class="round button dark"><i class="fa fa-power-off"></i></a></li>
			</ul>
		</div>	
	</div>
	<div id="header-with-tabs">
		<div class="page-full-width cf">
		</div>
	</div>
	<div id="content">
		<div class="page-full-width cf">
			<div class="side-menu fl">
				<h3>Account Actions</h3>
				<ul>
					<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>account.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Edit Settings</a></li>
					<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>servers.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> My Servers</a></li>
				</ul>
				<h3>Server Actions</h3>
				<ul>
					<li><a href="index.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Overview</a></li>
					<li><a href="console.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Live Console</a></li>
					<li><a href="settings.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Server Settings</a></li>
					<li><a href="plugins.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Server Plugins</a></li>
					<li><a href="files/index.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> File Manager</a></li>
					<li><a href="backup.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Backup Manager</a></li>
				</ul>
			</div>
			<div class="side-content fr">
				<div class="error-box round text_highlighted" style="display:none;">You have selected text in the console. The console will not auto-update when this occurs. This is done to allow you to easily copy or select text in the console. To allow for automatic refreshing again simply un-select the text.</div>
				<div class="content-module">
					<div class="content-module-heading cf">
						<h3 class="fl">Live Console</h3>
					</div> <!-- end content-module-heading -->
					<div class="content-module-main" id="server_info">
						<textarea id="live_console" disabled="disabled"><?php echo $core->framework->files->last_lines($core->framework->server->nodeData('server_dir').$core->framework->server->getData('ftp_user').'/server/logs/latest.log', 250); ?></textarea>
					</div> <!-- end content-module-main -->
				</div>
                <div class="error-box round text_highlighted" style="display:none;">You have selected text in the console. The console will not auto-update when this occurs. This is done to allow you to easily copy or select text in the console. To allow for automatic refreshing again simply un-select the text.</div>
			</div>
			<div class="side-content fr">
				<div class="half-size-column fl">
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Send Command</h3>
						</div>
						<div class="content-module-main" id="server_players">
                            <div class="error-box round" id="sc_resp" style="display:none;"></div>
							<form action="#" method="post" id="console_command">
								<fieldset>
									<p>
										<input type="text" name="command" id="ccmd" placeholder="Enter Command to Send" class="round full-width-input" />
									</p>
									<div class="stripe-separator"><!--  --></div>
									<input type="submit" value="Send Command" id="sending_command" class="round blue ic-right-arrow" />
								</fieldset>
							</form>	
						</div>
					</div>
				</div>
				<div class="half-size-column fr">
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Server Actions</h3>
						</div>
						<div class="content-module-main cf center" id="server_stats">
                            <div class="confirmation-box round" id="pw_resp" style="display:none;"></div>
							<a href="#start" id="start" class="poke round button blue text-upper">Start Server</a>
							<a href="#stop" id="stop" class="poke round button blue text-upper">Stop Server</a>
							<a href="#kill" id="kill" class="poke round button blue text-upper">Kill Server</a>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.4.2 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
	<script type="text/javascript">
			$(document).ready(function(){
				
				$('#live_console').scrollTop($('#live_console')[0].scrollHeight - $('#live_console').height());
				
				$("#console_command").submit(function(){
					$("#sending_command").css({'background' : '#2069B4 url(\'<?php echo $core->framework->settings->get('master_url'); ?>assets/images/load/small_blue.gif\') 5px center no-repeat', 'padding-left' : '25px', 'padding-right' : '0.833em'});
					$("#sending_command").attr('value', 'Sending...');
					$("#sending_command").fadeTo(500, 0.5);
					var ccmd = $("#ccmd").val();
					$.ajax({
						type: "POST",
						url: 'core/ajax/console/send.php',
						data: { command: ccmd },
					  		success: function(data) {
					    		$("#sending_command").css({'background' : '#2069B4 url("<?php echo $core->framework->settings->get('master_url'); ?>assets/images/icons/ic_right.png") right center no-repeat ', 'padding-left' : '0.833em', 'padding-right' : '3em'});
					    		$("#sending_command").attr('value', 'Send Command');
					    		$("#sending_command").fadeTo(500, 1);
					    		$("#ccmd").val('');
					    		
					    			if(data != ''){
					    			
					    				$("#sc_resp").html(data).slideDown().delay(5000).slideUp();
					    			
					    			}
					    		
					    		updateConsole();
					    		
					 		 }
					});
					return false;
				});
				
				/*
				 * Autoload Console
				 */
				function updateConsole() {
					
					if(isScroll !== true){
						
						$.ajax({
							type: "GET",
							url: 'core/ajax/console/update.php',
						  		success: function(data) {
						    		if(isTextSelected($('#live_console')[0]) === false){
							    		
										if(isBottom() === true){
											var b = true;
											var bottom = $("#live_console")[0].scrollHeight - $("#live_console").height();
										}else{
											b = false;
											curloc = $('#live_console').scrollTop();
										}
	
										$("#live_console").html(data);
							    		
										if(b === true){
							    			$('#live_console').scrollTop(bottom);
											bottom = 0;
							    		}else{
							    			$('#live_console').scrollTop(curloc);
											curloc = 0;
							    		}
							    	}else{ /*Do Nothing*/ }
						 		 }
						});
						
					}
					
				}
				
				var isScroll;
				$("#live_console").scroll($.debounce(100, true, function(){ isScroll = true;}));
				$("#live_console").scroll($.debounce(100, function(){ isScroll = false;}));
				
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
				}, 200);
				
				setInterval(function(){
					updateConsole();
				}, 1000);
				
			});
			
			var can_run = true;
			$(".poke").click(function(){
				var command = $(this).attr("id");
				if(command == 'stop'){ uCommand = 'Stopping...'; }else if(command == 'start'){ uCommand = 'Starting...';}else{ uCommand = 'Killing...';}
					
					if(can_run === true){
					
						$(this).css({'background' : '#2069B4 url(\'<?php echo $core->framework->settings->get('master_url'); ?>assets/images/load/small_blue.gif\') 5px center no-repeat', 'padding-left' : '25px'});
						$(this).html(uCommand);
						$("a.poke").fadeTo(1000, 0.5);
					
						$.ajax({
							type: "POST",
							url: "core/ajax/console/power.php",
							data: { process: "power", command: command },
						  		success: function(data) {
					    			if(data == "Server Started."){
					    				
					    				$("#"+command).css({'background' : '#2069B4', 'padding-left' : '0.833em'});
					    				$("#"+command).html('Start Server');
					    				$("a.poke").fadeTo(1000, 1);
					    				
					    				//make nicer
					    				$("#pw_resp").html("Server has been started successfully.").slideDown().delay(5000).slideUp();
					    				
					    				return false;
					    			
					    			}else if(data == "Server Stopped."){
					    				
					    				$("#"+command).css({'background' : '#2069B4', 'padding-left' : '0.833em'});
					    				$("#"+command).html('Stop Server');
					    				$("a.poke").fadeTo(1000, 1);
					    				
					    				//make nicer
					    				$("#pw_resp").html("Server has been stopped successfully.").slideDown().delay(5000).slideUp();
					    				
					    				return false;
					    			
					    			}else if(data == "Server Killed."){
					    			
					    				$("#"+command).css({'background' : '#2069B4', 'padding-left' : '0.833em'});
					    				$("#"+command).html('Kill Server');
					    				$("a.poke").fadeTo(1000, 1);
					    				
					    				//make nicer
					    				$("#pw_resp").html("The server java process has been killed. Please check your data for possible corruption.").slideDown().delay(5000).slideUp();
					    				
					    				return false;
					    					
					    			}else{
					    				
					    				//make nicer
					    				alert(data);
					    				
					    				$("#stop").css({'background' : '#2069B4', 'padding-left' : '0.833em'});
					    				$("#stop").html('Stop Server');
					    				$("#start").css({'background' : '#2069B4', 'padding-left' : '0.833em'});
					    				$("#start").html('Start Server');
					    				$("#kill").css({'background' : '#2069B4', 'padding-left' : '0.833em'});
					    				$("#kill").html('Kill Server');
					    				$("a.poke").fadeTo(1000, 1);
					    				
					    			}
						 		 }
						});
						
					}else{
					
					
					}
					
					
				return false;
				
			});
		</script>
</body>
</html>