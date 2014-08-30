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
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('../../../../src/include/header.php'); ?>
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
			<div class="col-3"><?php include('../../../../src/include/admin.php'); ?></div>
			<div class="col-9">
				<h3 class="nopad">Add New Node</h3><hr />
				<?php

					if(isset($_GET['disp']) && !empty($_GET['disp'])){

						switch($_GET['disp']){

							case 'agree_warn':
								echo '<div class="alert alert-danger">You must agree to the node warning before we can create the node.</div>';
								break;
							case 'missing_args':
								echo '<div class="alert alert-danger">Not all arguments were passed by the script.</div>';
								break;
							case 'n_fail':
								echo '<div class="alert alert-danger">The node name does not meet the requirements (1-15 characters, a-zA-Z0-9_.-).</div>';
								break;
							case 'url_fail':
								echo '<div class="alert alert-danger">The node URL provided is not valid. URLs must end with a trailing slash and must be a subdomain without any additional folders. (e.g. <strong>http://node.example.com/</strong>)</div>';
								break;
							case 'ip_fail':
								echo '<div class="alert alert-danger">The IP addresses provided were not valid.</div>';
								break;
							case 'dir_fail':
								echo '<div class="alert alert-danger">The directories you entered were not valid. They must end with a trailing slash.</div>';
								break;
							case 'user_fail':
								echo '<div class="alert alert-danger">SSH users must not be blank, and may not be \'root\'.</div>';
								break;
							case 'ip_port_space':
								echo '<div class="alert alert-danger">Unable to allocate the ports you inputted.</div>';
								break;

						}

					}

				?>
				<form action="ajax/new/create.php" method="POST">
					<fieldset>
						<div class="well">
							<div class="row">
								<div class="form-group col-6 nopad">
									<label for="node_name" class="control-label">Node Short Name</label>
									<div>
										<input type="text" name="node_name" placeholder="shortname" class="form-control" />
										<p class="text-muted" style="margin: 0 0 -10.5px;"><small><em>15 character maximum (a-zA-Z0-9_-.).</em></small></p>
									</div>
								</div>
								<div class="form-group col-6 nopad-right">

								</div>
							</div>
							<div class="row">
								<div class="form-group col-6 nopad">
									<label for="node_ip" class="control-label">Node IP Address</label>
									<div>
										<input type="text" name="node_ip" class="form-control" />
									</div>
								</div>
								<div class="form-group col-6 nopad-right">
									<label for="node_sftp_ip" class="control-label">Node SFTP IP Address</label>
									<div>
										<input type="text" name="node_sftp_ip" class="form-control" />
										<p class="text-muted" style="margin: 0 0 -10.5px;"><small><em>In most cases this is the same as the Node IP Address.</em></small></p>
									</div>
								</div>
							</div>
						</div>
						<div class="well">
							<div class="row">
								<div class="form-group">
									<label for="ip_port" class="control-label">Available IPs &amp; Ports</label>
									<div>
										<textarea name="ip_port" class="form-control" rows="5" placeholder="127.0.0.1|25565,25567,25569,25571,25573,25575"></textarea>
										<p class="text-muted" style="margin: 0 0 -10.5px;"><small><em>Enter one IP address per line, followed by a pipe (|) and then a list of each available port separated with commas.</em></small></p>
									</div>
								</div>
							</div>
						</div>
						<div class="alert alert-danger"><input type="checkbox" name="read_warning" /> By checking this box you are confirming that you have correctly set up your node to handle Minecraft&trade; servers created from this system. Do not add this node until you have correctly done so. Please consult the <a href="https://github.com/DaneEveritt/PufferPanel/wiki/Setting-up-a-New-Node" target="_blank">documentation</a> for how to do this if you are unsure.</div>
						<input type="submit" value="Create Node" id="disable_complete" class="btn btn-primary btn-sm disabled" />
					</fieldset>
				</form>
			</div>
		</div>
		<div class="footer">
			<?php include('../../../../src/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$(document).ready(function(){
			$("input[name='node_name']").keyup(function(){
				if($(this).val().length > 15 || /^[a-zA-Z0-9_.-]*$/.test($(this).val()) == false){
					$(this).parent().parent().removeClass('has-success');
					$(this).parent().parent().addClass('has-error');
				}else if($(this).val().length == 0){
					$(this).parent().parent().removeClass('has-success');
					$(this).parent().parent().removeClass('has-error');
				}else{
					$(this).parent().parent().removeClass('has-error');
					$(this).parent().parent().addClass('has-success');
				}
			});
			$("input[name='read_warning']").click(function(){
				if($("input[name='read_warning']").is(":checked"))
					$("#disable_complete").removeClass("disabled");
				else
					$("#disable_complete").addClass("disabled");
			});
			if($.urlParam('error') != null){
				var field = $.urlParam('error');
				var exploded = field.split('|');
					$.each(exploded, function(key, value) {
						$('[name="' + value + '"]').parent().parent().addClass('has-error');
					});
				var obj = $.parseJSON($.cookie('__TMP_pp_admin_newnode'));
					$.each(obj, function(key, value) {
						$('[name="' + key + '"]').val(value);
					});
			}
		});
	</script>
</body>
</html>
