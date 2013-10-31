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
	<meta charset="utf-8">
	<title>PufferPanel - Viewing Server `<?php echo $server['name']; ?>`</title>
	
	<!-- Stylesheets -->
	<link href='http://fonts.googleapis.com/css?family=Droid+Sans:400,700' rel='stylesheet'>
	<link rel="stylesheet" href="../../../assets/css/style.css">
	
	<!-- Optimize for mobile devices -->
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
	
	<!-- jQuery & JS files -->
	<script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
	
</head>
<body>
	<div id="top-bar">
		<div class="page-full-width cf">
			<ul id="nav" class="fl">
				<li><a href="../../../account.php" class="round button dark"><i class="fa fa-user"></i>&nbsp;&nbsp; <strong><?php echo $core->framework->user->getData('username'); ?></strong></a></li>
			</ul>
			<ul id="nav" class="fr">
				<li><a href="../../../servers.php" class="round button dark"><i class="fa fa-sign-out"></i></a></li>
				<li><a href="../../../logout.php" class="round button dark"><i class="fa fa-power-off"></i></a></li>
			</ul>
		</div>	
	</div>
	<div id="header-with-tabs">
		<div class="page-full-width cf">
		</div>
	</div>
	<div id="content">
		<div class="page-full-width cf">
			<?php include('../../../core/templates/admin_sidebar.php'); ?>
			<div class="side-content fr">
				<div class="content-module">
					<div class="content-module-heading cf">
						<h3 class="fl">Viewing Server `<?php echo $server['name']; ?>`</h3>
					</div>
				</div>
			</div>
			<div class="side-content fr">
				<div class="half-size-column fl">
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Server Information</h3>
						</div>
						<div class="content-module-main cf">
							<form action="ajax/server/connection.php" method="POST">
								<?php
									$selectData = $mysql->prepare("SELECT * FROM `nodes` WHERE `node` = :name");
									$selectData->execute(array(
										':name' => $server['node']
									));
									$node = $selectData->fetch();
								?>
								<fieldset>
									<p><a href="../../../servers.php?goto=<?php echo $server['hash']; ?>">Click here</a> to access server control tools.</p>
									<div class="information-box no-image round">If you want to change the Server IP then select an IP from the list below that has at least one available port. When you select a new IP you will be prompted to select a new port from a list. If you only want to change the port, and not the IP, then you can do so by simply selecting an available port.</div>
									<p>
										<label for="server_ip">Server IP</label>
										<input type="text" name="server_ip" id="server_ip" value="<?php echo $server['server_ip']; ?>" class="round full-width-input" />
										<div class="grey-box round">
											<?php
											
												$ips = json_decode($node['ips'], true);
												foreach($ips as $ip => $internal){
												
													if($internal['ports_free'] > -1)
														echo '<a href="#/!setIP" class="click_ip" id="'.$ip.'">'.$ip.'</a> has '.$internal['ports_free'].' avaliable port(s).<br />';
													else
														echo $ip.' has '.$internal['ports_free'].' avaliable port(s).<br />';
												
												}
											
											?>
										</div>
									</p>
									<p>
										<label for="server_port">Server Port</label>
										<input type="text" name="server_port" id="server_port" value="<?php echo $server['server_port']; ?>" class="round full-width-input" />
											<?php
											
												$ports = json_decode($node['ports'], true);
												
												foreach($ports as $ip => $internal){
												
													if($server['server_ip'] == $ip)
														echo '<div class="grey-box round" id="node_'.$ip.'">';
													else
														echo '<div class="grey-box round" style="display:none;" id="node_'.$ip.'">';
													
													foreach($internal as $port => $avaliable){
													
														if($avaliable == 1)
															echo '<a href="#/!setPort" class="click_port" id="'.$port.'">'.$port.'</a> is <span style="color:#52964f;">avaliable</span>.<br />';
														else
															echo $port.' is <span style="color:#cf4425;">unavaliable</span>.<br />';
														
													}
													
													echo '</div>';
												
												}
																							
											?>
									</p>
									<div class="stripe-separator"><!--  --></div>
									<input type="hidden" name="sid" value="<?php echo $_GET['id']; ?>" />
									<input type="submit" value="Update Server" class="round blue ic-right-arrow" />
								</fieldset>
							</form>
						</div>
					</div>
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Backup Settings</h3>
						</div>
						<div class="content-module-main">
							<form action="ajax/server/backup.php" method="post">
								<fieldset>
									<p>
										<label for="max_files">Maximum Files</label>
										<input type="text" readonly="readonly" value="<?php echo $server['backup_file_limit']; ?>" class="round full-width-input" />
									</p>
									<p>
										<label for="pass">Maximum Space (In Megabytes)</label>
										<input type="text" readonly="readonly" value="<?php echo $server['backup_disk_limit']; ?>" class="round full-width-input" />
									</p>
									<div class="stripe-separator"><!--  --></div>
									<input type="hidden" name="sid" value="<?php echo $_GET['id']; ?>" />
									<input type="submit" value="Update Backup Information" class="round blue ic-right-arrow" />
								</fieldset>
							</form>
						</div>
					</div>
				</div>
				<div class="half-size-column fr">
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">SFTP Information</h3>
						</div>
						<div class="content-module-main">
							<form action="ajax/server/sftp.php" method="post">
								<fieldset>
									<p>
										<label for="sftp_host">Host</label>
										<input type="text" readonly="readonly" value="<?php echo $server['ftp_host']; ?>" class="round full-width-input" />
									</p>
									<p>
										<label for="sftp_user">Username</label>
										<input type="text" readonly="readonly" value="<?php echo $server['ftp_user']; ?>" class="round full-width-input" />
									</p>
									<div class="stripe-separator"><!--  --></div>
									<div class="warning-box round" style="display: none;" id="gen_pass"></div>
									<p>
										<label for="sftp_pass">New Password (<a href="#" id="gen_pass_bttn">Generate</a>)</label>
										<input type="password" name="sftp_pass" class="round full-width-input" />
										<em>Minimum length is 8 characters. 12 or more is suggested for highest security.</em>
									</p>
									<p>
										<label for="sftp_pass_2">New Password (Again)</label>
										<input type="password" name="sftp_pass_2" class="round full-width-input" />
									</p>
									<p>
										<input type="checkbox" name="email_user" /> Email new password to user.
									</p>
									<div class="stripe-separator"><!--  --></div>
									<input type="hidden" name="sid" value="<?php echo $_GET['id']; ?>" />
									<input type="submit" value="Update Password" class="round blue ic-right-arrow" />
								</fieldset>
							</form>
						</div>
					</div>
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Server Settings</h3>
						</div>
						<div class="content-module-main">
							<form action="ajax/server/settings.php" method="post">
								<fieldset>
									<p>
										<label for="owner">Owner</label>
										<input type="text" readonly="readonly" value="<?php echo $user['username']; ?> (<?php echo $user['email']; ?>)" class="round full-width-input" />
										<input type="hidden" name="owner_id" value="<?php echo $server['owner_id']; ?>" />
									</p>
									<div class="stripe-separator"><!--  --></div>
									<div class="warning-box round" style="display: none;" id="gen_pass"></div>
									<p>
										<label for="alloc_ram">Allocate RAM (in Megabytes)</label>
										<input type="text" name="pass" value="<?php echo $server['max_ram']; ?>" class="round full-width-input" />
									</p>
									<p>
										<label for="alloc_disk">Disk Space (in Megabytes)</label>
										<input type="text" name="pass_2" value="<?php echo $server['disk_space']; ?>" class="round full-width-input" />
									</p>
									<div class="stripe-separator"><!--  --></div>
									<input type="hidden" name="sid" value="<?php echo $_GET['id']; ?>" />
									<input type="submit" value="Update Server Settings" class="round blue ic-right-arrow" />
								</fieldset>
							</form>
						</div>
					</div>
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Server Settings</h3>
						</div>
						<div class="content-module-main">
							<p>Do the diddily-do here for banning a server, or suspending it.</p>
						</div>
					</div>
				</div>
			</div>
			</div>
		</div>
	</div>
	<script type="text/javascript">
		$(".click_ip").click(function(e){
			e.preventDefault();
			var start = $('#server_ip').val().replace(/\./g, "\\.");
			var ip = $(this).attr('id').replace(/\./g, "\\.");
			$('#server_ip').val(ip.replace(/\\\./g, "."));
			$('#node_'+start).slideToggle();
			$('#node_'+ip).slideToggle();
		});
		$(".click_port").click(function(e){
			e.preventDefault();
			var port = $(this).attr('id');
			$('#server_port').val(port);
		});
		$("#gen_pass_bttn").click(function(){
			$.ajax({
				type: "GET",
				url: "view.php?do=generate_password",
				success: function(data) {
					$("#gen_pass").html('Generated Password: '+data);
					$("#gen_pass").slideDown();
					$('input[name="sftp_pass"]').val(data);
					$('input[name="sftp_pass_2"]').val(data);
					return false;
				}
			});
			return false;
		});
	</script>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.3 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>