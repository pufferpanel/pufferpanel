<?php
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
	<title>PufferPanel - Create New Server</title>
	
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
				<li><a href="../../../account.php" class="round button dark"><i class="icon-user"></i>&nbsp;&nbsp; <strong><?php echo $core->framework->user->getData('username'); ?></strong></a></li>
			</ul>
			<ul id="nav" class="fr">
				<li><a href="../../../servers.php" class="round button dark"><i class="icon-signout"></i></a></li>
				<li><a href="../../../logout.php" class="round button dark"><i class="icon-off"></i></a></li>
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
						<h3 class="fl">Create New Server</h3>
					</div>
					<div class="content-module-main cf">
					<?php 
						
						if(isset($_GET['disp']) && !empty($_GET['disp'])){
						
							switch($_GET['disp']){
							
								case 's_fail':
									echo '<div class="error-box">The server name you entered does not meet the requirements. Must be at least 4 characters, and no more than 35. Server name can only contain a-zA-Z0-9_-</div>';
									break;
								case 'n_fail':
									echo '<div class="error-box">The node selected does not seem to exist.</div>';
									break;
								case 'ip_fail':
									echo '<div class="error-box">The selected IP does not exist.</div>';
									break;
								case 'port_fail':
									echo '<div class="error-box">The selected port does not exist.</div>';
									break;
								case 'port_full':
									echo '<div class="error-box">The selected port is already in use.</div>';
									break;
								case 'e_fail':
									echo '<div class="error-box">The email you entered is invalid.</div>';
									break;
								case 'p_fail':
									echo '<div class="error-box">The passwords you entered did not match or were not at least 8 characters.</div>';
									break;
								case 'a_fail':
									echo '<div class="error-box">Account with that email does not exist in the system.</div>';
									break;
								case 'm_fail':
									echo '<div class="error-box">You entered a non-number for Disk and/or Memory.</div>';
									break;
								case 'b_fail':
									echo '<div class="error-box">You entered a non-number for Backup Files and/or Disk Space.</div>';
									break;
							
							}
						
						}
					
					?>
						<fieldset>
							<form action="ajax/new/create.php" method="POST">
								<p>
									<label for="server_name">Server Name</label>
									<input type="text" autocomplete="off" name="server_name" class="round default-width-input" />
									<em>Character Limits: a-zA-Z0-9_- (Max 35 characters)</em>
								</p>
								<p>
									<label for="node">Node</label>
									<select name="node" id="getNode" class="round default-width-input">
										<?php
											$selectData = $mysql->prepare("SELECT * FROM `nodes`");
											$selectData->execute(array());
											while($node = $selectData->fetch()){
											
												echo '<option value="'.$node['node'].'">'.$node['node'].'</option>';
											
											}
										?>
									</select><i class="icon-angle-down select-arrow"></i>
								</p>
								<p>
									<label for="email">Owner Email</label>
									<input type="text" autocomplete="off" name="email" value="<?php if(isset($_GET['email'])) echo $_GET['email']; ?>" class="round default-width-input" />
								</p>
								<div class="stripe-separator"><!--  --></div>
								<p><em><a href="ajax/new/viewPopup.php" onclick="window.open(this.href + '?node=' + $('#getNode').val(), 'View Avaliable IPs and Ports', 'width=500, height=600, left=24, top=24'); return false;">View available</a> IPs &amp; Ports as well as free RAM and Disk Space.</em></p>
								<p>
									<label for="email">Assign IP Address</label>
									<input type="text" autocomplete="off" name="server_ip" class="round default-width-input" />
								</p>
								<p>
									<label for="email">Assign Port</label>
									<input type="text" autocomplete="off" name="server_port" class="round default-width-input" />
								</p>
								<p>
									<label for="email">Allocate Memory (in MB)</label>
									<input type="text" autocomplete="off" name="alloc_mem" class="round default-width-input" />
								</p>
								<p>
									<label for="email">Allocate Disk Space (in MB)</label>
									<input type="text" autocomplete="off" name="alloc_disk" class="round default-width-input" />
								</p>
								<div class="stripe-separator"><!--  --></div>
								<div class="warning-box round" style="display: none;" id="gen_pass"></div>
								<p>
									<label for="pass">SFTP Password (<a href="#" id="gen_pass_bttn">Generate</a>)</label>
									<input type="password" autocomplete="off" name="sftp_pass" class="round default-width-input" />
									<em>Minimum Length 8 characters. Suggested 12.</em>
								</p>
								<p>
									<label for="pass_2">SFTP Password (Again)</label>
									<input type="password" autocomplete="off" name="sftp_pass_2" class="round default-width-input" />
								</p>
								<div class="stripe-separator"><!--  --></div>
								<p>
									<label for="backup_disk">Backup Disk Space (in MB)</label>
									<input type="text" autocomplete="off" name="backup_disk" class="round default-width-input" />
								</p>
								<p>
									<label for="backup_files">Backup Max Files</label>
									<input type="text" autocomplete="off" name="backup_files" class="round default-width-input" />
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
				url: "add.php?do=generate_password",
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
		$.urlParam = function(name){
		    var results = new RegExp('[\\?&]' + name + '=([^&#]*)').exec(window.location.href);
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
			$("#errorDetected").slideDown();
			
				$.each(exploded, function(key, val) {
					
					$('[name="' + val + '"]').addClass('error-input');
					
				});
		
		}
	</script>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.3 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>
