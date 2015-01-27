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
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<style>
	<?php include PANEL_DIR . 'assets/css/bootstrap.css'; ?>
	</style>
	<script src="//ajax.googleapis.com/ajax/libs/jquery/1.11.0/jquery.min.js"></script>
	<title>PufferPanel Installer</title>
</head>
<body>
	<div class="container">
		<div class="alert alert-danger">
			<strong>WARNING:</strong> Do not run this version on a live environment! There are known security holes that we are working on getting patched. This is extremely beta software and this version is to get the features in place while we work on security enhancements.
		</div>
		<div class="navbar navbar-default">
			<div class="navbar-header">
				<a class="navbar-brand" href="#">Install PufferPanel</a>
			</div>
		</div>
		<div class="col-12">
			<div class="row">
				<div class="col-2"></div>
				<div class="col-8">
					<div class="well">
						<p>When this installer finishes please manually change the permissions back to <code>0755</code> for the <code>/src/core</code> folder, and delete this installer. Please set the permissions on the <code>configuration.php</code> file to <code>0444</code>.</p>
					</div>
					<h3>Detailed Version Information</h3>
					<div class="well well-sm">
						<code><?php echo \PufferPanel\Core\Version::get(); ?></code>
					</div>
					<h3>Installer Comments</h3>
					<?php

						$continue = true;
						/*
						 * Fail if not Installed with Git
						 */
						if(!is_dir(BASE_DIR.'.git')) { ?>

							<div class="panel panel-danger">
								<div class="panel-heading">
									<h3 class="panel-title">Incompatible Install Method</h3>
								</div>
								<div class="panel-body">
									<p class="text-danger">This panel <strong>must</strong> be installed using <code>git clone</code> on your server. Please re-read the documentatioin and follow the directions correctly.</p>
								</div>
							</div>
							<?php $continue = false;

						}

						/*
						* Fail if not composer hasn't been run
						*/
						if(!is_dir(BASE_DIR.'vendor')) { ?>

							<div class="panel panel-danger">
								<div class="panel-heading">
									<h3 class="panel-title">Run Composer</h3>
								</div>
								<div class="panel-body">
									<p class="text-danger">This panel <strong>must</strong> have composer run before being installed. Please double check the documentation for instructions on doing this.</p>
								</div>
							</div>';
							<?php $continue = false;

						}

						/*
						 * Check to make sure PHP is at least 5.5.0
						 */
						if(version_compare(PHP_VERSION, "5.5.0") < 0) { ?>

							<div class="panel panel-danger">
								<div class="panel-heading">
									<h3 class="panel-title">PHP Version too Low</h3>
								</div>
								<div class="panel-body">
									<p class="text-danger">Minimum Required Version: <code>5.5.0</code><br />
									Currently Installed: <code><?php echo PHP_VERSION; ?></code></p>
								</div>
							</div>
							<?php $continue = false;

						}

						$failedchmod = array();

						/*
						 * Check Folder CHMOD Permissions
						 */
						if(substr(sprintf('%o', fileperms(SRC_DIR.'core')), -4) != "0777") {
							$failedchmod[] = '<p class="text-danger"><code>/src/core</code> is improperly CHMOD\'d. It should be 0777.</p>';
						}

						if(substr(sprintf('%o', fileperms(SRC_DIR.'cache')), -4) != "0777") {
							$failedchmod[] = '<p class="text-danger"><code>/src/cache</code> is improperly CHMOD\'d. It should be 0777.</p>';
						}

						if(!empty($failedchmod)){ ?>

							<div class="panel panel-danger">
								<div class="panel-heading">
									<h3 class="panel-title">Failed CHMOD Checks</h3>
								</div>
								<div class="panel-body">
									<?php echo implode("", $failedchmod); ?>
								</div>
							</div>

							<?php $continue = false;

						}

						/* Check for Required Dependencies */
						$failedextensions = array();
						$list = array('curl', 'hash', 'openssl', 'mcrypt', 'PDO', 'pdo_mysql');
						foreach($list as $extension) {
							if(!extension_loaded($extension)) {
								$failedextensions[] = '<p class="text-danger">The <code>php-'.$extension.'</code> extension was not able to be loaded.</p>';
							}
						}

						if(!empty($failedextensions)){ ?>

							<div class="panel panel-danger">
								<div class="panel-heading">
									<h3 class="panel-title">Failed Dependencies Checks</h3>
								</div>
								<div class="panel-body">
									<?php echo implode("", $failedextensions); ?>
								</div>
							</div>

							<?php $continue = false;

						}

						/* Check for Required Functions */
						$failedList = array();
						$functions = array('fopen', 'fclose', 'fwrite', 'session_start', 'socket_set_option', 'socket_send', 'socket_connect', 'socket_create', 'stream_set_timeout', 'fsockopen', 'crypt', 'hash', 'curl_init', 'curl_setopt', 'curl_exec', 'curl_close');
						foreach($functions as $function) {
							if(!function_exists($function)) {
								$failedList[] = '<p class="text-danger"><code>'.$function.'()</code> is not enabled.</p>';
							}
						}

						if(!empty($failedList)) { ?>

							<div class="panel panel-danger">
								<div class="panel-heading">
									<h3 class="panel-title">Failed Function Checks</h3>
								</div>
								<div class="panel-body">
									<?php echo implode("", $failedList); ?>
								</div>
							</div>

							<?php $continue = false;

						}
						if($continue) { ?>
							<div class="well"><p style="margin-bottom:0;">Everything looks good here captain!</p></div><a onclick="cont()" href="#">Continue &rarr;</a>
						<?php } else { ?>
							<div class="alert alert-info">Please fix the errors above before continuing.</div>
						<?php }
					?>
				</div>
				<div class="col-2"></div>
			</div>
		</div>
		<div class="footer">
			<div class="col-8 nopad"><p>PufferPanel is licensed under a <a href="https://github.com/PufferPanel/PufferPanel/blob/master/LICENSE">GPL-v3 License</a>.<br />Running <?php echo \PufferPanel\Core\Version::get(); ?> distributed by <a href="http://pufferpanel.com">PufferPanel Development</a>.</p></div>
		</div>
	</div>
</body>

<script>
	function cont() {
		$.post('install', function() {
			window.location.href = "install/start";
		});		
	}
</script>
</html>
