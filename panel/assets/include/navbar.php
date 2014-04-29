<div class="pull-right" style="margin-top: -26px;">
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=de" class="language">Deutsch</a> 
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=en" class="language">English</a> 
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=es" class="language">Espa&ntilde;ol</a> 
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=fr" class="language">Fran&ccedil;ais</a> 
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=pt" class="language">Portugu&ecirc;s</a> 
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=ru" class="language">&#1088;&#1091;&#1089;&#1089;&#1082;&#1080;&#1081;</a> 
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=se" class="language">Svenska</a> 
	<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=zh" class="language">&#20013;&#22269;&#30340;的</a> 
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