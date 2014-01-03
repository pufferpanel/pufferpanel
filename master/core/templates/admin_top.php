<div class="navbar navbar-default">
	<div class="navbar-header">
		<a class="navbar-brand" href="<?php echo $core->framework->settings->get('master_url'); ?>account.php"><?php echo $core->framework->user->getData('username'); ?></a>
	</div>
	<div class="navbar-collapse navbar-responsive-collapse collapse" style="height: 1px;">
		<ul class="nav navbar-nav navbar-right">
			<li class="dropdown">
				<a href="#" class="dropdown-toggle" data-toggle="dropdown">Account <b class="caret"></b></a>
					<ul class="dropdown-menu">
						<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>servers.php">Servers</a></li>
						<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>logout.php">Logout</a></li>
					</ul>
			</li>
		</ul>
	</div>
</div>