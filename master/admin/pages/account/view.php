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
	exit($core->framework->auth->keygen(rand(12,18)));

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
				<div class="row">
					<h3 class="nopad">Viewing User: <?php echo $user['username']; ?></h3><hr />
					<?php 
						
						if(isset($_GET['disp']) && !empty($_GET['disp'])){
						
							switch($_GET['disp']){
							
								case 'p_updated':
									echo '<div class="alert alert-success">Account password has been updated.</div>';
									break;
								case 'd_updated':
									echo '<div class="alert alert-success">Account email & administrator status updated.</div>';
									break;
							
							}
						
						}
					
					?>
					<div class="col-6 nopad">
						<form action="ajax/account/update.php" method="POST">
							<fieldset>
								<div class="form-group">
									<label for="username" class="control-label">Username</label>
									<div>
										<input type="text" name="username" value="<?php echo $user['username']; ?>" readonly="readonly" class="form-control" />
									</div>
								</div>
								<div class="form-group">
									<label for="email" class="control-label">Email</label>
									<div>
										<input type="text" name="email" value="<?php echo $user['email']; ?>" class="form-control" />
									</div>
								</div>
								<div class="form-group">
									<?php
									
										$date1 = new DateTime(date('Y-m-d', $user['register_time']));
										$date2 = new DateTime(date('Y-m-d', time()));
										
										$totalDays = $date2->diff($date1)->format("%a Days Ago");
									
									?>
									<label for="registered" class="control-label">Registered</label>
									<div>
										<input type="text" name="registered" value="<?php echo date('F j, Y g:ia', $user['register_time']); ?> (<?php echo $totalDays; ?>)" readonly="readonly" class="form-control" />
									</div>
								</div>
								<div class="form-group">
									<?php
									
										$isSelected = array();
										$isSelected['no'] = ($user['root_admin'] == 0) ? 'selected="selected"' : '';
										$isSelected['yes'] = ($user['root_admin'] == 1) ? 'selected="selected"' : '';
									
									?>
									<label for="email" class="control-label">Root Administrator</label>
									<div>
										<select name="root_admin" class="form-control">
											<option value="0" <?php echo $isSelected['no']; ?>>No</option>
											<option value="1" <?php echo $isSelected['yes']; ?>>Yes</option>
										</select>
										<p><small class="text-muted"><em><strong>Warning:</strong> Setting this to "Yes" gives a user full administrative access to PufferPanel.</em></small></p>
									</div>
								</div>
								<div class="form-group">
									<div>
										<input type="hidden" name="uid" value="<?php echo $_GET['id']; ?>" />
										<input type="hidden" name="action" value="details" />
										<input type="submit" value="Update User" class="btn btn-primary" />
									</div>
								</div>
							</fieldset>
						</form>
					</div>
					<div class="col-6 nopad-right">
						<div class="well">
							<h5 class="nopad">Update Password</h5><hr />
							<form action="ajax/account/password.php" method="post">
								<fieldset>
									<div class="alert alert-success" style="display:none;margin-bottom:10px;" id="gen_pass"></div>
									<div class="form-group">
										<label for="pass" class="control-label">New Password</label>
										<div>
											<input type="password" name="pass" class="form-control" />
										</div>
									</div>
									<div class="form-group">
										<label for="pass_2" class="control-label">New Password (Again)</label>
										<div>
											<input type="password" name="pass_2" class="form-control" />
											<div class="checkbox">
												<label>
													<input type="checkbox" name="email_user" /> Email new password to user.
												</label>
											</div>
											<div class="checkbox">
												<label>
													<input type="checkbox" name="clear_session" /> Clear the user's session.
												</label>
											</div>
										</div>
									</div>
									<div class="form-group">
										<div>
											<input type="hidden" name="uid" value="<?php echo $_GET['id']; ?>" />
											<input type="hidden" name="email" value="<?php echo $user['email']; ?>" />
											<button class="btn btn-primary" type="submit">Change Password</button>
											<button class="btn btn-default" id="gen_pass_bttn" type="button">Generate Password</button>
										</div>
									</div>
								</fieldset>
							</form>
						</div>
					</div>
				</div>
				<div class="row">
					<h3>Account Servers</h3><hr />
					<?php
						$listServers = '';
						
						while($row = $selectServers->fetch()){
						
							($row['active'] == '1') ? $isActive = '<span class="label label-success">Enabled</span>' : $isActive = '<span class="label label-danger">Disabled</span>';
							$listServers .= '
											<tr>
												<td><a href="../../../servers.php?goto='.$row['hash'].'"><i class="fa fa-tachometer"></i></a></td>
												<td><a href="../server/view.php?id='.$row['id'].'">'.$row['name'].'</a></td>
												<td>'.$core->framework->settings->nodeName($row['node']).'</td>
												<td>'.$row['server_ip'].':'.$row['server_port'].'</td>
												<td>'.$isActive.'</td>
											</tr>
											';
						
						}
					?>
					<table class="table table-striped table-bordered table-hover">
						<thead>
							<tr>
								<th style="width:2%;"></th>
								<th>Server Name</th>
								<th>Node</th>
								<th>Connection</th>
								<th style="width:10%;"></th>
							</tr>
						</thead>
						<tbody>
							<?php echo $listServers; ?>
						</tbody>
					</table>
					<button onclick="location.href='../server/add.php?email=<?php echo $user['email']; ?>';" class="btn btn-success btn-sm">Add New Server</button>
				</div>
			</div>
		</div>
		<div class="footer">
			<?php include('../../../assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$(document).ready(function(){
			setActiveOption('account-find');
			$("#gen_pass_bttn").click(function(){
				$.ajax({
					type: "GET",
					url: "view.php?do=generate_password",
					success: function(data) {
						$("#gen_pass").html('<strong>Generated Password:</strong> '+data);
						$("#gen_pass").slideDown();
						$('input[name="pass"]').val(data);
						$('input[name="pass_2"]').val(data);
						return false;
					}
				});
				return false;
			});
		});
	</script>
</body>
</html>