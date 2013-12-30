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
						<td>'.$core->framework->settings->nodeName($row['node']).'</td>
						<td><a href="servers.php?goto='.$row['hash'].'">'.$row['name'].'</a></td>
						<td>'.$row['server_ip'].':'.$row['server_port'].'</td>
						<td>'.$isActive.'</td>
					</tr>
					';

}

?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('assets/include/header.php'); ?>
	<title>PufferPanel - Your Servers</title>
</head>
<body>
	<div class="container">
		<div class="navbar navbar-default">
			<div class="navbar-header">
				<a class="navbar-brand" href="#"><?php echo $core->framework->settings->get('company_name'); ?></a>
			</div>
			<div class="navbar-collapse navbar-responsive-collapse collapse" style="height: 1px;">
				<ul class="nav navbar-nav navbar-right">
					<li class="dropdown">
						<a href="#" class="dropdown-toggle" data-toggle="dropdown">Account <b class="caret"></b></a>
							<ul class="dropdown-menu">
								<li><a href="logout.php">Logout</a></li>
								<?php if($core->framework->user->getData('root_admin') == 1){ echo '<li><a href="admin/index.php">Admin CP</a></li>'; } ?>
							</ul>
					</li>
				</ul>
			</div>
		</div>
		<div class="row">
			<div class="col-3">
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Account Actions</strong></a>
					<a href="account.php" class="list-group-item">Settings</a>
					<a href="servers.php" class="list-group-item active">My Servers</a>
				</div>
			</div>
			<div class="col-9">
				<table class="table table-striped table-bordered table-hover">
					<thead>
						<tr>
							<th>Node</th>
							<th>Name</th>
							<th>Connect</th>
							<th>Status</th>
						</tr>
					</thead>
					<tbody>
						<?php echo $listServers; ?>
					</tbody>
				</table>
			</div>
		</div>
		<div class="footer">
			<?php include('assets/include/footer.php'); ?>
		</div>
	</div>
</body>
</html>