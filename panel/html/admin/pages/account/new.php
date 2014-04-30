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
	Page\components::redirect('../../../index.php?login');
}

if(isset($_GET['do']) && $_GET['do'] == 'generate_password')
	exit($core->auth->keygen(rand(12, 18)));
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
				<h3 class="nopad">Create New Account</h3><hr />
				<?php 
					
					if(isset($_GET['disp']) && !empty($_GET['disp'])){
					
						switch($_GET['disp']){
						
							case 'u_fail':
								echo '<div class="alert alert-danger">The username you entered does not meet the requirements. Must be at least 4 characters, and no more than 35. Username can only contain a-zA-Z0-9_-</div>';
								break;
							case 'e_fail':
								echo '<div class="alert alert-danger">The email you entered is invalid.</div>';
								break;
							case 'p_fail':
								echo '<div class="alert alert-danger">The passwords you entered did not match or were not at least 8 characters.</div>';
								break;
							case 'a_fail':
								echo '<div class="alert alert-danger">Account with that username or email already exists in the system.</div>';
								break;
						
						}
					
					}
				
				?>
				<form action="ajax/new/create.php" method="post">
					<fieldset>
						<div class="form-group">
							<label for="username" class="control-label">Username</label>
							<div>
								<input type="text" autocomplete="off" name="username" class="form-control" />
							</div>
						</div>
						<div class="form-group">
							<label for="email" class="control-label">Email</label>
							<div>
								<input type="text" autocomplete="off" name="email" class="form-control" />
							</div>
						</div>
						<div id="gen_pass" class="alert alert-success" style="display:none;margin-bottom: 10px;"></div>
						<div class="form-group col-6 nopad">
							<label for="pass" class="control-label">Password</label>
							<div>
								<input type="password" name="pass" class="form-control" />
							</div>
						</div>
						<div class="form-group col-6 nopad-right">
							<label for="pass_2" class="control-label">Password Again</label>
							<div>
								<input type="password" name="pass_2" class="form-control" />
							</div>
						</div>
						<div class="form-group">
							<div>
								<button class="btn btn-primary" type="submit">Create Account</button>
								<button class="btn btn-default" id="gen_pass_bttn" type="button">Generate Password</button>
							</div>
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
		$(document).ready(function(){
			$("#gen_pass_bttn").click(function(e){
				e.preventDefault();
				$.ajax({
					type: "GET",
					url: "new.php?do=generate_password",
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