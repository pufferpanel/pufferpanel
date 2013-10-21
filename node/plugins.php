<?php
session_start();
require_once('core/framework/framework.core.php');

$filesIncluded = true;

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), $core->framework->auth->getCookie('pp_server_hash')) === false){

	$core->framework->page->redirect($core->framework->settings->get('master_url').'index.php');
	exit();
}

/*
 * Are we on the correct node?
 */
//$url = parse_url($_SERVER["SERVER_NAME"], PHP_URL_PATH);
//$parts = explode('.', $url);
//
//	if($parts[0] != $core->framework->server->getData('node')){
//		$core->framework->page->redirect($core->framework->settings->get('master_url').'index.php');
//	}

if(isset($_GET['do']) && $_GET['do'] == 'view' && isset($_GET['slug']) && !empty($_GET['slug'])){

	/*
	 * Viewing Plugin
	 */
	$_GET['slug'] = str_replace(array(' ', '+', '%20'), '', $_GET['slug']);
	$data = file_get_contents('http://api.bukget.org/3/plugins/bukkit/'.$_GET['slug']);
	$data = json_decode($data, true);
	
	$data['description'] = (strlen($data['description']) == 0) ? 'No description is avaliable for this plugin.' : $data['description'];
	
	if(empty($data['authors'])){
		$pluginAuthors = 'none specified';
	}else{
		$pluginAuthors = '';
		foreach($data['authors'] as $id => $name){ $pluginAuthors .= $name.', '; }
		$pluginAuthors = rtrim($pluginAuthors, ', ');
	}
	
	$downloads = '';
	$i = 0;
	
	$data['versions'] = (is_array($data['versions'])) ? $data['versions'] : array($data['versions']);
	foreach($data['versions'] as $id => $value){
			
		$gameVersions = '';
		foreach($value['game_versions'] as $gid => $ver){
			$gameVersions .= 'CraftBukkit '.$ver.'<br />';
		}
		
		$gameVersions = substr($gameVersions, 0, -6);
		$gameVersions = str_replace("CB", "", $gameVersions);
		
		$downloads .= '
						<tr>
							<td><a href="#install" id="'.$_GET['slug'].'|'.$i.'" class="install"><i class="icon-download"></i></a></td>
							<td>'.$value['filename'].'</td>
							<td>'.$value['version'].'</td>
							<td>'.date('M n, Y \a\t g:ia', $value['date']).'</td>
							<td>'.$gameVersions.'</td>
							<td>'.$value['md5'].'</td>
						</tr>
						';
		
		$i++;
	
	}

	$pluginName = $data['plugin_name'];
	$pluginInformation = '
							<p><em>This plugin was created and is maintained by: '.$pluginAuthors.'</em></p>
							<table>
								<thead>
									<th style="width:5%;">DL</th>
									<th style="width:20%;">File Name</th>
									<th style="width:5%;">Version</th>
									<th style="width:20%;">Publish Date</th>
									<th style="width:20%;">Bukkit Versions</th>
									<th style="width:30%;">MD5</th>
								</thead>
								<tbody>
									'.$downloads.'
								</tbody>
							</table>
						';
						
		$divContent = '<p>'.$data['description'].'</p>';			

}else if(isset($_GET['do']) && $_GET['do'] == 'search' && isset($_GET['term']) && !empty($_GET['term'])){

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
			
			$searchResults .= '	<tr class="hover_fade_row">
									<td><a href="#install" id="'.$value['slug'].'|0" class="install"><i class="icon-download"></i></a></td>
									<td><a href="plugins.php?do=view&slug='.$value['slug'].'">'.$value['plugin_name'].'</a></td>
									<td>'.$value['description'].'</td>
								</tr>';
		
		}
		
	if(count($data) == 25){
		
		if(isset($_GET['start']) && $_GET['start'] > 24){
		
			$pageOptions = '<a href="plugins.php?do=search&term='.$_GET['term'].'&start='.($_GET['start'] - 25).'" class="round button blue text-upper small-button">Previous Page</a>&nbsp;&nbsp;<a href="plugins.php?do=search&term='.$_GET['term'].'&start='.($_GET['start'] + 25).'" class="round button blue text-upper small-button">Next Page</a>';
		
		}else{
		
			$pageOptions = '<a href="plugins.php?do=search&term='.$_GET['term'].'&start='.($_GET['start'] + 25).'" class="round button blue text-upper small-button">Next Page</a>';
		
		}
		
	}else{
		
		if(isset($_GET['start']) && $_GET['start'] != 0){
		
			$pageOptions = '<a href="plugins.php?do=search&term='.$_GET['term'].'&start='.($_GET['start'] - 25).'" class="round button blue text-upper small-button">Previous Page</a>';
		
		}else{
		
			$pageOptions = '';
		
		}
		
	}

}

?>
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title><?php echo $core->framework->settings->get('company_name'); ?> - Server Dashboard</title>
	
	<!-- Stylesheets -->
	<link href='http://fonts.googleapis.com/css?family=Droid+Sans:400,700' rel='stylesheet'>
	<link rel="stylesheet" href="<?php echo $core->framework->settings->get('master_url'); ?>assets/css/style.css">
	
	<!-- Optimize for mobile devices -->
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
	
	<!-- jQuery & JS files -->
	<script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
	
	<script type="text/javascript" src="<?php echo $core->framework->settings->get('master_url'); ?>assets/javascript/jquery.ba-throttle-debounce.min.js"></script>
