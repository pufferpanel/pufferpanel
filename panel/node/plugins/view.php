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

$filesIncluded = true;

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false){

	Page\components::redirect($core->settings->get('master_url').'index.php');
	exit();
}

if(isset($_GET['slug']) && !empty($_GET['slug'])){

	/*
	 * Viewing Plugin
	 */
	$_GET['slug'] = str_replace(array(' ', '+', '%20'), '', $_GET['slug']);
	$context = stream_context_create(array(
		"http" => array(
			"method" => "GET",
			"header" => 'User-Agent: PufferPanel',
			"timeout" => 5
		)
	));
	$data = file_get_contents('http://api.bukget.org/3/plugins/bukkit/'.$_GET['slug'], false, $context);
	$data = json_decode($data, true);
	
	$data['description'] = (strlen($data['description']) == 0) ? 'No description is avaliable for this plugin.' : $data['description'];
	
	if(empty($data['authors'])){
		$pluginAuthors = 'none specified';
	}else{
		$pluginAuthors = '';
		foreach($data['authors'] as $id => $name){ $pluginAuthors .= '<a href="http://dev.bukkit.org/profiles/'.$name.'/" target="_blank">'.$name.'</a>, '; }
		$pluginAuthors = rtrim($pluginAuthors, ', ');
	}
	
	$downloads = '';
	$i = 0;
	
	$data['versions'] = (is_array($data['versions'])) ? $data['versions'] : array($data['versions']);
	foreach($data['versions'] as $id => $value){
			
		$gameVersions = '';
		foreach($value['game_versions'] as $gid => $ver){
			$gameVersions .= $ver.'<br />';
		}
		
		$gameVersions = substr($gameVersions, 0, -6);
		$gameVersions = str_replace("CB", "", $gameVersions);
		
		$downloads .= '
						<tr>
							<td><a href="#install" id="'.$_GET['slug'].'|'.$i.'" class="install"><i class="fa fa-download"></i></a></td>
							<td>'.$value['filename'].'</td>
							<td>'.$value['version'].'</td>
							<td>'.date('M n, Y', $value['date']).'</td>
							<td>'.$gameVersions.'</td>
							<td><a href="#" data-toggle="popover" data-content="'.$value['md5'].'" data-original-title="MD5 Checksum">...</a></td>
						</tr>
						';
		
		$i++;
	
	}

}else{

	Page\components::redirect('search.php');
	exit();
	
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
					<a href="#" class="list-group-item list-group-item-heading"><strong>Account Actions</strong></a>
					<a href="../../account.php" class="list-group-item">Settings</a>
					<a href="../../servers.php" class="list-group-item">My Servers</a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Server Actions</strong></a>
					<a href="../index.php" class="list-group-item">Overview</a>
					<a href="../console.php" class="list-group-item">Live Console</a>
					<a href="../files/index.php" class="list-group-item">File Manager</a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Server Settings</strong></a>
					
					<a href="../settings.php" class="list-group-item">Server Management</a>
					<a href="index.php" class="list-group-item active">Server Plugins</a>
				</div>
			</div>
			<div class="col-9">
				<div class="panel panel-primary">
					<div class="panel-heading">
						<h3 class="panel-title">About <?php echo $data['plugin_name']; ?></h3>
					</div>
					<div class="panel-body">
						<p class="text-muted"><small>This plugin is created and maintained by <?php echo $pluginAuthors; ?>. View this plugin <a href="http://dev.bukkit.org/bukkit-plugins/<?php echo $_GET['slug']; ?>/" target="_blank">on Bukkit</a>.</small></p>
						<p><?php echo $data['description']; ?></p>
					</div>
				</div>
				<div id="p_install_one" class="alert alert-warning" style="display:none;">
					<i class="fa fa-spinner fa fa-spin"></i> Please wait while your plugin is installing. This process could take about a minute to complete. <strong>Do not navigate away from this page!</strong>
				</div>
				<div id="p_install_two" class="alert alert-success" style="display:none;">
					Your plugin has been installed successfully.
				</div>
				<table class="table table-striped table-bordered table-hover">
					<thead>
						<tr>
							<th></th>
							<th>Name</th>
							<th>Version</th>
							<th>Published</th>
							<th>Versions</th>
							<th>MD5</th>
						</tr>
					</thead>
					<tbody>
						<?php echo $downloads; ?>
					</tbody>
				</table>
			</div>
		</div>
		<div class="footer">
			<?php include('../../assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$(document).ready(function(){
			$('a[data-toggle="popover"]').popover({'placement' : 'left', 'trigger' : 'hover'});
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