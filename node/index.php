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

/*
 * Are we on the correct node?
 * NOTE: MAKE SURE NODES ARE ON A SUBDOMAIN STARTING WITH THE NODE NAME! node1.example.com. example.com/node1 will not work.
 */

//Look for more graceful method.
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
	<title><?php echo $core->framework->settings->get('company_name'); ?> - Server Dashboard</title>
	
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
				<div class="half-size-column fl">
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Players Online</h3>
						</div>
						<div class="content-module-main" id="server_players">
							<p id="server_players_loading" style="margin: 1.25em;text-align: center;color: #000;"><i class="fa fa-cog fa-3x fa-spin"></i></p>
						</div>
					</div>
				</div>
				<div class="half-size-column fr">
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Server Statistics</h3>
						</div>
						<div class="content-module-main cf" id="server_stats">
							<p id="server_stats_loading" style="margin: 1.25em;text-align: center;color: #000;"><i class="fa fa-cog fa-3x fa-spin"></i></p>
						</div>
					</div>
				</div>
			</div>
			<div class="side-content fr">
				<div class="content-module">
					<div class="content-module-heading cf">
						<h3 class="fl">Server Information</h3>
					</div> <!-- end content-module-heading -->
					<div class="content-module-main" id="server_info">
						<p id="server_info_loading" style="margin: 1.25em;text-align: center;color: #000;"><i class="fa fa-cog fa-3x fa-spin"></i></p>
					</div> <!-- end content-module-main -->
				</div>
			</div>
		</div>
	</div>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.4.2 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
	<script type="text/javascript">
	$(document).ready(function(){
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
		$.ajax({
			type: "POST",
			url: "core/ajax/overview/data.php",
			data: { command: 'players' },
		  		success: function(data) {
					$("#server_players_loading").slideUp("slow", function(){
						$("#server_players").hide();
						$("#server_players").html(data);
						$("#server_players").slideDown("slow");				
					});
		 		}
		});
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
	});
	</script>
</body>
</html>