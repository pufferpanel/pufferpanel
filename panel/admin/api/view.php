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
				<ul class="nav nav-tabs" id="config_tabs">
					<li class="active"><a href="#view" data-toggle="tab">List Keys</a></li>
					<li><a href="#add" data-toggle="tab">New Key</a></li>
				</ul>
				<div class="tab-content">
					<div class="tab-pane active" id="view">
						<h3>Current API Keys</h3><hr />
						<table class="table table-striped table-bordered table-hover">
							<thead>
								<tr>
									<th>API Key</th>
									<th>Permissions</th>
                                    <th></th>
								</tr>
							</thead>
							<tbody>
								<?php
									$select = $mysql->prepare("SELECT * FROM `api`");
									$select->execute();

									while($row = $select->fetch()){

										echo '<tr>
											<td>
												<code>'.$row['key'].'</code>
											</td>
											<td>
												'.json_encode(json_decode($row['permissions'], true), JSON_PRETTY_PRINT).'
											</td>
                                            <td style="vertical-align: middle;">
                                                <a href="#/trash"><i class="fa fa-trash"></i></a>
                                            </td>
										</tr>';

									}
								?>
							</tbody>
						</table>
					</div>
					<div class="tab-pane" id="add">
						<h3>Add New API Key</h3><hr />
						<form action="ajax/add.php" method="POST">
							<fieldset>
								<div class="form-group">
									<label for="main_url" class="control-label">Allowed Connection IPs</label>
									<div>
										<input type="text" name="allowed_ips" class="form-control" />
										<p><small class="text-muted"><em>Please enter a comma separated list of IPv4 addresses that are allowed to connect to the API to send and receive data. Enter * or leave blank to allow any IP to connect.</em></small></p>
									</div>
								</div>
								<div class="form-group">
									<h3>Permissions</h3><hr />
									<div class="checkbox">
										<label><input type="checkbox" onclick="toggle(this);" /> Check All</label><br />
										<hr />
										<label><input type="checkbox" name="permissions[]" value="list_nodes" /> List Nodes</label><br />
										<label><input type="checkbox" name="permissions[]" value="list_ips" /> List IPs</label><br />
										<label><input type="checkbox" name="permissions[]" value="list_ports" /> List Ports</label><br />
										<label><input type="checkbox" name="permissions[]" value="list_user" /> List User Information</label><br />
										<label><input type="checkbox" name="permissions[]" value="list_server" /> List Server Information</label><br />
										<hr />
										<label><input type="checkbox" name="permissions[]" value="add_server" /> Add New Server</label><br />
										<label><input type="checkbox" name="permissions[]" value="add_user" /> Add New User</label><br />
										<label><input type="checkbox" name="permissions[]" value="add_node" /> Add New Node</label><br />
										<hr />
										<label><input type="checkbox" name="permissions[]" value="update_server" /> Update Server Information</label><br />
										<label><input type="checkbox" name="permissions[]" value="update_sftp_pass" /> Update SFTP Password for Server</label><br />
										<label><input type="checkbox" name="permissions[]" value="update_user_pass" /> Update User Password</label><br />
										<hr />
										<label><input type="checkbox" name="permissions[]" value="delete_user" /> Delete User</label><br />
										<label><input type="checkbox" name="permissions[]" value="delete_server" /> Delete Server</label><br />
									</div>
								</div>
								<div class="form-group">
									<div>
										<input type="submit" value="Add New Key" class="btn btn-primary" />
									</div>
								</div>
							</fieldset>
						</form>
					</div>
				</div>
			</div>
		</div>
		<script type="text/javascript">
		function toggle(source) {
			checkboxes = document.getElementsByName('permissions[]');
				for(var i=0, n=checkboxes.length;i<n;i++) {
					checkboxes[i].checked = source.checked;
			}
		}
		</script>
		<div class="footer">
			<?php include('../../../../src/include/footer.php'); ?>
		</div>
	</div>
</body>
</html>
