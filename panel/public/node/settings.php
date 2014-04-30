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
require_once('../core/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false){

	Page\components::redirect($core->settings->get('master_url').'index.php?login');
	exit();
	
}

if(isset($_GET['do']) && $_GET['do'] == 'generate_password')
	exit($core->auth->keygen(rand(12, 18)));
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('../assets/include/header.php'); ?>
	<title>PufferPanel - Manage Your Server</title>
</head>
<body>
	<div class="container">
		<?php include('../assets/include/navbar.php'); ?>
		<div class="row">
			<div class="col-3">
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.acc_actions'); ?></strong></a>
					<a href="../account.php" class="list-group-item"><?php echo $_l->tpl('sidebar.settings'); ?></a>
					<a href="../servers.php" class="list-group-item"><?php echo $_l->tpl('sidebar.servers'); ?></a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.server_acc'); ?></strong></a>
					<a href="index.php" class="list-group-item"><?php echo $_l->tpl('sidebar.overview'); ?></a>
					<a href="console.php" class="list-group-item"><?php echo $_l->tpl('sidebar.console'); ?></a>
					<a href="files/index.php" class="list-group-item"><?php echo $_l->tpl('sidebar.files'); ?></a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.server_sett'); ?></strong></a>
					
					<a href="settings.php" class="list-group-item active"><?php echo $_l->tpl('sidebar.manage'); ?></a>
					<a href="plugins/index.php" class="list-group-item"><?php echo $_l->tpl('sidebar.plugins'); ?></a>
				</div>
			</div>
			<div class="col-9">
				<ul class="nav nav-tabs" id="config_tabs">
					<li class="active"><a href="#server_sett" data-toggle="tab"><?php echo $_l->tpl('string.modpacks'); ?></a></li>
					<li class=""><a href="#sftp_sett" data-toggle="tab"><?php echo $_l->tpl('string.sftp'); ?></a></li>
				</ul>
				<div class="tab-content">
					<div class="tab-pane active" id="server_sett">
					<br>
					<h3 class="nopad"><?php echo $_l->tpl('node.settings.modpack.modpack_h1'); ?></h3><hr />
						<form action="#" method="post" id="updateModpack">
							<fieldset>
								<div class="row" id="installingModpack" style="display:none;"></div>
								<div class="row">
									<div class="well well-sm">
										<?php
											$packs = $mysql->prepare("SELECT hash, name, version, deleted FROM `modpacks` WHERE `hash` = :hash");
											$packs->execute(array(
												':hash' => $core->server->getData('modpack')
											));
											
											$pack = $packs->fetch();
											$isDeleted = ($pack['deleted'] == 1) ? '[DELETED] ' : null;
										?>
										<small><?php echo $_l->tpl('node.settings.modpack.current'); ?>: <strong><?php echo $isDeleted.$pack['name'].' ('.$pack['version'].')'; ?></strong></small>
									</div>
								</div>
								<div class="form-group col-8 nopad">
									<div>
										<select class="form-control" name="new_pack">
											<option disabled="disabled"><?php echo $_l->tpl('node.settings.modpack.select'); ?></option>
											<?php
												$packs = $mysql->prepare("SELECT hash, name, version, server_jar FROM `modpacks` WHERE `deleted` = 0 AND `min_ram` <= :ram");
												$packs->execute(array(
													':ram' => $core->server->getData('max_ram')
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
										<input type="submit" id="install_modpack_submit" value="<?php echo $_l->tpl('node.settings.modpack.install'); ?>" class="btn btn-primary" />
									</div>
								</div>
							</fieldset>
						</form>
						<div class="row">
							<h3><?php echo $_l->tpl('node.settings.jar.jar_h1'); ?></h3><hr />
								<div class="well">
									<form action="ajax/settings/jarname.php" method="post">
										<fieldset>
										<div class="form-group">
											<label for="jarfile" class="control-label"><?php echo $_l->tpl('node.settings.jar.jarfile_name'); ?></label>
											<div class="input-group">
												<input type="text" autocomplete="off" name="jarfile" class="form-control" value="<?php echo str_replace('.jar' , '', $core->server->getData('server_jar')); ?>"/>
												<span class="input-group-addon">.jar</span>
												<span class="input-group-btn">
													<button class="btn btn-primary"><?php echo $_l->tpl('node.settings.jar.update'); ?></button>
												</span>
											</div>
										</div>
										</fieldset>
									</form>
								</div>
						</div>
					</div>
					<div class="tab-pane" id="sftp_sett">
						<h3><?php echo $_l->tpl('node.settings.jar.update'); ?></h3><hr />
						<form action="ajax/settings/sftp.php" method="post">
							<fieldset>
								<div class="form-group">
									<label for="sftp_host" class="control-label"><?php echo $_l->tpl('string.host'); ?></label>
									<div>
										<input type="text" readonly="readonly" value="<?php echo $core->server->nodeData('sftp_ip'); ?>" class="form-control" />
									</div>
								</div>
								<div class="form-group">
									<label for="sftp_user" class="control-label"><?php echo $_l->tpl('string.username'); ?></label>
									<div>
										<input type="text" readonly="readonly" value="<?php echo $core->server->getData('ftp_user'); ?>" class="form-control" />
									</div>
								</div>
								<div class="well">
									<div class="row">
										<div class="alert alert-success" style="display: none;margin-bottom:10px;" id="gen_pass"></div>
										<div class="form-group col-6 nopad">
											<label for="sftp_pass" class="control-label"><?php echo $_l->tpl('settings.new_pass'); ?></label>
											<div class="input-group">
												<input type="password" autocomplete="off" name="sftp_pass" class="form-control" />
												<span class="input-group-btn">
													<button class="btn btn-success" id="gen_pass_bttn" type="button"><?php echo $_l->tpl('string.generate'); ?></button>
												</span>
											</div>
										</div>
										<div class="form-group col-6 nopad-right">
											<label for="sftp_pass_2" class="control-label"><?php echo $_l->tpl('settings.new_pass').' '.$_l->tpl('string.again'); ?></label>
											<div>
												<input type="password" autocomplete="off" name="sftp_pass_2" class="form-control" />
											</div>
										</div>
									</div>
								</div>
								<input type="submit" value="<?php echo $_l->tpl('settings.change_pass'); ?>" class="btn btn-primary btn-sm" />
							</fieldset>
						</form>
					</div>
				</div>
			</div>
		</div>
		<div class="footer">
			<?php include('../assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$("#updateModpack").submit(function(e){
			e.preventDefault();
			var pack = $('select[name="new_pack"]').val();
			
			if(pack != null){
				$('#install_modpack_submit').addClass('disabled');
				$("#installingModpack").html('<div class="alert alert-info"><i class="fa fa-spinner fa-spin"></i> <?php echo $_l->tpl('node.settings.ajax.installing_modpack'); ?></div>').show();
				$.ajax({
					type: "POST",
					url: 'ajax/settings/modpack.php',
					data: { new_pack: pack },
			  		success: function(data) {
			    		$("#installingModpack").hide().html(data).show();
			    		$('#install_modpack_submit').removeClass('disabled');
			 		}
				});
			}else
				alert('<?php echo $_l->tpl('node.settings.ajax.modpack_error'); ?>');
		});
		$("#gen_pass_bttn").click(function(){
				$.ajax({
					type: "GET",
					url: "settings.php?do=generate_password",
					success: function(data) {
						$("#gen_pass").html('<strong><?php echo $_l->tpl('string.generated_pass'); ?>:</strong> '+data);
						$("#gen_pass").slideDown();
						$('input[name="sftp_pass"]').val(data);
						$('input[name="sftp_pass_2"]').val(data);
						return false;
					}
				});
				return false;
			});
	</script>
</body>
</html>