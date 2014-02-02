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
					<a href="console.php" class="list-group-item">Live Console</a>
					<a href="files/index.php" class="list-group-item">File Manager</a>
					<a href="backup.php" class="list-group-item">Backup Manager</a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Server Settings</strong></a>
					<a href="settings.php" class="list-group-item active">Modpack Management</a>
					<a href="plugins/index.php" class="list-group-item">Server Plugins</a>
				</div>
			</div>
			<div class="col-9">
				<h3 class="nopad">Server Modpack Management</h3><hr />
				<form action="#" method="post" id="updateModpack">
					<fieldset>
						<div class="row" id="installingModpack" style="display:none;"></div>
						<div class="row">
							<div class="well well-sm">
								<?php
									$packs = $mysql->prepare("SELECT hash, name, version, deleted FROM `modpacks` WHERE `hash` = :hash");
									$packs->execute(array(
										':hash' => $core->framework->server->getData('modpack')
									));
									
									$pack = $packs->fetch();
									$isDeleted = ($pack['deleted'] == 1) ? '[DELETED] ' : null;
								?>
							 	<small>Current Modpack: <strong><?php echo $isDeleted.$pack['name'].' ('.$pack['version'].')'; ?></strong></small>
							</div>
						</div>
						<div class="form-group col-8 nopad">
							<div>
								<select class="form-control" name="new_pack">
									<option disabled="disabled">-- Select a Modpack</option>
									<?php
										$packs = $mysql->prepare("SELECT hash, name, version FROM `modpacks` WHERE `deleted` = 0 AND `min_ram` <= :ram");
										$packs->execute(array(
											':ram' => $core->framework->server->getData('max_ram')
										));
										
										while($row = $packs->fetch()){
											
											if($row['hash'] != $pack['hash'])
												echo '<option value="'.$row['hash'].'">'.$row['name'].' ('.$row['version'].')</option>';
											else
												echo '<option disabled="disabled">'.$row['name'].' ('.$row['version'].')</option>';
										}
									?>
								</select>
							</div>
						</div>
						<div class="form-group col-4 nopad-right">
							<div>
								<input type="submit" id="install_modpack_submit" value="Install Modpack" class="btn btn-primary" />
							</div>
						</div>
					</fieldset>
				</form>
			</div>
		</div>
		<div class="footer">
			<?php include('assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$("#updateModpack").submit(function(e){
			e.preventDefault();
			var pack = $('select[name="new_pack"]').val();
			
			if(pack != null){
				$('#install_modpack_submit').addClass('disabled');
				$("#installingModpack").html('<div class="alert alert-info"><i class="fa fa-spinner fa-spin"></i> Installing Modpack.</div>').show();
				$.ajax({
					type: "POST",
					url: 'core/ajax/settings/modpack.php',
					data: { new_pack: pack },
			  		success: function(data) {
			    		$("#installingModpack").hide().html(data).show();
			    		$('#install_modpack_submit').removeClass('disabled');
			 		}
				});
			}else
				alert('You can not use a disabled form input as a modpack!');
		});
	</script>
</body>
</html>