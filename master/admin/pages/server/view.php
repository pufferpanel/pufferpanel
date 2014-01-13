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

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	$core->framework->page->redirect('../../../index.php');
}

if(isset($_GET['do']) && $_GET['do'] == 'generate_password')
	exit($core->framework->auth->keygen(rand(12, 18)));

if(!isset($_GET['id']))
	$core->framework->page->redirect('find.php');

/*
 * Select Servers Information
 */
$selectServers = $mysql->prepare("SELECT * FROM `servers` WHERE `id` = :id ORDER BY `active` DESC");
$selectServers->execute(array(
	':id' => $_GET['id']
));

	if($selectServers->rowCount() != 1)
		$core->framework->page->redirect('find.php?error=no_server');
	else
		$server = $selectServers->fetch();
		
$select = $mysql->prepare("SELECT * FROM `users` WHERE `id` = :oid");
$select->execute(array(
	':oid' => $server['owner_id']
));

	if($select->rowCount() != 1)
		$core->framework->page->redirect('find.php?error=no_server_user');
	else
		$user = $select->fetch();

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
					<li class="active"><a href="#info" data-toggle="tab">Connection</a></li>
					<li><a href="#server_sett" data-toggle="tab">Settings</a></li>
					<li><a href="#backup_sett" data-toggle="tab">Backups</a></li>
					<li><a href="#sftp_sett" data-toggle="tab">SFTP</a></li>
					<li><a href="../../../servers.php?goto=<?php echo $server['hash']; ?>">Server Control</a></li>
				</ul>
				<div class="tab-content">
					<div class="tab-pane active" id="info">
						<h3>Connection Information</h3><hr />
						<form action="ajax/server/connection.php" method="POST">
							<?php
								$selectData = $mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :name");
								$selectData->execute(array(
									':name' => $server['node']
								));
								$node = $selectData->fetch();
							?>
							<fieldset>
								<div class="panel panel-default">
									<div class="panel-heading">Changing IP &amp; Port</div>
									<div class="panel-body">
										<p>If you want to change the Server IP then select an IP from the list below that has at least one available port. When you select a new IP you will be prompted to select a new port from a list. If you only want to change the port, and not the IP, then you can do so by simply selecting an available port.</p>
									</div>
								</div>
								<div class="form-group col-6 nopad">
									<label for="server_ip" class="control-label">Server IP</label>
									<div>
										<select name="server_ip" id="server_ip" class="form-control">
										    <?php
											
												$ips = json_decode($node['ips'], true);
												foreach($ips as $ip => $internal){
												
										            if($server['server_ip'] == $ip)
										                echo '<option value="'.$ip.'" selected="selected">'.$ip.' ('.$internal['ports_free'].' Avaliable Ports)</option>';
										            else{
										            
										                if($internal['ports_free'] > 0)
														  echo '<option value="'.$ip.'">'.$ip.' ('.$internal['ports_free'].' Avaliable Ports)</option>';
										                else
														  echo '<option disabled="disabled">'.$ip.' ('.$internal['ports_free'].' Avaliable Ports)</option>';
										            
										            }
										            												
												}
											
											?>
										</select>
									</div>
								</div>
								<div class="form-group col-6 nopad-right">
									<label for="server_ip" class="control-label">Server Port</label>
									<div>
										<?php
										
										    $ports = json_decode($node['ports'], true);
										    
										    foreach($ports as $ip => $internal){
										    
										        if($server['server_ip'] == $ip)
										            echo '<span id="node_'.$ip.'"><select name="server_port_'.$ip.'" class="form-control">';
										        else
										            echo '<span style="display:none;" id="node_'.$ip.'"><select name="server_port_'.$ip.'" class="form-control">';
										        
										            foreach($internal as $port => $avaliable){
										            
										                if($server['server_port'] == $port)
										                    echo '<option value="'.$port.'" selected="selected">'.$port.'</option>';
										                else{
										                    
										                        if($avaliable == 1)
										                            echo '<option value="'.$port.'">'.$port.'</option>';   
										                        else
										                            echo '<option disabled="disabled">'.$port.'</option>';
										                    
										                }
										                
										            }
										        
										        echo '</select></span>';
										    
										    }
																							
										?>
									</div>
								</div>
								<input type="hidden" name="sid" value="<?php echo $_GET['id']; ?>" />
								<input type="hidden" name="nid" value="<?php echo $node['id']; ?>" />
								<input type="submit" value="Update Server" class="btn btn-primary btn-sm" />
							</fieldset>
						</form>
					</div>
					<div class="tab-pane" id="server_sett">
						<h3>Server Settings</h3><hr />
						<form action="ajax/server/allocate.php" method="post">
							<fieldset>
								<div class="form-group">
									<label class="control-label">Owner Email</label>
									<div>
										<input type="text" readonly="readonly" value="<?php echo $user['username']; ?> (<?php echo $user['email']; ?>)" class="form-control" />
									</div>
								</div>
								<div class="row">
									<div class="form-group col-6 nopad">
										<label for="alloc_mem" class="control-label">Allocate RAM</label>
										<div class="input-group">
											<input type="text" autocomplete="off" name="alloc_mem" value="<?php echo $server['max_ram']; ?>" class="form-control" />
											<span class="input-group-addon">MB</span>
										</div>
									</div>
									<div class="form-group col-6 nopad-right">
										<label for="alloc_disk" class="control-label">Disk Space</label>
										<div class="input-group">
											<input type="text" name="alloc_disk" value="<?php echo $server['disk_space']; ?>" class="form-control" />
											<span class="input-group-addon">MB</span>
										</div>
									</div>
								</div>
								<input type="hidden" name="sid" value="<?php echo $_GET['id']; ?>" />
								<input type="submit" value="Update Server Settings" class="btn btn-primary btn-sm" />
							</fieldset>
						</form>
					</div>
					<div class="tab-pane" id="backup_sett">
						<h3>Backup Settings</h3><hr />
						<form action="ajax/server/backup.php" method="post">
							<fieldset>
								<div class="row">
									<div class="form-group col-6 nopad">
										<label for="backup_disk" class="control-label">Maximum Backup Space (in MB)</label>
										<div class="input-group">
											<input type="text" autocomplete="off" readonly="readonly" value="<?php echo $server['backup_disk_limit']; ?>" name="backup_disk" class="form-control" />
											<span class="input-group-addon">MB</span>
										</div>
									</div>
									<div class="form-group col-6 nopad-right">
										<label for="max_files" class="control-label">Maximum Compressed Backups</label>
										<div>
											<input type="text" autocomplete="off" readonly="readonly" name="max_files" value="<?php echo $server['backup_file_limit']; ?>" class="form-control" />
										</div>
									</div>
								</div>
								<input type="hidden" name="sid" value="<?php echo $_GET['id']; ?>" />
								<input type="submit" disabled="disabled" value="Update Backup Information" class="btn btn-primary btn-sm disabled" />
							</fieldset>
						</form>
					</div>
					<div class="tab-pane" id="sftp_sett">
						<h3>SFTP Settings</h3><hr />
						<form action="ajax/server/sftp.php" method="post">
							<fieldset>
								<div class="form-group">
									<label for="sftp_host" class="control-label">Host</label>
									<div>
										<input type="text" readonly="readonly" value="<?php echo $server['ftp_host']; ?>" class="form-control" />
									</div>
								</div>
								<div class="form-group">
									<label for="sftp_user" class="control-label">Username</label>
									<div>
										<input type="text" readonly="readonly" value="<?php echo $server['ftp_user']; ?>" class="form-control" />
									</div>
								</div>
								<div class="well">
									<div class="row">
										<div class="alert alert-success" style="display: none;margin-bottom:10px;" id="gen_pass"></div>
										<div class="form-group col-6 nopad">
											<label for="sftp_pass" class="control-label">New Password</label>
											<div class="input-group">
												<input type="password" autocomplete="off" name="sftp_pass" class="form-control" />
												<span class="input-group-btn">
													<button class="btn btn-success" id="gen_pass_bttn" type="button">Generate</button>
												</span>
											</div>
										</div>
										<div class="form-group col-6 nopad-right">
											<label for="sftp_pass_2" class="control-label">New Password (Again)</label>
											<div>
												<input type="password" autocomplete="off" name="sftp_pass_2" class="form-control" />
											</div>
											<div class="checkbox" style="margin-bottom:-20px;">
												<label>
													<input type="checkbox" name="email_user" /> Email new password to user.
												</label>
											</div>
										</div>
									</div>
								</div>
								<input type="hidden" name="sid" value="<?php echo $_GET['id']; ?>" />
								<input type="hidden" name="nid" value="<?php echo $node['id']; ?>" />
								<input type="submit" value="Update Password" class="btn btn-primary btn-sm" />
							</fieldset>
						</form>
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
			setActiveOption('server-find');
			if($.urlParam('error') != null){
				var field = $.urlParam('error');
				var exploded = field.split('|');
					$.each(exploded, function(key, value) {
						$('[name="' + value + '"]').parent().parent().addClass('has-error');
					});
			}
			if($.urlParam('tab') != null){
				$('#config_tabs a[href="#'+$.urlParam('tab')+'"]').tab('show');
			}
			$("#server_ip").change(function() {
			    var ip = $(this).val().replace(/\./g, "\\.");
			    $("[id^=node_]").hide();
			    $("#node_"+ip).show();
			});
			$("#gen_pass_bttn").click(function(){
				$.ajax({
					type: "GET",
					url: "add.php?do=generate_password",
					success: function(data) {
						$("#gen_pass").html('<strong>Generated Password:</strong> '+data);
						$("#gen_pass").slideDown();
						$('input[name="sftp_pass"]').val(data);
						$('input[name="sftp_pass_2"]').val(data);
						return false;
					}
				});
				return false;
			});
		});
	</script>
</body>
</html>