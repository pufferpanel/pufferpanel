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
	exit($core->framework->auth->keygen(12));

if(!isset($_GET['id']))
	$core->framework->page->redirect('find.php');

/*
 * Select User Information
 */
$selectUser = $mysql->prepare("SELECT * FROM `users` WHERE `id` = :id LIMIT 1");
$selectUser->execute(array(
	':id' => $_GET['id']
));

	if($selectUser->rowCount() != 1)
		$core->framework->page->redirect('find.php?error=no_user');
	else
		$user = $selectUser->fetch();
		
/*
 * Select Servers Owned by the User
 */
$selectServers = $mysql->prepare("SELECT * FROM `servers` WHERE `owner_id` = :id ORDER BY `active` DESC");
$selectServers->execute(array(
	':id' => $user['id']
));

?>
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>PufferPanel - Viewing Account `<?php echo $user['username']; ?>`</title>
	
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
						<h3 class="fl">Viewing Account Information for `<?php echo $user['username']; ?>`</h3>
					</div>
					<?php 
						
						if(isset($_GET['disp']) && !empty($_GET['disp'])){
						
							switch($_GET['disp']){
							
								case 'p_updated':
									echo '<div class="content-module-main" id="fadeOut"><div class="confirmation-box">Account password has been updated.</div></div>';
									break;
								case 'd_updated':
									echo '<div class="content-module-main" id="fadeOut"><div class="confirmation-box">Account email & administrator status updated.</div></div>';
									break;
							
							}
						
						}
					
					?>
				</div>
			</div>
			<div class="side-content fr">
				<div class="half-size-column fl">
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Account Information</h3>
						</div>
						<div class="content-module-main cf">
							<form action="ajax/account/update.php" method="POST">
								<fieldset>
									<p>
										<label for="username">Account Username</label>
										<input type="text" name="username" value="<?php echo $user['username']; ?>" readonly="readonly" class="round full-width-input" />
									</p>
									<p>
										<?php
										
											$date1 = new DateTime(date('Y-m-d', $user['register_time']));
											$date2 = new DateTime(date('Y-m-d', time()));
											
											$totalDays = $date2->diff($date1)->format("%a Days");
										
										?>
										<label for="registered">Account Registered</label>
										<input type="text" name="registered" value="<?php echo date('F j, Y g:ia', $user['register_time']); ?> (<?php echo $totalDays; ?>)" readonly="readonly" class="round full-width-input" />
									</p>
									<p>
										<label for="email">Account Email</label>
										<input type="text" name="email" value="<?php echo $user['email']; ?>" class="round full-width-input" />
									</p>
									<p>
										<?php
										
											$isSelected = array();
											$isSelected['no'] = ($user['root_admin'] == 0) ? 'selected="selected"' : '';
											$isSelected['yes'] = ($user['root_admin'] == 1) ? 'selected="selected"' : '';
										
										?>
										<label for="root_admin">Root Administrator</label>
										<select name="root_admin">
											<option value="0" <?php echo $isSelected['no']; ?>>No</option>
											<option value="1" <?php echo $isSelected['yes']; ?>>Yes</option>
										</select><i class="fa fa-angle-down pull-right select-arrow select-halfbox"></i>
										<em>Warning: setting this to yes will give the user full access to this Admin CP. Only set this to yes if the person is a member of your company.</em>
									</p>
									<div class="stripe-separator"><!--  --></div>
									<input type="hidden" name="uid" value="<?php echo $_GET['id']; ?>" />
									<input type="hidden" name="action" value="details" />
									<input type="submit" value="Update User" class="round blue ic-right-arrow" />
								</fieldset>
							</form>
						</div>
					</div>
				</div>
				<div class="half-size-column fr">
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Password Update</h3>
						</div>
						<div class="content-module-main">
							<form action="ajax/account/password.php" method="post">
								<fieldset>
									<div class="warning-box round" style="display: none;" id="gen_pass"></div>
									<p>
										<label for="pass">New Password (<a href="#" id="gen_pass_bttn">Generate</a>)</label>
										<input type="password" name="pass" class="round full-width-input" />
									</p>
									<p>
										<label for="pass_2">New Password (Again)</label>
										<input type="password" name="pass_2" class="round full-width-input" />
									</p>
									<p>
										<input type="checkbox" name="email_user" /> Email new password to user.
									</p>
									<p>
										<input type="checkbox" name="clear_session" /> Clear the user's session.
									</p>
									<div class="stripe-separator"><!--  --></div>
									<input type="hidden" name="uid" value="<?php echo $_GET['id']; ?>" />
									<input type="hidden" name="email" value="<?php echo $user['email']; ?>" />
									<input type="submit" value="Update User" class="round blue ic-right-arrow" />
								</fieldset>
							</form>
						</div>
					</div>
				</div>
			</div>
			<div class="side-content fr">
				<div class="content-module">
					<div class="content-module-heading cf">
						<h3 class="fl">Servers</h3>
					</div>
					<div class="content-module-main cf">
						<?php
							$listServers = '';
							
							while($row = $selectServers->fetch()){
							
								($row['active'] == '1') ? $isActive = 'Enabled' : $isActive = 'Disabled';
								$listServers .= '
												<tr>
													<td><a href="../server/view.php?id='.$row['id'].'">'.$row['name'].'</a></td>
													<td>'.$row['node'].'</td>
													<td>'.$row['server_ip'].'</td>
													<td>'.$row['server_port'].'</td>
													<td>'.$row['max_ram'].' MB</td>
													<td>'.$row['disk_space'].' MB</td>
													<td>'.$isActive.'</td>
												</tr>
												';
							
							}
						?>
						<table>
							<thead>
								<tr>
									<th style="width:15%">Server Name</th>
									<th style="width:10%">Node</th>
									<th style="width:20%">IP Address</th>
									<th style="width:10%">Port</th>
									<th style="width:15%">Memory</th>
									<th style="width:15%">Disk Space</th>
									<th style="width:15%">Status</th>
								</tr>
							</thead>
							<tbody>
								<?php echo $listServers; ?>
							</tbody>
						</table>
						<p><a href="../server/add.php?email=<?php echo $user['email']; ?>">Add New Server</a></p>
					</div>
				</div>
			</div>
		</div>
	</div>
	<script type="text/javascript">
		$("#gen_pass_bttn").click(function(){
			$.ajax({
				type: "GET",
				url: "view.php?do=generate_password",
				success: function(data) {
					$("#gen_pass").html('Generated Password: '+data);
					$("#gen_pass").slideDown();
					$('input[name="pass"]').val(data);
					$('input[name="pass_2"]').val(data);
					return false;
				}
			});
			return false;
		});
		$(document).ready(function(){
			$('#fadeOut').delay(5000).slideUp();
		});
	</script>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.3 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>