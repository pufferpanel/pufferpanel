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

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token')) !== true){
	Page\components::redirect('index.php?login');
	exit();
}

/*
 * Redirect
 */
if(isset($_GET['goto']) && !empty($_GET['goto'])){

	$core->server->nodeRedirect($_GET['goto']);
	
}

if($core->user->getData('root_admin') == '1'){
	$query = $mysql->prepare("SELECT * FROM `servers` ORDER BY `active` DESC");
	$query->execute();
}else{
	$query = $mysql->prepare("SELECT * FROM `servers` WHERE `owner_id` = :oid ORDER BY `active` DESC");
	$query->execute(array(':oid' => $core->user->getData('id')));
}

/*
 * List Servers
 */
$listServers = '';
while($row = $query->fetch()){
	
	$listServers .= '
					<tr class="dynUpdate" id="'.$row['id'].'">
						<td>'.$core->settings->nodeName($row['node']).'</td>
						<td><a href="servers.php?goto='.$row['hash'].'">'.$row['name'].'</a></td>
						<td>'.$row['server_ip'].':'.$row['server_port'].'</td>
						<td class="applyUpdate" style="width:5%;"><span class="label label-warning"> <i class="fa fa-refresh fa-spin"></i> </span></td>
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
		<?php include('assets/include/navbar.php'); ?>
		<div class="row">
			<div class="col-3">
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.acc_actions'); ?></strong></a>
					<a href="account.php" class="list-group-item"><?php echo $_l->tpl('sidebar.settings'); ?></a>
					<a href="servers.php" class="list-group-item active"><?php echo $_l->tpl('sidebar.servers'); ?></a>
				</div>
			</div>
			<div class="col-9">
				<table class="table table-striped table-bordered table-hover">
					<thead>
						<tr>
							<th><?php echo $_l->tpl('string.node'); ?></th>
							<th><?php echo $_l->tpl('string.name'); ?></th>
							<th><?php echo $_l->tpl('string.connect'); ?></th>
							<th style="width:5%;"><?php echo $_l->tpl('string.status'); ?></th>
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
	<script type="text/javascript">
		$(document).ready(function(){
				$(".dynUpdate").each(function(index, data){
				    var connection = $(this).attr("id");
				    var element = $(this);
				    $.ajax({
				    	type: "POST",
				    	url: "core/ajax/get_status.php",
				    	data: { server: connection },
				      		success: function(data) {
				    			element.find(".applyUpdate").html(data);
				     		}
				    });
				});
		});
	</script>
</body>
</html>