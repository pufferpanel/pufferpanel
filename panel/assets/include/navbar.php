<p class="pull-right" style="margin-top: -26px;">
	<img src="http://localhost/PufferPan/panel/assets/flags/da.png" />
	<img src="http://localhost/PufferPan/panel/assets/flags/de.png" />
	<img src="http://localhost/PufferPan/panel/assets/flags/en.png" />
	<img src="http://localhost/PufferPan/panel/assets/flags/es.png" />
	<img src="http://localhost/PufferPan/panel/assets/flags/fr.png" />
	<img src="http://localhost/PufferPan/panel/assets/flags/ja.png" />
	<img src="http://localhost/PufferPan/panel/assets/flags/nl.png" />
	<img src="http://localhost/PufferPan/panel/assets/flags/pt.png" />
	<img src="http://localhost/PufferPan/panel/assets/flags/se.png" />
	<img src="http://localhost/PufferPan/panel/assets/flags/zh.png" />
</p>
<div class="navbar navbar-default">
	<div class="navbar-header">
		<a class="navbar-brand" href="#"><?php echo $core->settings->get('company_name'); ?></a>
	</div>
	<div class="navbar-collapse navbar-responsive-collapse collapse" style="height: 1px;">
		<ul class="nav navbar-nav navbar-right">
			<li class="dropdown">
				<a href="#" class="dropdown-toggle" data-toggle="dropdown"><?php echo $_l->tpl('header.account'); ?> <b class="caret"></b></a>
					<ul class="dropdown-menu">
						<li><a href="/logout.php"><?php echo $_l->tpl('header.logout'); ?></a></li>
						<?php if($core->user->getData('root_admin') == 1){ echo '<li><a href="admin/index.php">'.$_l->tpl('header.admin').'</a></li>'; } ?>
					</ul>
			</li>
		</ul>
	</div>
</div>
