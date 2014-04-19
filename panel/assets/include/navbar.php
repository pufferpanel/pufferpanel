<div class="pull-right" style="margin-top: -26px;">
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=ar"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/ar.png" /></a>
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=da"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/da.png" /></a>
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=de"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/de.png" /></a>
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=en"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/en.png" /></a>
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=es"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/es.png" /></a>
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=fr"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/fr.png" /></a>
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=ja"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/ja.png" /></a>
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=nl"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/nl.png" /></a>
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=pt"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/pt.png" /></a>
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=se"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/se.png" /></a>
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=zh"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/zh.png" /></a>
</div>
<div class="navbar navbar-default">
	<div class="navbar-header">
		<a class="navbar-brand" href="#"><?php echo $core->settings->get('company_name'); ?></a>
	</div>
	<div class="navbar-collapse navbar-responsive-collapse collapse" style="height: 1px;">
		<ul class="nav navbar-nav navbar-right">
			<li class="dropdown">
				<a href="#" class="dropdown-toggle" data-toggle="dropdown"><?php echo $_l->tpl('header.account'); ?> <b class="caret"></b></a>
					<ul class="dropdown-menu">
						<li><a href="<?php echo $core->settings->get('master_url'); ?>logout.php"><?php echo $_l->tpl('header.logout'); ?></a></li>
						<?php if($core->user->getData('root_admin') == 1){ echo '<li><a href="'.$core->settings->get('master_url').'admin/index.php">'.$_l->tpl('header.admin').'</a></li>'; } ?>
					</ul>
			</li>
		</ul>
	</div>
</div>