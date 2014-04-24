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

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Page\components::redirect('../../../index.php');
}

$packs = $mysql->prepare("SELECT * FROM `modpacks` WHERE `hash` = :mid");
$packs->execute(array(':mid' => $_GET['mid']));

if($packs->rowCount() == 1)
	$pack = $packs->fetch();
else
	Page\components::redirect('modpacks.php');

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
			<div class="col-3"><?php include('../../../assets/include/admin.php'); ?></div>
			<div class="col-9">
				<h3>Modpack: <strong><?php echo $pack['name'].' ('.$pack['version'].')'; ?></strong></h3>
				<h5>Modpack File Name: <strong><?php echo $pack['hash']; ?>.zip</strong></h5><hr />
				<?php 
					
					if(isset($_GET['disp']) && !empty($_GET['disp'])){
					
						switch($_GET['disp']){
						
							case 'missing_args':
								echo '<div class="alert alert-danger">Not all arguments were passed by the script.</div>';
								break;
							case 'pn_fail':
								echo '<div class="alert alert-danger">Invalid pack name was provided. Names must be at least one character, and no more than sixty-four characters. Allowed characters: <strong>a-zA-Z0-9._-()</strong> and <strong>[space]</strong></div>';
								break;
							case 'num_fail':
								echo '<div class="alert alert-danger">Minimum RAM and Permgen must be numeric.</div>';
								break;
							case 'ver_fail':
								echo '<div class="alert alert-danger">Invalid pack version was provided. Versions must be at least one character, and no more than thirty-two characters. Allowed characters: <strong>a-zA-Z0-9._-</strong></div>';
								break;						
						}
					
					}else if(isset($_GET['success'])){
					
						echo '<div class="alert alert-success">Successfully updated the Modpack settings.</div>';
					
					}
				
				?>
				<form action="ajax/modpack/update.php" method="post">
					<fieldset>
						<div class="well">
							<div class="form-group">
								<label for="pack_name" class="control-label">Modpack Name</label>
								<div>
									<input type="text" name="pack_name" class="form-control" value="<?php echo $pack['name']; ?>" />
								</div>
								<div class="checkbox">
									<label>
										<input type="checkbox" name="pack_default" <?php if($pack['default'] == 1){ echo 'checked="checked"'; } ?>/> Set this as the default modpack.
									</label>
								</div>
							</div>
							<div class="row">
								<div class="form-group col-4 nopad">
									<label for="pack_version" class="control-label">Modpack Version</label>
									<div class="input-group">
										<span class="input-group-addon">v.</span>
										<input type="text" name="pack_version" class="form-control" value="<?php echo $pack['version']; ?>" />
									</div>
								</div>
								<div class="form-group col-4 nopad">
									<label for="pack_minram" class="control-label">Minimum RAM Allocation to Install</label>
									<div class="input-group">
										<input type="text" name="pack_minram" class="form-control" value="<?php echo $pack['min_ram']; ?>" />
										<span class="input-group-addon">MB</span>
									</div>
								</div>
								<div class="form-group col-4 nopad">
									<label for="pack_permgen" class="control-label">PermGen Default Size</label>
									<div class="input-group">
										<input type="text" name="pack_permgen" class="form-control" value="<?php echo $pack['permgen']; ?>" />
										<span class="input-group-addon">MB</span>
									</div>
								</div>
							</div>
						</div>
						<div class="form-group">
							<div>
								<input type="hidden" name="pack_hash" value="<?php echo $pack['hash']; ?>" />
								<input type="submit" value="Update Modpack" class="btn btn-primary" />
								<button class="btn btn-danger btn-sm" data-toggle="modal" data-target="#delModpack">Delete</button>
							</div>
						</div>
					</fieldset>
				</form>
				<h3>Servers Using this Modpack</h3><hr />
				<table class="table table-striped table-bordered table-hover">
					<thead>
						<tr>
							<th>Name</th>
							<th>Node</th>
							<th>Connection Address</th>
						</tr>
					</thead>
					<tbody>
						<?php
							
							$servers = $mysql->prepare("SELECT * FROM `servers` WHERE `modpack` = :mpack");
							$servers->execute(array(
								':mpack' => $pack['hash']
							));
							
							while($row = $servers->fetch()){
								
								echo '<tr>
										<td><a href="../server/view.php?id='.$row['id'].'">'.$row['name'].'</a></td>
										<td>'.$core->settings->nodeName($row['node']).'</td>
										<td>'.$row['server_ip'].':'.$row['server_port'].'</td>
									</tr>';
							
							}
						
						?>
					</tbody>
				</table>
			</div>
			<div class="modal fade" id="delModpack" tabindex="-1" role="dialog" aria-labelledby="DeleteModpack" aria-hidden="true">
				<div class="modal-dialog">
					<form action="ajax/modpack/delete.php" method="post">
						<div class="modal-content">
							<div class="modal-header">
								<button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
								<h4 class="modal-title" id="DeleteModpack">Delete Modpack</h4>
							</div>
							<div class="modal-body">
								<?php 
									
									if(isset($_GET['error']) && $_GET['error'] == "pack_delete" && isset($_GET['disp']) && !empty($_GET['disp'])){
									
										switch($_GET['disp']){
										
											case 'missing_params':
												echo '<div class="alert alert-danger">Not all arguments were passed. This is often caused by not setting a new default pack. You must set a new default pack if deleting the current default pack. If there are no packs to select please upload a new one and set it as default.</div>';
												break;
											case 'pack_hash_mismatch':
												echo '<div class="alert alert-danger">The hash provided by the script did not match the hash that you manually entered.</div>';
												break;
											case 'pack_hash':
												echo '<div class="alert alert-danger">That pack hash does not appear to exist in the system.</div>';
												break;
											case 'no_new_default':
												echo '<div class="alert alert-danger">You must set a new default pack if deleting the current default pack. If there are no packs to select please upload a new one and set it as default.</div>';
												break;
											case 'new_default_noexist':
												echo '<div class="alert alert-danger">The new default pack you selected does not exist in the system.  If there are no packs to select please upload a new one and set it as default.</div>';
												break;						
										}
									
									}
								
								?>
								<div class="alert alert-danger">Deleting this Modpack is <strong>irreversible</strong>, all servers currently using it will be moved to the Modpack that you define below.</div>
								<?php if($pack['default'] == 1){ ?>
								<div class="row">
									<div class="form-group col-6">
										<label for="pack_newdefault" class="control-label">Set New Default Modpack</label>
										<div>
											<select class="form-control" name="pack_newdefault">
												<option value="no-continue" disabled="disbaled">-- Select a Modpack</option>
												<?php
												
													$packs = $mysql->prepare("SELECT * FROM `modpacks` WHERE `hash` != :mid");
													$packs->execute(array(':mid' => $_GET['mid']));
													
													if($packs->rowCount() > 0){
														
														while($row = $packs->fetch())
															echo '<option value="'.$row['hash'].'">'.$row['name'].' ('.$row['version'].')</option>';
															
													}else
														echo '<option value="no-continue" disabled="disbaled">No Avaliable Modpacks</option>';
												
												?>
											</select>
										</div>
									</div>
								</div>
								<?php } ?>
								<div class="row">
									<div class="form-group col-6">
										<label for="conf_pack_hash" class="control-label">Enter the Modpack Hash: <strong><?php echo $pack['hash']; ?></strong></label>
										<div>
											<input type="text" name="conf_pack_hash" class="form-control" />
										</div>
									</div>
								</div>
								<div class="row">
									<div class="form-group col-6">
										<div class="checkbox">
											<label>
												<input type="checkbox" name="confirm_delete" /> I have read and understand the above statements. I have correctly set a new Modpack to transfer servers to, and (if applicable) I have set a new default Modpack.
											</label>
										</div>
									</div>
								</div>
							</div>
							<div class="modal-footer">
								<input type="hidden" name="pack_del" value="<?php echo $pack['hash']; ?>" />
								<input type="submit" value="Delete" id="disable_complete" class="btn btn-danger" />
								<button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
							</div>
						</div>
					</form>
				</div>
			</div>
		</div>
		<div class="footer">
			<?php include('../../../assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$(document).ready(function(){
			setActiveOption('configuration-modpacks');
			if($.urlParam('error') != null){
				if($.urlParam('error') == "pack_delete"){
					$("#delModpack").modal('show')
				}
				var field = $.urlParam('error');
				var exploded = field.split('|');
					$.each(exploded, function(key, value) {
						$('[name="' + value + '"]').parent().parent().addClass('has-error');
					});
				var obj = $.parseJSON($.cookie('__TMP_pp_admin_updatemodpack'));
					$.each(obj, function(key, value) {
						$('[name="' + key + '"]').val(value);
					});			
			}
		});
	</script>
</body>
</html>