<?php
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
						<h3 class="fl">Add New Node</h3>
					</div>
					<div class="content-module-main cf">
					<?php 
						
						if(isset($_GET['disp']) && !empty($_GET['disp'])){
						
							switch($_GET['disp']){
							
							
							}
						
						}
					
					?>
						<fieldset>
							<form action="ajax/new/create.php" method="POST">
								<p>
									<label for="node_name">Node Short Name</label>
									<input type="text" autocomplete="off" name="node_name" placeholder="nyc-1" class="round default-width-input" />
								</p>
								<p>
									<label for="node_url">Node URL</label>
									<input type="text" autocomplete="off" name="node_url" placeholder="http://nyc-1.example.com" class="round default-width-input" />
								</p>
								<div class="stripe-separator"><!--  --></div>
								<p>
									<label for="node_ip">Node IP Address</label>
									<input type="text" autocomplete="off" name="node_ip" class="round default-width-input" />
								</p>
								<p>
									<label for="node_sftp_ip">Node SFTP IP Address</label>
									<input type="text" autocomplete="off" name="node_sftp_ip" class="round default-width-input" />
								</p>
								<p>
									<label for="s_dir">Server Directory</label>
									<input type="text" autocomplete="off" name="s_dir" placeholder="/srv/servers/" class="round default-width-input" />
								</p>
								<p>
									<label for="s_dir_backup">Server Backup Directory</label>
									<input type="text" autocomplete="off" name="s_dir_backup" placeholder="/second/backups/" class="round default-width-input" />
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
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.3 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>