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
session_start();
require_once('../../core/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false)
	Page\components::redirect($core->settings->get('master_url').'index.php?login');

$errorMessage = null;
if(isset($_GET['term']) && !empty($_GET['term'])){

	/*
	 * Searching for Plugin
	 */
	((isset($_GET['start']) && $_GET['start'] >= 1) ? $_GET['start'] = $_GET['start'] : $_GET['start'] = '0');
	$_GET['term'] = str_replace(array(' ', '+', '%20'), '', $_GET['term']);
	$data = file_get_contents('http://api.bukget.org/3/search/plugin_name/like/'.$_GET['term'].'?start='.$_GET['start'].'&size=25');
	$data = json_decode($data, true);
	
		$searchResults = '';
		foreach($data as $item => $value){
				
			$value['description'] = (strlen($value['description']) == 0) ? 'No description is avaliable for this plugin.' : $value['description'];
			$value['description'] = (strlen($value['description']) > 200) ? substr($value['description'], 0, 197).'...' : $value['description'];
			
			$searchResults .= '	<tr>
									<td><a href="view.php?slug='.$value['slug'].'">'.$value['plugin_name'].'</a></td>
									<td>'.$value['description'].'</td>
									<td><a href="#install" id="'.$value['slug'].'|0" class="install"><i class="fa fa-download"></i></a></td>
								</tr>';
		
		}
		
	if(count($data) == 25){
		
		if(isset($_GET['start']) && $_GET['start'] > 24){
		
			$pageOptions = '<li><a href="search.php?term='.$_GET['term'].'&start='.($_GET['start'] - 25).'">'.$_l->tpl('string.previous').'</a></li><li><a href="search.php?term='.$_GET['term'].'&start='.($_GET['start'] + 25).'">'.$_l->tpl('string.next').'</a></li>';
		
		}else{
		
			$pageOptions = '<li class="disabled"><a href="#">'.$_l->tpl('string.previous').'</a></li><li><a href="search.php?term='.$_GET['term'].'&start='.($_GET['start'] + 25).'">'.$_l->tpl('string.next').'</a></li>';
		
		}
		
	}else{
		
		if(isset($_GET['start']) && $_GET['start'] != 0){
		
			$pageOptions = '<li><a href="search.php?term='.$_GET['term'].'&start='.($_GET['start'] - 25).'">'.$_l->tpl('string.previous').'</a></li><li class="disabled"><a href="#">'.$_l->tpl('string.next').'</a></li>';
		
		}else{
		
			$pageOptions = null;
		
		}
		
	}

}else{

	$errorMessage = '<div class="alert alert-warning">'.$_l->tpl('node.plugins.error.term').'</div>';
	$searchResults = null;
	$pageOptions = null;

}
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('../../assets/include/header.php'); ?>
	<title>PufferPanel - Manage Your Server</title>
</head>
<body>
	<div class="container">
		<?php include('../../assets/include/navbar.php'); ?>
		<div class="row">
			<div class="col-3">
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.acc_actions'); ?></strong></a>
					<a href="../../account.php" class="list-group-item"><?php echo $_l->tpl('sidebar.settings'); ?></a>
					<a href="../../servers.php" class="list-group-item"><?php echo $_l->tpl('sidebar.servers'); ?></a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.server_acc'); ?></strong></a>
					<a href="../index.php" class="list-group-item"><?php echo $_l->tpl('sidebar.overview'); ?></a>
					<a href="../console.php" class="list-group-item"><?php echo $_l->tpl('sidebar.console'); ?></a>
					<a href="../files/index.php" class="list-group-item"><?php echo $_l->tpl('sidebar.files'); ?></a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.server_sett'); ?></strong></a>
					
					<a href="../settings.php" class="list-group-item"><?php echo $_l->tpl('sidebar.manage'); ?></a>
					<a href="index.php" class="list-group-item active"><?php echo $_l->tpl('sidebar.plugins'); ?></a>
				</div>
			</div>
			<div class="col-9">
				<?php echo $errorMessage; ?>
				<div id="p_install_one" class="alert alert-warning" style="display:none;">
					<i class="fa fa-spinner fa fa-spin"></i> <?php echo $_l->tpl('node.plugins.installing'); ?> <strong><?php echo $_l->tpl('node.plugins.installing_warning'); ?></strong>
				</div>
				<div id="p_install_two" class="alert alert-success" style="display:none;">
					<?php echo $_l->tpl('node.plugins.installed'); ?>
				</div>
				<table class="table table-striped table-bordered table-hover">
					<thead>
						<tr>
							<th><?php echo $_l->tpl('string.name'); ?></th>
							<th><?php echo $_l->tpl('string.description'); ?></th>
							<th></th>
						</tr>
					</thead>
					<tbody>
						<?php echo $searchResults; ?>
					</tbody>
				</table>
				<ul class="pager">
					<?php echo $pageOptions; ?>
				</ul>
			</div>
		</div>
		<div class="footer">
			<?php include('../../assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$(document).ready(function(){
			$(".install").click(function(){
				var sendme = $(this).attr("id");
				$("#p_install_one").slideDown();
				$.ajax({
					type: "POST",
					url: "../ajax/plugins/install.php",
					data: { plugin: sendme },
			  		success: function(data) {
			  			$("#p_install_one").slideUp("fast", function(){$("#p_install_two").slideDown();});
			 		}
				});
			});
		});
	</script>
</body>
</html>