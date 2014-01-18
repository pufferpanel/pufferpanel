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
				<ul class="nav nav-tabs" id="config_tabs">
					<li class="active"><a href="#installed" data-toggle="tab">Installed Packs</a></li>
					<li><a href="#install" data-toggle="tab">Install New Pack</a></li>
				</ul>
				<div class="tab-content">
					<div class="tab-pane active" id="installed">
						<h3>Installed Modpacks</h3><hr />
							<table class="table table-striped table-bordered table-hover">
								<thead>
									<tr>
										<th>Modpack</th>
										<th>Version</th>
										<th>Added</th>
									</tr>
								</thead>
								<tbody>
									<?php
										
										$packs = $mysql->prepare("SELECT * FROM `modpacks` WHERE `deleted` = 0");
										$packs->execute();
										
										while($row = $packs->fetch()){
											
											$isDefault = ($row['default'] ==  1) ? '<span class="label label-success">Default</span> ' : '';
											echo '<tr>
													<td>'.$isDefault.'<a href="edit.php?mid='.$row['hash'].'">'.$row['name'].'</a></td>
													<td>'.$row['version'].'</td>
													<td>'.date('M d, Y \a\t H:i',$row['added']).'</td>
												</tr>';
										
										}
									
									?>
								</tbody>
							</table>
					</div>
					<div class="tab-pane" id="install">
						<h3>Install New Modpack</h3><hr />
						<?php 
							
							if(isset($_GET['disp']) && !empty($_GET['disp'])){
							
								switch($_GET['disp']){
								
									case 'missing_args':
										echo '<div class="alert alert-danger">Not all arguments were passed by the script.</div>';
										break;
									case 'pn_fail':
										echo '<div class="alert alert-danger">Invalid pack name was provided. Names must be at least one character, and no more than sixty-four characters. Allowed characters: <strong>a-zA-Z0-9._-()</strong> and <strong>[space]</strong></div>';
										break;
									case 'num_fail':
										echo '<div class="alert alert-danger">Minimum RAM and Permgen must be numeric.</div>';
										break;
									case 'ver_fail':
										echo '<div class="alert alert-danger">Invalid pack version was provided. Versions must be at least one character, and no more than thirty-two characters. Allowed characters: <strong>a-zA-Z0-9._-</strong></div>';
										break;
									case 'file_error':
										echo '<div class="alert alert-danger">An unidentified error occured when attempting to process the file upload.</div>';
										break;
									case 'no_file':
										echo '<div class="alert alert-danger">No file was provided to upload.</div>';
										break;
									case 'file_size':
										echo '<div class="alert alert-danger">The file provided exceeded the maximum size of 35MB.</div>';
										break;
									case 'file_type':
										echo '<div class="alert alert-danger">The file type of the upload was invalid. Only <strong>.zip</strong> files are allowed to be uploaded.</div>';
										break;
									case 'file_nomove':
										echo '<div class="alert alert-danger">The file was unable to me moved into the correct directory on the server. Please double check permissions.</div>';
										break;
								
								}
							
							}
						
						?>
						<form action="ajax/modpack/new.php" method="post" enctype="multipart/form-data">
							<fieldset>
								<div class="well">
									<div class="form-group">
										<label for="pack_name" class="control-label">Modpack Name</label>
										<div>
											<input type="text" name="pack_name" class="form-control" />
										</div>
										<div class="checkbox">
											<label>
												<input type="checkbox" name="pack_default" /> Set this as the default modpack.
											</label>
										</div>
									</div>
									<div class="row">
										<div class="form-group col-4 nopad">
											<label for="pack_version" class="control-label">Modpack Version</label>
											<div class="input-group">
												<span class="input-group-addon">v.</span>
												<input type="text" name="pack_version" class="form-control" />
											</div>
										</div>
										<div class="form-group col-4 nopad">
											<label for="pack_minram" class="control-label">Minimum RAM Allocation to Install</label>
											<div class="input-group">
												<input type="text" name="pack_minram" class="form-control" />
												<span class="input-group-addon">MB</span>
											</div>
										</div>
										<div class="form-group col-4 nopad">
											<label for="pack_permgen" class="control-label">PermGen Default Size</label>
											<div class="input-group">
												<input type="text" name="pack_permgen" class="form-control" />
												<span class="input-group-addon">MB</span>
											</div>
										</div>
									</div>
								</div>
								<div class="well">
									<div class="form-group">
										<label for="pack_jar" class="control-label">Upload Modpack</label>
										<div>
											<input type="file" name="pack_jar" class="form-control" />
											<em><p class="text-muted"><small>Please put all of the files for the Modpack into a compressed (.zip) file. Upload the compressed file. The Modpack installer will unzip the file and move them into the server directory when users install the Modpack.</small></p></em>
										</div>
									</div>
								</div>
								<div class="form-group">
									<div>
										<input type="submit" value="Upload Modpack" class="btn btn-primary" />
										<input type="reset" value="Clear Form" class="btn btn-default" />
									</div>
								</div>
							</fieldset>
						</form>
					</div>
				</div>
			</div>
		</div>
		<div class="footer">
			<?php include('../../../assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$(document).ready(function(){
			if($.urlParam('error') != null){
				var field = $.urlParam('error');
				var exploded = field.split('|');
					$.each(exploded, function(key, value) {
						$('[name="' + value + '"]').parent().parent().addClass('has-error');
					});
				var obj = $.parseJSON($.cookie('__TMP_pp_admin_newmodpack'));
					$.each(obj, function(key, value) {
						$('[name="' + key + '"]').val(value);
					});			
			}
			if($.urlParam('tab') != null){
				$('#config_tabs a[href="#'+$.urlParam('tab')+'"]').tab('show');
			}
		});
	</script>
</body>
</html>