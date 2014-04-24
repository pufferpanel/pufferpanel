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
require_once('../../core/framework/framework.core.php');

$filesIncluded = true;

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false)
	Page\components::redirect($core->settings->get('master_url').'index.php?login');

?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('../../assets/include/header.php'); ?>
	<title>PufferPanel - Manage Your Server</title>
</head>
<body>
	<div class="container">
		<?php include('../../assets/include/navbar.php'); ?>
		<div class="row">
			<div class="col-3">
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.acc_actions'); ?></strong></a>
					<a href="../../account.php" class="list-group-item"><?php echo $_l->tpl('sidebar.settings'); ?></a>
					<a href="../../servers.php" class="list-group-item"><?php echo $_l->tpl('sidebar.servers'); ?></a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.server_acc'); ?></strong></a>
					<a href="../index.php" class="list-group-item active"><?php echo $_l->tpl('sidebar.overview'); ?></a>
					<a href="../console.php" class="list-group-item"><?php echo $_l->tpl('sidebar.console'); ?></a>
					<a href="../files/index.php" class="list-group-item"><?php echo $_l->tpl('sidebar.files'); ?></a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.server_sett'); ?></strong></a>
					
					<a href="../settings.php" class="list-group-item"><?php echo $_l->tpl('sidebar.manage'); ?></a>
					<a href="index.php" class="list-group-item active"><?php echo $_l->tpl('sidebar.plugins'); ?></a>
				</div>
			</div>
			<div class="col-9">
				<div class="panel panel-primary">
					<div class="panel-heading">
						<h3 class="panel-title"><?php echo $_l->tpl('node.plugins.welcome_h1'); ?></h3>
					</div>
					<div class="panel-body">
						<p><?php echo $_l->tpl('node.plugins.welcome'); ?></p>
					</div>
				</div>
				<form action="search.php" method="get">
					<legend><?php echo $_l->tpl('node.plugins.search_h1'); ?></legend>
					<fieldset>
						<div class="form-group">
							<label for="email" class="control-label"><?php echo $_l->tpl('node.plugins.plugin_name'); ?></label>
							<div>
								<input type="text" class="form-control" name="term" placeholder="e.g. worldedit" />
							</div>
						</div>
						<div class="form-group">
							<div>
								<input type="submit" class="btn btn-primary" value="<?php echo $_l->tpl('string.search'); ?>" />
								<input type="reset" class="btn btn-default" value="<?php echo $_l->tpl('string.clear'); ?>" />
							</div>
						</div>
					</fieldset>
				</form>
			</div>
		</div>
		<div class="footer">
			<?php include('../../assets/include/footer.php'); ?>
		</div>
	</div>
</body>
</html>