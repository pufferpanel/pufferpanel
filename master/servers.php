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
$error = '';

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token')) !== true){
	$core->framework->page->redirect('index.php', $core->framework->page->genRedirect());
	exit();
}

/*
 * Redirect
 */
if(isset($_GET['goto']) && !empty($_GET['goto'])){

	$core->framework->page->nodeRedirect($_GET['goto']);
	
}

if($core->framework->user->getData('root_admin') == '1'){
	$query = $mysql->prepare("SELECT * FROM `servers` ORDER BY `active` DESC");
	$query->execute();
}else{
	$query = $mysql->prepare("SELECT * FROM `servers` WHERE `owner_id` = :oid ORDER BY `active` DESC");
	$query->execute(array(':oid' => $core->framework->user->getData('id')));
}

/*
 * List Servers
 */
$listServers = '';
while($row = $query->fetch()){
	
	($row['active'] == '1') ? $isActive = 'Enabled' : $isActive = 'Disabled';
	$listServers .= '
					<tr>
						<td><a href="servers.php?goto='.$row['hash'].'">'.$row['name'].'</a></td>
						<td>'.$row['node'].'</td>
						<td>'.$row['server_ip'].'</td>
						<td>'.$row['server_port'].'</td>
						<td>'.$row['max_ram'].' MB</td>
						<td>'.$row['disk_space'].' MB</td>
						<td>'.$isActive.'</td>
					</tr>
					';

}

?>
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>PufferPanel - My Servers</title>
	
	<!-- Stylesheets -->
	<link href='http://fonts.googleapis.com/css?family=Droid+Sans:400,700' rel='stylesheet'>
	<link rel="stylesheet" href="assets/css/style.css">
	
	<!-- Optimize for mobile devices -->
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
	
	<!-- jQuery & JS files -->
	<script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
	
</head>
<body>
	<div id="top-bar">
		<div class="page-full-width cf">
			<ul id="nav" class="fl">
				<li><a href="#" class="round button dark"><i class="fa fa-user"></i>&nbsp;&nbsp; <strong><?php echo $core->framework->user->getData('username'); ?></strong></a></li>
			</ul>
			<ul id="nav" class="fr">
				<?php if($core->framework->user->getData('root_admin') == 1){ echo '<li><a href="admin/index.php" class="round button dark"><i class="fa fa-bar-chart-o"></i>&nbsp;&nbsp; Admin CP</a></li>'; } ?>
				<li><a href="logout.php" class="round button dark"><i class="fa fa-power-off"></i></a></li>
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
					<li><a href="account.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Edit Settings</a></li>
					<li><a href="servers.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> My Servers</a></li>
				</ul>
				<h3>Server Actions</h3>
			</div>
			<div class="side-content fr">
				<?php if(isset($_GET['error'])){ echo '<div class="error-box round">Unable to locate that specific server. Has it been suspended?</div>'; } ?>
				<div class="content-module">
					<div class="content-module-heading cf">
						<h3 class="fl">My Servers</h3>
					</div>
					<div class="content-module-main">
						<table>
							<thead>
								<tr>
									<th style="width:15%">Server Name</th>
									<th style="width:10%">Node</th>
									<th style="width:20%">IP Address</th>
									<th style="width:10%">Port</th>
									<th style="width:15%">Memory</th>
									<th style="width:15%">Disk Space</th>
									<th style="width:15%">Status</th>
								</tr>
							</thead>
							<tbody>
								<?php echo $listServers; ?>
							</tbody>
						</table>
					</div>
				</div>
			</div>
		</div>
	</div>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.4.2 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>