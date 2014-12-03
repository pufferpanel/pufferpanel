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
if(isset($_POST['version']) && file_exists('do/'.$_POST['version'].'.php'))
	header('Location: do/'.$_POST['version'].'.php');
?>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<link rel="stylesheet" href="../../assets/css/bootstrap.css">
	<title>PufferPanel Upgrade Script</title>
</head>
<body>
	<div class="container">
		<div class="navbar navbar-default">
			<div class="navbar-header">
			</div>
		</div>
		<div class="row">
			<div class="col-10">
				<h1 class="nopad">PufferPanel Upgrader</h1>
				<p>Please select the option below that reflects which verison of PufferPanel you are upgrading from. This will run you through the upgrade process.</p>
				<form action="index.php" method="POST">
					<div>
						<label><input type="checkbox" name="version" disabled="disabled" value="0.7.2-beta"/> 0.7.2 Beta [Disabled]</label><br />
					</div>
					<div>
						<label><input type="checkbox" name="version" disabled="disabled" value="0.7.2-beta-bugfix"/> 0.7.2 Beta (Bugfix) [Disabled]</label><br />
					</div>
					<div>
						<label><input type="checkbox" name="version" value="0.7.3-beta"/> 0.7.3 Beta</label><br />
					</div>
					<div>
						<label><input type="checkbox" name="version" value="0.7.4-beta"/> 0.7.4 Beta</label><br />
					</div>
					<div>
						<label><input type="checkbox" name="version" value="0.7.4.1-beta"/> 0.7.4.1 Beta</label><br />
					</div>
					<div>
						<input type="submit" name="do" value="Run Upgrader" class="btn btn-sm btn-primary" />
					</div>
				</form>
			</div>
		</div>
		<div class="footer">
			<div class="row">
				<div class="col-8">
					<p>PufferPanel is licensed under a <a href="https://github.com/PufferPanel/PufferPanel/blob/master/LICENSE">GPL-v3 License</a>.<br />
						Running <?php echo trim(file_get_contents('../../../src/versions/current')).' ('.substr(trim(@file_get_contents('../../../.git/HEAD')), 0, 8).')'; ?>
						distributed by <a href="http://pufferpanel.com">PufferPanel Development</a>.
					</p>
				</div>
			</div>
		</div>
	</div>
</body>
</html>