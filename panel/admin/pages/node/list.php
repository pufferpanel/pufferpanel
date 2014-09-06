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
require_once('../../../../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Page\components::redirect('../../../index.php?login');
}
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('../../../../src/include/header.php'); ?>
	<title>PufferPanel Admin Control Panel</title>
</head>
<body>
	<div class="container">
		<div class="navbar navbar-default">
			<div class="navbar-header">
				<a class="navbar-brand" href="#"><?php echo $core->settings->get('company_name'); ?></a>
			</div>
			<div class="navbar-collapse navbar-responsive-collapse collapse" style="height: 1px;">
				<ul class="nav navbar-nav navbar-right">
					<li class="dropdown">
						<a href="#" class="dropdown-toggle" data-toggle="dropdown">Account <b class="caret"></b></a>
							<ul class="dropdown-menu">
								<li><a href="<?php echo $core->settings->get('master_url'); ?>logout.php">Logout</a></li>
								<li><a href="<?php echo $core->settings->get('master_url'); ?>servers.php">View All Servers</a></li>
							</ul>
					</li>
				</ul>
			</div>
		</div>
		<div class="row">
			<div class="col-3"><?php include('../../../../src/include/admin.php'); ?></div>
			<div class="col-9">
				<table class="table table-striped table-bordered table-hover">
					<thead>
						<tr>
							<th>Node Name</th>
							<th>GSD</th>
						</tr>
					</thead>
					<tbody>
						<?php
						
							$find = $mysql->prepare("SELECT `id`, `node`, `node_ip` FROM `nodes`");
							$find->execute(array());
							
							while($row = $find->fetch())
								{
								
									echo '
									<tr class="dynUpdate" id="'.$row['node_ip'].'">
										<td><a href="view.php?id='.$row['id'].'">'.$row['node'].'</a></td>
										<td class="applyUpdate" style="width:5%;"><span class="label label-warning"> <i class="fa fa-refresh fa-spin"></i> </span></td>
									</tr>
									';
								
								}
						
						?>
					</tbody>
				</table>
			</div>
		</div>
		<div class="footer">
			<?php include('../../../../src/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$(document).ready(function(){
				$(".dynUpdate").each(function(index, data){
				    var connection = $(this).attr("id");
				    var element = $(this);
				    $.ajax({
				    	type: "POST",
				    	url: "ajax/get_status.php",
				    	data: { ip: connection },
				      		success: function(data) {
				    			element.find(".applyUpdate").html(data);
				     		}
				    });
				});
		});
	</script>
</body>
</html>
