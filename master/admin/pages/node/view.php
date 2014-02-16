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
require_once('../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), null, true) !== true){
	$core->framework->page->redirect('../../../index.php');
}

if(isset($_GET['do']) && $_GET['do'] == 'redirect' && isset($_GET['node'])){

	$select = $mysql->prepare("SELECT `id` FROM `nodes` WHERE `node` = :name");
	$select->execute(array(':name' => $_GET['node']));
	
	if($select->rowCount() == 1) {
		$n = $select->fetch();
		$core->framework->page->redirect('view.php?id='.$n['id']);
	}else{
		$core->framework->page->redirect('list.php');
	}

}

if(!isset($_GET['id']))
	$core->framework->page->redirect('list.php');

/*
 * Select Node Information
 */
$selectNode = $mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :id");
$selectNode->execute(array(
	':id' => $_GET['id']
));

	if($selectNode->rowCount() != 1)
		$core->framework->page->redirect('list.php?error=no_node');
	else
		$node = $selectNode->fetch();

?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('../../../assets/include/header.php'); ?>
	<title>PufferPanel Admin Control Panel</title>
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
								<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>logout.php">Logout</a></li>
								<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>servers.php">View All Servers</a></li>
							</ul>
					</li>
				</ul>
			</div>
		</div>
		<div class="row">
			<div class="col-3"><?php include('../../../assets/include/admin.php'); ?></div>
			<div class="col-9">
				<ul class="nav nav-tabs" id="config_tabs">
					<li class="active"><a href="#info" data-toggle="tab">Information</a></li>
					<li><a href="#allocation" data-toggle="tab">Allocation</a></li>
					<li><a href="#sftp" data-toggle="tab">SFTP</a></li>
				</ul>
				<?php 
					
					if(isset($_GET['disp']) && !empty($_GET['disp'])){
					
						echo '<div class="alert alert-danger" style="margin-top:15px;">';
						switch($_GET['disp']){
							
							case 'missing_warn':
								echo 'You must agree to the warning before updating the information.';
								break;
							case 'missing_args':
								echo 'Not all arguments were passed by the script.';
								break;
							case 'ip_fail':
								echo 'The IP address provided for SFTP was invalid.';
								break;
							case 'user_fail':
								echo 'SFTP users must not be blank, and may not be \'root\'.';
								break;
							case 'pass_fail':
								echo 'SSH passwords must be at least 12 characters.';
								break;
							case 'n_fail':
								echo 'The node name does not meet the requirements (1-15 characters, a-zA-Z0-9_.-).';
								break;
							case 'url_fail':
								echo 'The node URL provided is not valid.';
								break;
							case 'add_port_fail':
								echo 'The port list entered was invalid.';
								break;
						
						}
						echo '</div>';
					
					}
				
				?>
				<div class="tab-content">
					<div class="tab-pane active" id="info">
						<h3>Basic Information</h3><hr />
						<form action="ajax/update/basic.php" method="post">
							<fieldset>
								<div class="form-group">
									<label for="name" class="control-label">Node Name</label>
									<div>
										<input type="text" name="name" value="<?php echo $node['node']; ?>" class="form-control" />
									</div>
								</div>
								<div class="form-group">
									<label for="ip" class="control-label">Node IP</label>
									<div>
										<input type="text" name="ip" value="<?php echo $node['node_ip']; ?>" class="form-control" />
									</div>
								</div>
								<div class="form-group">
									<div>
										<input type="hidden" name="nid" value="<?php echo $_GET['id']; ?>" />
										<input type="submit" value="Update Information" class="btn btn-primary" />
									</div>
								</div>
							</fieldset>
						</form>
					</div>
					<div class="tab-pane" id="allocation">
						<h3>IP &amp; Port Allocation</h3><hr />
						<form action="ajax/ports/add_port.php" id="addPorts" style="display: none;" method="post">
							<div class="form-group">
								<label for="add_ports" class="control-label" id="setTitle"></label>
								<div class="input-group">
									<input type="text" name="add_ports" value="" placeholder="enter a comma separated list of ports to add; enter to submit" class="form-control" />
									<span class="input-group-btn">
										<input type="hidden" name="add_ports_ip" value=""/>
										<input type="hidden" name="add_ports_node" value=""/>
										<button class="btn btn-primary" type="submit">&rarr;</button>
									</span>
								</div>
							</div>
						</form>
						<table class="table table-striped table-bordered table-hover">
							<thead>
								<tr>
									<th>IP Address</th>
									<th>Ports</th>
									<th></th>
								</tr>
							</thead>
							<tbody>
								<?php
								
									foreach(json_decode($node['ports'], true) as $ip => $ports)
										{
								
											echo "<tr><td style=\"vertical-align:top;\">{$ip}<br /><a href=\"#/add/{$ip}/{$node['id']}\" class=\"clickToAdd\" onclick=\"return false;\">Add Port(s)</a></td>";
											$counter = 1;
											$row1 = null; $row2 = null;
											foreach($ports as $port => $avaliable)
												{
										
													if($counter & 1)
														{
														
															$row1 .= ($avaliable == 1) ? "<span><a href=\"#/delete/{$ip}/{$port}/{$node['id']}\" class=\"deletePort\" onclick=\"return false;\"><i class=\"fa fa-circle-o\"></i></a>" : "<i class=\"fa fa-dot-circle-o\"></i>";
															$row1 .= "&nbsp;&nbsp;&nbsp; {$port}<br /></span>";
															
														}else{
														
															$row2 .= ($avaliable == 1) ? "<span><a href=\"#/delete/{$ip}/{$port}/{$node['id']}\" class=\"deletePort\" onclick=\"return false;\"><i class=\"fa fa-circle-o\"></i></a>" : "<i class=\"fa fa-dot-circle-o\"></i>";
															$row2 .= "&nbsp;&nbsp;&nbsp; {$port}<br /></span>";
														
														}
													
													$counter++;
									
												}
											echo "<td style=\"vertical-align:top;\">{$row1}</td><td style=\"vertical-align:top;\">{$row2}</td></tr>";
								
										}
								
								?>
							</tbody>
						</table>
						<div class="panel panel-default">
							<div class="panel-heading">Key</div>
							<div class="panel-body">
								<p><i class="fa fa-circle-o"></i> (Port Available; Click to Delete Port)</p><p><i class="fa fa-dot-circle-o"></i> (Port Used; Cannot Delete)</p>
							</div>
						</div>
					</div>
					<div class="tab-pane" id="sftp">
						<h3>SFTP Settings</h3><hr />
						<div class="well">
							<form action="ajax/update/sftp.php?do=ipuser" method="post">
								<fieldset>
									<div class="form-group">
										<label for="sftp_ip" class="control-label">SFTP IP</label>
										<div>
											<input type="text" name="sftp_ip" value="<?php echo $node['sftp_ip']; ?>" class="form-control" />
										</div>
									</div>
									<div class="form-group">
										<label for="sftp_user" class="control-label">SFTP Username</label>
										<div>
											<input type="text" name="sftp_user" value="<?php echo $node['username']; ?>" class="form-control" />
										</div>
									</div>
									<div class="form-group">
										<div>
											<div class="alert alert-warning">Editing your username will require that you also update the account password below.</div>
										</div>
										<div class="checkbox">
											<label>
												<input type="checkbox" id="warning_1" name="warning" /> I have read and understand the above statement.
											</label>
										</div>
									</div>
									<div class="form-group">
										<div>
											<input type="hidden" name="nid" value="<?php echo $_GET['id']; ?>" />
											<input type="submit" value="Update SFTP Information" id="disable_complete" class="btn btn-primary disabled" />
										</div>
									</div>
								</fieldset>
							</form>
						</div>
						<div class="well">
							<form action="ajax/update/sftp.php?do=pass" method="post">
								<fieldset>
									<div class="form-group">
										<label for="sftp_ip" class="control-label">New SFTP Password</label>
										<div>
											<input type="password" name="pass" autocomplete="off" class="form-control" />
										</div>
									</div>
									<div class="form-group">
										<div>
											<div class="alert alert-warning">Please ensure that you have entered the above information correctly. Changing this wrongly could result in multiple clients being unable to access their server(s).</div>
										</div>
										<div class="checkbox">
											<label>
												<input type="checkbox" id="warning_2" name="warning" /> I have read and understand the above statement.
											</label>
										</div>
									</div>
									<div class="form-group">
										<div>
											<input type="hidden" name="nid" value="<?php echo $_GET['id']; ?>" />
											<input type="submit" value="Update SFTP Password" id="disable_complete_pass" class="btn btn-primary disabled" />
										</div>
									</div>
								</fieldset>
							</form>
						</div>
					</div>
				</div>
			</div>
		</div>
		<div class="footer">
			<?php include('../../../assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$(document).ready(function(){
			setActiveOption('node-list');
			$(".clickToAdd").click(function(){
				var rawUrl = $(this).attr("href");
				var exploded = rawUrl.split('/');
				var ip = exploded[2];
				var node_id = exploded[3];
				$("#addPorts").slideUp(function(){
					$("#setTitle").html('Add Ports for '+ip);
					$("input[name='add_ports']").val("");
					$("input[name='add_ports_ip']").val(ip);
					$("input[name='add_ports_node']").val(node_id);
					$("#addPorts").slideDown();
				});
			});
			$(".deletePort").click(function(){
				
				var rawUrl = $(this).attr("href");
				var exploded = rawUrl.split('/');
				var ip = exploded[2];
				var port = exploded[3];
				var node_id = exploded[4];
				var conf = confirm("Are you sure you want to delete "+ip+":"+port);
				
					if(conf == true)
						{
							$.ajax({
								type: "POST",
								url: "ajax/ports/delete.php",
								data: { ip: ip, port: port, node: node_id},
								success: function(data) {
									$(".deletePort[href='#/delete/"+ip+"/"+port+"/"+node_id+"']").parent().fadeOut();
								}
							});
						}else{
							return false;
						}
				
			});
			if($.urlParam('error') != null){
				var field = $.urlParam('error');
				var exploded = field.split('|');
					$.each(exploded, function(key, value) {
						$('[name="' + value + '"]').parent().parent().addClass('has-error');
					});
			}
			$("#warning_1").click(function(){
				if($(this).is(":checked"))
					$("#disable_complete").removeClass("disabled");
				else
					$("#disable_complete").addClass("disabled");
			});
			$("#warning_2").click(function(){
				if($(this).is(":checked"))
					$("#disable_complete_pass").removeClass("disabled");
				else
					$("#disable_complete_pass").addClass("disabled");
			});
			if($.urlParam('tab') != null){
				$('#config_tabs a[href="#'+$.urlParam('tab')+'"]').tab('show');
			}
		});
	</script>
</body>
</html>