</head>
<body>
	<div id="top-bar">
		<div class="page-full-width cf">
			<ul id="nav" class="fl">
				<li><a href="#" class="round button dark"><i class="icon-user"></i>&nbsp;&nbsp; <strong><?php echo $core->framework->user->getData('username'); ?></strong></a></li>
				<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>servers.php" class="round button dark"><i class="icon-hdd"></i></a></li>
			</ul>
			<ul id="nav" class="fr">
				<?php if($core->framework->user->getData('root_admin') == 1){ echo '<li><a href="'.$core->framework->settings->get('master_url').'admin/index.php" class="round button dark"><i class="icon-bar-chart"></i>&nbsp;&nbsp; Admin CP</a></li>'; } ?>
				<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>logout.php" class="round button dark"><i class="icon-off"></i></a></li>
			</ul>
		</div>	
	</div>
	<div id="header-with-tabs">
		<div class="page-full-width cf">
		</div>
	</div>
	<div id="content">
		<div class="page-full-width cf">
			<div class="side-menu fl">
				<h3>Account Actions</h3>
				<ul>
					<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>account.php"><i class="icon-double-angle-right pull-right menu-arrows"></i> Edit Settings</a></li>
					<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>servers.php"><i class="icon-double-angle-right pull-right menu-arrows"></i> My Servers</a></li>
				</ul>
				<h3>Server Actions</h3>
				<ul>
					<li><a href="index.php"><i class="icon-double-angle-right pull-right menu-arrows"></i> Overview</a></li>
					<li><a href="console.php"><i class="icon-double-angle-right pull-right menu-arrows"></i> Live Console</a></li>
					<li><a href="settings.php"><i class="icon-double-angle-right pull-right menu-arrows"></i> Server Settings</a></li>
					<li><a href="plugins.php"><i class="icon-double-angle-right pull-right menu-arrows"></i> Server Plugins</a></li>
					<li><a href="files.php"><i class="icon-double-angle-right pull-right menu-arrows"></i> File Manager</a></li>
					<li><a href="backup.php"><i class="icon-double-angle-right pull-right menu-arrows"></i> Backup Manager</a></li>
				</ul>
			</div>
			<?php if(isset($_GET['do']) && $_GET['do'] == 'view' && isset($_GET['slug']) && !empty($_GET['slug'])){ ?>
			<div class="side-content fr">
				<div class="content-module">
					<div class="content-module-heading cf">
						<h3 class="fl">Plugin: <?php echo $pluginName; ?></h3>
					</div>
					<div id="p_install_one" class="content-module-main cf" style="display:none;">
						<div class="warning-box round"><i class="icon-spinner icon-spin"></i> Please wait while your plugin is installing. This process could take about a minute to complete. <strong>Do not navigate away from this page!</strong></div>
					</div>
					<div id="p_install_two" class="content-module-main cf" style="display:none;">
						<div class="confirmation-box round">Your plugin has been installed.</div>
					</div>
					<div class="content-module-main cf">
						<?php echo $pluginInformation; ?>
					</div>
				</div>
			</div>
			<div class="side-content fr">
				<div class="content-module">
					<div class="content-module-heading cf">
						<h3 class="fl">Plugin Description</h3>
					</div>
					<div class="content-module-main cf" id="plugin_desc">
						<?php print $divContent; ?>
					</div>
				</div>
			</div>
			<?php }else if(isset($_GET['do']) && $_GET['do'] == 'search' && isset($_GET['term']) && !empty($_GET['term'])){ ?>
			<div class="side-content fr">
				<div class="content-module">
					<div class="content-module-heading cf">
						<h3 class="fl">Search results for: "<?php echo $_GET['term']; ?>"</h3>
					</div>
					<div id="p_install_one" class="content-module-main cf" style="display:none;">
						<div class="warning-box round"><i class="icon-spinner icon-spin"></i> Please wait while your plugin is installing. This process could take about a minute to complete. <strong>Do not navigate away from this page!</strong></div>
					</div>
					<div id="p_install_two" class="content-module-main cf" style="display:none;">
						<div class="confirmation-box round">Your plugin has been installed.</div>
					</div>
					<div class="content-module-main cf">
						<table>
							<thead>
								<th style="width:5%;">DL</th>
								<th style="width:25%;">Name</th>
								<th style="width:70%;">Description</th>
							</thead>
							<tbody>
								<?php echo $searchResults; ?>
							</tbody>
						</table>
						<div class="center"><?php echo $pageOptions; ?></div>
					</div>
				</div>
			</div>
			<?php }else{ ?>
			<div class="side-content fr">
				<div class="half-size-column fl">
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Plugin Search</h3>
						</div>
						<div class="content-module-main">
							<form action="plugins.php" id="plugin_search_form" method="get">
								<fieldset>
									<input type="hidden" name="do" value="search"/>
									<p>
										<label for="term">Plugin Name</label>
										<input type="text" name="term" class="round full-width-input"/>
									</p>
									<div class="stripe-separator"><!--  --></div>
									<input type="submit" value="Search Plugin Repository" class="round blue ic-right-arrow" />
								</fieldset>
							</form>	
						</div>
					</div>
				</div>
				<div class="half-size-column fr">
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Plugin Results</h3>
						</div>
						<div class="content-module-main cf" id="plugin_results">
							<p id="hide_plugin_results" class="nomargin">No search term was entered.</p>
						</div>
					</div>
				</div>
			</div>
			<?php } ?>
		</div>
	</div>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.3 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
	<script type="text/javascript">
		$(document).ready(function(){
			$(".install").click(function(){
				var sendme = $(this).attr("id");
				$("#p_install_one").slideDown();
				$.ajax({
					type: "POST",
					url: "core/ajax/plugins/install.php",
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