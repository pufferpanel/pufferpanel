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
						<?php if($core->framework->user->getData('root_admin') == 1){ echo '<li><a href="'.$core->framework->settings->get('master_url').'admin/index.php">Admin CP</a></li>'; } ?>
					</ul>
			</li>
		</ul>
	</div>
</div>