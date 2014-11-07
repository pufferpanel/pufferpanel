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
				<a class="navbar-brand" href="#">PufferPanel Upgrader</a>
			</div>
		</div>
		<div class="row">
			<div class="col-3">
			</div>
		</div>
		<div class="footer">
			<div class="row">
				<div class="col-8 nopad">
					<p>PufferPanel is licensed under a <a href="https://github.com/DaneEveritt/PufferPanel/blob/master/LICENSE">GPL-v3 License</a>.<br />
						Running <?php echo trim(file_get_contents('../../../src/versions/current')).' ('.substr(trim(@file_get_contents('../../../.git/HEAD')), 0, 8).')'; ?>
						distributed by <a href="http://pufferpanel.com">PufferPanel Development</a>.
					</p>
				</div>
			</div>
		</div>
	</div>
</body>
</html>