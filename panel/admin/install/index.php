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
if(file_exists('install.lock'))
	exit('Installer is Locked.');
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<link rel="stylesheet" href="../../assets/css/bootstrap.css">
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
					<div class="panel panel-info">
						<div class="panel-heading">
							<h3 class="panel-title">Installer Information</h3>
						</div>
						<div class="panel-body">
							<p>When this installer finishes please manually change the permissions back to <code>0755</code> for the <code>/src/framework</code> folder, and delete this installer. Please set the permissions on the <code>configuration.php</code> file to <code>0444</code>.</p>
							<p>You are installing version <code>git-<?php echo trim(@file_get_contents('../../../.git/HEAD')); ?></code>. Please keep track of this information as we may request it when you report bugs.</p>
						</div>
					</div>
					<?php

						/* Permissions Checking */
						$successList = null; $failedList = null;
						$continue = true;

						/*
						 * Fail if not installed with git
						 */
						if(!@file_get_contents('../../../src/versions/current')){

							echo '<div class="panel panel-danger">
								<div class="panel-heading">
									<h3 class="panel-title">Incompatable Install Method</h3>
								</div>
								<div class="panel-body">
									<p class="text-danger">This panel <strong>must</strong> be installed using `git clone` on your server. Please re-read the wiki and follow the directions correctly.</p>
								</div>
							</div><hr />';
							$continue = false;

						}

						/*
						* Fail if not composer hasn't been run
						*/
						if(!@file_get_contents('../../../vendor/autoload.php')){

							echo '<div class="panel panel-danger">
								<div class="panel-heading">
									<h3 class="panel-title">Run Composer</h3>
								</div>
								<div class="panel-body">
									<p class="text-danger">This panel <strong>must</strong> has composer run before being installed. Please double check the wiki for instructions on doing this.</p>
								</div>
							</div><hr />';
							$continue = false;

						}

						/*
						 * Check to make sure PHP is at least 5.5.0
						 */
						if(version_compare(PHP_VERSION, "5.5.0") < 0){

							echo '<div class="panel panel-danger">
								<div class="panel-heading">
									<h3 class="panel-title">PHP version is not compatible</h3>
								</div>
								<div class="panel-body">
									<p class="text-danger">Minimum Required Version: <code>5.5.0</code><br />
									Currently Installed: <code>'.PHP_VERSION.'</code></p>
								</div>
							</div><hr />';
							$continue = false;

						}

						/*
						 * Check Configuration File
						 */
						if(substr(sprintf('%o', fileperms('../../../src/framework/configuration.php.dist')), -4) == "0666")
							$successList .= '<p class="text-success"><code>/src/framework/configuration.php.dist</code> is correctly CHMOD\'d.</p>';
						else
							$failedList .= '<p class="text-danger"><code>/src/framework/configuration.php.dist</code> is improperly CHMOD\'d. It should be 0666.</p>';


						/*
						 * Check Framework Folder
						 */
						if(substr(sprintf('%o', fileperms('../../../src/framework')), -4) == "0777")
							$successList .= '<p class="text-success"><code>/src/framework</code> is correctly CHMOD\'d.</p>';
						else
							$failedList .= '<p class="text-danger"><code>/src/framework</code> is improperly CHMOD\'d. It should be 0777.</p>';

						/*
						 * Check Installer Folder
						 */
						if(substr(sprintf('%o', fileperms('../install')), -4) == "0777")
							$successList .= '<p class="text-success"><code>/panel/admin/install</code> is correctly CHMOD\'d.</p>';
						else
							$failedList .= '<p class="text-danger"><code>/panel/admin/install</code> is improperly CHMOD\'d. It should be 0777.</p>';

						/*
						 * Check Installer Process Folder
						 */
						if(substr(sprintf('%o', fileperms('do')), -4) == "0777")
							$successList .= '<p class="text-success"><code>/panel/admin/install/do</code> is correctly CHMOD\'d.</p><br />';
						else
							$failedList .= '<p class="text-danger"><code>/panel/admin/install/do</code> is improperly CHMOD\'d. It should be 0777.</p>';

						/*
						 * Output
						 */
						if(!is_null($failedList)){

							echo '<div class="panel panel-danger">
								<div class="panel-heading">
									<h3 class="panel-title">Failed CHMOD Checks</h3>
								</div>
								<div class="panel-body">
									'.$failedList.'
								</div>
							</div><hr />';

						}

						if(!is_null($successList)){

							echo '<div class="panel panel-success">
								<div class="panel-heading">
									<h3 class="panel-title">Passed CHMOD Checks</h3>
								</div>
								<div class="panel-body">
									'.$successList.'
								</div>
							</div>';

						}

					?>

					<hr />
					<p><?php

						$successList = null; $failedList = null;

						/* List of Required Dependencies */
						$list = array(
							'curl',
							'hash',
							'openssl',
							'mcrypt',
							'PDO',
							'pdo_mysql'
						);

						/*
						 * Check for the Dependencies
						 */
						foreach($list as $ext) {

							if(extension_loaded($ext))
								$successList .= '<p class="text-success">The php-'.$ext.' extension was loaded.</p>';
							else {
								$failedList .= '<p class="text-danger"><strong>The php-'.$ext.' extension was not loaded.</strong></p>';
								$continue = false;
							}

						}

						/*
						 * Output
						 */
						if(!is_null($failedList)){

							echo '<div class="panel panel-danger">
								<div class="panel-heading">
									<h3 class="panel-title">Failed Dependencies Checks</h3>
								</div>
								<div class="panel-body">
									'.$failedList.'
								</div>
							</div>';

						}

						if(!is_null($successList)){

							echo '<div class="panel panel-success">
								<div class="panel-heading">
									<h3 class="panel-title">Passed Dependencies Checks</h3>
								</div>
								<div class="panel-body">
									'.$successList.'
								</div>
							</div>';

						}

					?></p>
					<hr />
					<p><?php

						$successList = null; $failedList = null;

						/* List of Required Functions */
						$functions = array(
							'fopen',
							'fclose',
							'fwrite',
							'session_start',
							'socket_set_option',
							'socket_send',
							'socket_connect',
							'socket_create',
							'stream_set_timeout',
							'fsockopen',
							'crypt',
							'hash'
						);

						/*
						 * Check for the Functions
						 */
						foreach($functions as $fct) {

							if(function_exists($fct))
								$successList .= '<p class="text-success">'.$fct.'() is enabled.</p>';
							else {
								$failedList .= '<p class="text-danger"><strong>'.$fct.'() is not enabled.</strong></p>';
								$continue = false;
							}

						}

						/*
						 * Output
						 */
						if(!is_null($failedList)){

							echo '<div class="panel panel-danger">
								<div class="panel-heading">
									<h3 class="panel-title">Failed Function Checks</h3>
								</div>
								<div class="panel-body">
									'.$failedList.'
								</div>
							</div>';

						}

						if(!is_null($successList)){

							echo '<div class="panel panel-success">
								<div class="panel-heading">
									<h3 class="panel-title">Passed Function Checks</h3>
								</div>
								<div class="panel-body">
									'.$successList.'
								</div>
							</div>';

						}

					?></p>
					<hr />
					<?php echo ($continue === true) ? '<a href="do/start.php">Continue &rarr;</a>' : '<div class="alert alert-info">Please fix missing extensions and functions before continuing.</div>'; ?>
				</div>
				<div class="col-2"></div>
			</div>
		</div>
		<div class="footer">
			<div class="col-8 nopad"><p>PufferPanel is licensed under a <a href="https://github.com/DaneEveritt/PufferPanel/blob/master/LICENSE">GPL-v3 License</a>.<br />Running <?php echo trim(file_get_contents('../../../src/versions/current')).' ('.substr(trim(@file_get_contents('../../../.git/HEAD')), 0, 8).')'; ?> distributed by <a href="http://pufferpanel.com">PufferPanel Development</a>.</p></div>
		</div>
	</div>
</body>
</html>
