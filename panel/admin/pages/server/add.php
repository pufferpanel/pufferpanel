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

if(isset($_GET['do']) && $_GET['do'] == 'generate_password')
	exit($core->auth->keygen(12));
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
				<h3 class="nopad">Create New Server</h3><hr />
				<?php 
					
					if(isset($_GET['disp']) && !empty($_GET['disp'])){
					
						switch($_GET['disp']){
						
							case 'missing_args':
								echo '<div class="alert alert-danger">Not all arguments were passed by the script.</div>';
								break;
							case 's_fail':
								echo '<div class="alert alert-danger">The server name you entered does not meet the requirements. Must be at least 4 characters, and no more than 35. Server name can only contain a-zA-Z0-9_-</div>';
								break;
							case 'n_fail':
								echo '<div class="alert alert-danger">The node selected does not seem to exist.</div>';
								break;
							case 'no_modpack':
								echo '<div class="alert alert-danger">The modpack hash passed is not valid.</div>';
								break;
							case 'modpack_ram':
								echo '<div class="alert alert-danger">The modpack selected requires more RAM be allocated to it. Miimum amount of RAM that can be allocated: '.$_GET['min_ram'].' MB.</div>';
								break;
							case 'ip_fail':
								echo '<div class="alert alert-danger">The selected IP does not exist.</div>';
								break;
							case 'port_fail':
								echo '<div class="alert alert-danger">The selected port does not exist.</div>';
								break;
							case 'port_full':
								echo '<div class="alert alert-danger">The selected port is already in use.</div>';
								break;
							case 'e_fail':
								echo '<div class="alert alert-danger">The email you entered is invalid.</div>';
								break;
							case 'p_fail':
								echo '<div class="alert alert-danger">The passwords you entered did not match or were not at least 8 characters.</div>';
								break;
							case 'a_fail':
								echo '<div class="alert alert-danger">Account with that email does not exist in the system.</div>';
								break;
							case 'm_fail':
								echo '<div class="alert alert-danger">You entered a non-number for Disk and/or Memory.</div>';
								break;
							case 'b_fail':
								echo '<div class="alert alert-danger">You entered a non-number for Backup Files and/or Disk Space.</div>';
								break;
						
						}
					
					}
				
				?>
				<form action="ajax/new/create.php" method="POST">
					<fieldset>
						<div class="well">
							<div class="row">
								<div class="form-group col-6 nopad">
									<label for="server_name" class="control-label">Server Name</label>
									<div>
										<input type="text" autocomplete="off" name="server_name" class="form-control" />
										<p class="text-muted" style="margin: 0 0 -10.5px;"><small><em>Character Limits: a-zA-Z0-9_- (Max 35 characters)</em></small></p>
									</div>
								</div>
								<div class="form-group col-6 nopad-right">
									<label for="email" class="control-label">Owner Email</label>
									<div>
										<input type="text" autocomplete="off" name="email" value="<?php if(isset($_GET['email'])) echo $_GET['email']; ?>" class="form-control" />
									</div>
								</div>
							</div>
						</div>
						<div class="well">
							<div class="row">
								<div class="form-group col-6 nopad">
									<label for="server_name" class="control-label">Server Node</label>
									<div>
										<select name="node" id="getNode" class="form-control">
											<?php
												$selectData = $mysql->prepare("SELECT * FROM `nodes`");
												$selectData->execute(array());
												while($node = $selectData->fetch()){
												
													echo '<option value="'.$node['id'].'">'.$node['node'].'</option>';
												
												}
											?>
										</select>
									</div>
								</div>
								<div class="form-group col-6 nopad-right">
									<label for="modpack" class="control-label">Server Modpack</label>
									<div>
										<select name="modpack"class="form-control">
											<option value="none">None</option>
											<?php
												$packs = $mysql->prepare("SELECT * FROM `modpacks` WHERE `deleted` = 0");
												$packs->execute();
												while($pack = $packs->fetch()){
												
													echo '<option value="'.$pack['hash'].'">'.$pack['name'].' ('.$pack['version'].')</option>';
												
												}
											?>
										</select>
									</div>
								</div>
							</div>
							<div class="alert alert-warning" id="noPorts" style="display:none;margin-bottom:10px;"><strong>Error:</strong> This node does not have any free ports available. Please select another node.</div>
							<div class="row">
								<span id="updateList">
									<div class="form-group col-6 nopad">
										<label for="server_ip" class="control-label">Server IP</label>
										<div>
											<select name="server_ip" class="form-control">
												<option value="---">Select a Node</option>
											</select>
										</div>
									</div>
									<div class="form-group col-6 nopad-right">
										<label for="server_port" class="control-label">Server Port</label>
										<div>
											<select name="server_port" class="form-control">
												<option value="---">Select a Node</option>
											</select>
										</div>
									</div>
								</span>
							</div>
						</div>
						<div class="well">
							<div class="row">
								<div class="form-group col-4 nopad">
									<label for="alloc_mem" class="control-label">Allocate Memory (in MB)</label>
									<div class="input-group">
										<input type="text" autocomplete="off" name="alloc_mem" class="form-control" />
										<span class="input-group-addon">MB</span>
									</div>
								</div>
								<div class="form-group col-4">
									<label for="alloc_disk" class="control-label">Allocate Disk Space (in MB)</label>
									<div class="input-group">
										<input type="text" autocomplete="off" name="alloc_disk" class="form-control" />
										<span class="input-group-addon">MB</span>
									</div>
								</div>
								<div class="form-group col-4 nopad-right">
									<label for="cpu_limit" class="control-label">CPU Limit</label>
									<div class="input-group">
										<input type="text" autocomplete="off" name="cpu_limit" class="form-control" />
										<span class="input-group-addon">%</span>
									</div>
								</div>
							</div>
							<div class="row"><p><small>If you do not want to limit CPU usage set the value to 0. To determine a value, take the number <em>physical</em> cores and multiply it by 100. For example, on a quad core system <code>(4 * 100 = 400)</code> there is 400% available. To limit a server to using half of a single core, you would set the value to <code>50</code>. To allow a server to use up to two physical cores, set the value to <code>200</code>.</small></p></div>
						</div>
						<div class="well">
							<div class="row">
								<div class="alert alert-success" style="display: none;margin-bottom:10px;" id="gen_pass"></div>
								<div class="form-group col-6 nopad">
									<label for="sftp_pass" class="control-label">SFTP Password</label>
									<div class="input-group">
										<input type="password" autocomplete="off" name="sftp_pass" class="form-control" />
										<span class="input-group-btn">
											<button class="btn btn-success" id="gen_pass_bttn" type="button">Generate</button>
										</span>
									</div>
								</div>
								<div class="form-group col-6 nopad-right">
									<label for="sftp_pass_2" class="control-label">SFTP Password (Again)</label>
									<div>
										<input type="password" autocomplete="off" name="sftp_pass_2" class="form-control" />
									</div>
								</div>
							</div>
						</div>
						<div class="well">
							<input type="submit" value="Create Server" class="btn btn-primary btn-sm" />
						</div>
					</fieldset>
				</form>
			</div>
		</div>
		<div class="footer">
			<?php include('../../../assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		function updatePortList(){
			$("#server_ip").change(function() {
			    var ip = $(this).val().replace(/\./g, "\\.");
			    $("[id^=node_]").hide();
			    $("#node_"+ip).show();
			});
		}
		function updateList(){
			var activeNode = $('#getNode').val();
			$("#noPorts").hide();
			$.ajax({
				type: "POST",
				url: "ajax/new/load_list.php",
				data: {'node' : activeNode},
				success: function(data) {
					$('#updateList').html(data);
					updatePortList();
					return false;
				}
			});
			return false;
		}
		$(document).ready(function(){
			updateList();
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
			$('#getNode').change(function(){
				updateList();
			});
			if($.urlParam('error') != null){
			
				var field = $.urlParam('error');
				var exploded = field.split('|');
				
					$.each(exploded, function(key, value) {
						
						$('[name="' + value + '"]').parent().parent().addClass('has-error');
						
					});
					
				var obj = $.parseJSON($.cookie('__TMP_pp_admin_newserver'));
				
					$.each(obj, function(key, value) {
						
						$('[name="' + key + '"]').val(value);
						
					});
			
			}
		});
	</script>
</body>
</html>