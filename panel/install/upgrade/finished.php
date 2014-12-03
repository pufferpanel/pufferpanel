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
				<h1 class="nopad">PufferPanel Upgrader has Finished</h1>
				<p>Your panel has been successfully upgraded to the latest version of PufferPanel.</p>
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