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
require_once('../core/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Page\components::redirect('../index.php?login');
}
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('../assets/include/header.php'); ?>
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
			<div class="col-3"><?php include('../assets/include/admin.php'); ?></div>
			<div class="col-9">
				<?php
				    if(is_dir('install'))
				        echo '<div class="alert alert-danger"><strong>WARNING!</strong> Please remove the install/ directory from PufferPanel immediately to prevent any possible security holes.</div>';
				?>
					<h3 class="nopad">PufferPanel Admin Control Panel</h3><hr />
					<p>Welcome to the most advanced, lightweight, and user-friendly control panel for Minecraft. You are currently running version <code><?php echo file_get_contents('../assets/versions/current'); ?></code>.</p>
					
					<p>Please include the following information in any bug reports that you submit:</p>
					<pre><?php
						
						$extensions = "";
						
						foreach(get_loaded_extensions() as $id => $val){
							if(phpversion($val) == "")
								$extensions .= $val.", ";
							else
								$extensions .= $val." (".phpversion($val)."), ";
						}
					
echo "=================[ PufferPanel Debug Output ]=================
\r\r=====[ System Information ]=====
\rOperating System: ".php_uname('s').
"\rOperating System Release: ".php_uname('r').
"\rHost Name: ".php_uname('n').
"\rVersion Information: ".php_uname('v').
"\rMachine Type: ".php_uname('m').
"\r\r=====[ PHP Information ]=====
\rPHP Version: ".phpversion().
"\rPHP SAPI: ".php_sapi_name().
"\rZend Engine Version: ".zend_version().
"\rLoaded PHP Extentions: ".rtrim($extensions, ", ").
"\r\r=====[ Panel Information ]=====
\rPanel Version: ".file_get_contents('../assets/versions/current').
"\rCurrent Directory: ".__DIR__.
"\rCurrent URL: http://".$_SERVER['HTTP_HOST'].$_SERVER['REQUEST_URI'].
"\rMaster URL: ".$core->settings->get('master_url').
"\rAssets URL: ".$core->settings->get('assets_url').
"\rCookie Domain: ".str_replace("_notfound_", "NULL", $core->settings->get('cookie_domain')).
"\rModpack Directory: ".$core->settings->get('modpack_dir')
					?></pre>
			</div>
		</div>
		<div class="footer">
			<?php include('../assets/include/footer.php'); ?>
		</div>
	</div>
</body>
</html>