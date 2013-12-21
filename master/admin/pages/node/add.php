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
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>PufferPanel - Add New Node</title>
	
	<!-- Stylesheets -->
	<link href='http://fonts.googleapis.com/css?family=Droid+Sans:400,700' rel='stylesheet'>
	<link rel="stylesheet" href="../../../assets/css/style.css">
	
	<!-- Optimize for mobile devices -->
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
	
	<!-- jQuery & JS files -->
	<script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
	<script src="../../../assets/javascript/jquery.cookie.js"></script>
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
						<h3 class="fl">Add New Node</h3>
					</div>
					<div class="content-module-main cf">
					<?php 
						
						if(isset($_GET['disp']) && !empty($_GET['disp'])){
						
							switch($_GET['disp']){
								
								case 'agree_warn':
									echo '<div class="error-box">You must agree to the node warning before we can create the node.</div>';
									break;
								case 'missing_args':
									echo '<div class="error-box">Not all arguments were passed by the script.</div>';
									break;
								case 'n_fail':
									echo '<div class="error-box">The node name does not meet the requirements (1-15 characters, a-zA-Z0-9_.-).</div>';
									break;
								case 'url_fail':
									echo '<div class="error-box">The node URL provided is not valid.</div>';
									break;
								case 'ip_fail':
									echo '<div class="error-box">The IP addresses provided were not valid.</div>';
									break;
								case 'dir_fail':
									echo '<div class="error-box">The director(y/ies) you entered were not valid. They must end with a trailing slash.</div>';
									break;
								case 'fir_match_fail':
									echo '<div class="error-box">The main server directory and backup directory cannot be the same, and the backup directory cannot be located inside the main server directory.</div>';
									break;
								case 'user_fail':
									echo '<div class="error-box">SSH users must not be blank, and may not be \'root\'.</div>';
									break;
								case 'pass_fail':
									echo '<div class="error-box">SSH passwords must be at least 12 characters.</div>';
									break;
							
							}
						
						}
					
					?>
						<fieldset>
							<form action="ajax/new/create.php" method="POST">
								<p>
									<label for="node_name">Node Short Name</label>
									<input type="text" autocomplete="off" name="node_name" placeholder="1.nyc.us" class="round default-width-input" />
									<em>15 character maximum (a-zA-Z0-9_-.)
								</p>
								<p>
									<label for="node_url">Node URL</label>
									<input type="text" autocomplete="off" name="node_url" placeholder="http://1.nyc.us.example.com" class="round default-width-input" />
								</p>
								<div class="stripe-separator"><!--  --></div>
								<p>
									<label for="node_ip">Node IP Address</label>
									<input type="text" autocomplete="off" name="node_ip" class="round default-width-input" />
								</p>
								<p>
									<label for="node_sftp_ip">Node SFTP IP Address</label>
									<input type="text" autocomplete="off" name="node_sftp_ip" class="round default-width-input" />
									<em>In most cases this is the same as the Node IP Address</em>
								</p>
								<p>
									<label for="s_dir">Server Directory</label>
									<input type="text" autocomplete="off" name="s_dir" placeholder="/srv/servers/" class="round default-width-input" />
								</p>
								<p>
									<label for="s_dir_backup">Server Backup Directory</label>
									<input type="text" autocomplete="off" name="s_dir_backup" placeholder="/second/backups/" class="round default-width-input" />
									<em>Remote backup directories are not <strong>currently</strong> supported.</em>
								</p>
								<div class="stripe-separator"><!--  --></div>
								<div class="warning-box round" style="display: none;" id="gen_pass"></div>
								<p>
									<label for="ssh_user">SSH Username</label>
									<input type="text" autocomplete="off" name="ssh_user" class="round default-width-input" />
								</p>
								<p>
									<label for="ssh_pass">SSH Password</label>
									<input type="password" autocomplete="off" name="ssh_pass" class="round default-width-input" />
								</p>
								<div class="stripe-separator"><!--  --></div>
								<p>
									<label for="ip_port">Available IPs &amp; Ports</label>
									<textarea name="ip_port" class="round full-width-input" rows="10" placeholder="127.0.0.1|25565,25566,25567,25568,25569,25570"></textarea>
									<em>Enter one IP address per line, followed by a pipe (|) and then a list of each available port separated with commas.</em>
								</p>
								<div class="stripe-separator"><!--  --></div>
								<div class="warning-box round"><input type="checkbox" name="read_warning" /> By checking this box you are confirming that you have correctly set up your node to handle Minecraft&trade; servers created from this system. Do not add this node until you have correctly done so. Please consult the <a href="#documentation-404">documentation</a> for how to do this if you are unsure.</div>
								<input type="submit" value="Create Node" class="round blue ic-right-arrow" />
							</form>
						</fieldset>
					</div>
				</div>
			</div>
		</div>
	</div>
	<script type="text/javascript">
		$.urlParam = function(name){
		    var results = new RegExp('[\\?&]' + name + '=([^&#]*)').exec(decodeURIComponent(window.location.href));
		    if (results==null){
		       return null;
		    }
		    else{
		       return results[1] || 0;
		    }
		}
		if($.urlParam('error') != null){
		
			var field = $.urlParam('error');
			var exploded = field.split('|');
			
				$.each(exploded, function(key, value) {
					
					$('[name="' + value + '"]').addClass('error-input');
					
				});
				
			var obj = $.parseJSON($.cookie('__TMP_pp_admin_newnode'));
			
				$.each(obj, function(key, value) {
					
					$('[name="' + key + '"]').val(value);
					
				});
		
		}
	</script>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.4 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>