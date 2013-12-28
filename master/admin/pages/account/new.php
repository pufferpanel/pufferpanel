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
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>PufferPanel - Create New Account</title>
	
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
						<h3 class="fl">Create New User Account</h3>
					</div>
					<div class="content-module-main cf">
					<?php 
						
						if(isset($_GET['disp']) && !empty($_GET['disp'])){
						
							switch($_GET['disp']){
							
								case 'u_fail':
									echo '<div class="error-box">The username you entered does not meet the requirements. Must be at least 4 characters, and no more than 35. Username can only contain a-zA-Z0-9_-</div>';
									break;
								case 'e_fail':
									echo '<div class="error-box">The email you entered is invalid.</div>';
									break;
								case 'p_fail':
									echo '<div class="error-box">The passwords you entered did not match or were not at least 8 characters.</div>';
									break;
								case 'a_fail':
									echo '<div class="error-box">Account with that username or email already exists in the system.</div>';
									break;
							
							}
						
						}
					
					?>
						<fieldset>
							<form action="ajax/new/create.php" method="POST">
								<p>
									<label for="username">Username</label>
									<input type="text" autocomplete="off" name="username" class="round default-width-input" />
									<em>Character Limits: a-zA-Z0-9_- (Max 35 characters)</em>
								</p>
								<p>
									<label for="email">Email</label>
									<input type="text" autocomplete="off" name="email" class="round default-width-input" />
								</p>
								<div class="stripe-separator"><!--  --></div>
								<div class="warning-box round" style="display: none;" id="gen_pass"></div>
								<p>
									<label for="pass">Password (<a href="#" id="gen_pass_bttn">Generate</a>)</label>
									<input type="password" autocomplete="off" name="pass" class="round default-width-input" />
									<em>Minimum Length 8 characters. Suggested 12.</em>
								</p>
								<p>
									<label for="pass_2">Password (Again)</label>
									<input type="password" autocomplete="off" name="pass_2" class="round default-width-input" />
								</p>
								<div class="stripe-separator"><!--  --></div>
								<p><em>To add a server to this user please go to the add new server page.</em></p>
								<input type="submit" value="Create User" class="round blue ic-right-arrow" />
							</form>
						</fieldset>
					</div>
				</div>
			</div>
		</div>
	</div>
	<script type="text/javascript">
		$("#gen_pass_bttn").click(function(){
			$.ajax({
				type: "GET",
				url: "new.php?do=generate_password",
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
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.4.2 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>
