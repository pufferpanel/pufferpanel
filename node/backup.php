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
$backupStatus = '';
if(isset($_GET['do']) && $_GET['do'] == 'create'){

	if(isset($_POST['dbu'])){
	
		$nodeSQLConnect = $mysql->prepare("SELECT * FROM `nodes` WHERE `node` = ? LIMIT 1");
		$nodeSQLConnect->execute(array($core->framework->server->getData('node')));
		
		$node = $nodeSQLConnect->fetch();
	
		/*
		 * Does User have Space!
		 */
		$files = glob($node['backup_dir'].$core->framework->server->getData('name').'/*', GLOB_MARK);
		
			$totalSpace = 0;
			$totalFiles = 0;
			foreach($files as $file){
			
				$stat = stat($file);
				$bytes = ($stat['size'] / 1024) / 1024;
				
				$totalFiles = $totalFiles + 1;
				$totalSpace = floor($totalSpace + $bytes);
				
			}
			
				/*
				 * Show Warnings
				 */
				$continue = true;
				if($totalSpace >= ($core->framework->server->getData('backup_disk_limit') - 500)) {
				
					$continue = false;
					$backupStatus .= '<div class="error-box round">You have reached your allotment of backup space. Please delete older backups before continuing.</div>';
					
				}
				
				if($totalFiles >= $core->framework->server->getData('backup_file_limit')){
				
					$continue = false;
					$backupStatus .= '<div class="error-box round">You have reached your limit for stored backups. Please delete older ones before continuing.</div>';
				
				} 
				
		
		/*
		 * Do Backup
		 */
		
		$user = $core->framework->server->getData('name');
		$filename = str_replace(" ", "", stripslashes(trim($_POST['backup_name'])));
		$filename = ltrim($filename, '/');
		$filename = rtrim($filename, '/');
		$filename = rtrim($filename, '.');
		$filename = preg_replace('/[^A-Za-z0-9_-]/', '', $filename);
		$filename = preg_replace('/(\.){2,}/', '.', $filename);
		$filename = $filename.'_'.$core->framework->auth->keygen('10');
		if(isset($_POST['save_options'])){ $saveOptions = true; }else{ $saveOptions = false; }
		$backupToken = 'TOKEN-'.$core->framework->auth->keygen('5').'-'.$core->framework->auth->keygen('10').'-'.$core->framework->auth->keygen('5');
		
			/*
			 * Save Options
			 */
			if($saveOptions === true){
			
				$checkSave = $mysql->prepare("SELECT * FROM `backup_datastore` WHERE `server` = ?");
				$checkSave->execute(array($core->framework->server->getData('hash')));
				
					if($checkSave->rowCount() > 0){
					
						/*
						 * Update
						 */
						$updateDatastore = $mysql->prepare("UPDATE `backup_datastore` SET `backup_pattern` = :params WHERE `server` = :server");
						$updateDatastore->execute(array(
						 	':params' => $_POST['backup_params'],
						 	':server' => $core->framework->server->getData('hash')
						));
						 
					
					}else{
					
						/*
						 * Insert
						 */
						$insertDatastore = $mysql->prepare("INSERT INTO `backup_datastore` VALUES(NULL, :server, :params)");
						$insertDatastore->execute(array(
							':params' => $_POST['backup_params'],
							':server' => $core->framework->server->getData('hash')
						));
						
					}
			
			}
		
			/*
			 * Clean up Backup Data
			 */
			$data = $_POST['backup_params'];
			
			if($data == ''){ $data = '*'; }
			
			if($_POST['backup_params'] != '' && $_POST['backup_params'] != '*') {
			
				$lines = explode("\n", $data);
				
					$backup = '';
					$skip = '';
					foreach($lines as $line){
					
						/*
						 * Remove Whitespace
						 */
						$line = str_replace(array(" ", "\\"), "", trim($line));
						$line = preg_replace("#/+#", "/", $line);
						$firstChara = substr($line, 0, 1);
							/* 
							 * Check if line is commented out & that it is also not blank
							 */
							if($firstChara != "#" && preg_match('~[0-9a-z]~i', $line) > 0){
							
								if($firstChara == '!'){
								
									/*
									 * Skip Folder or File
									 */
									if($skip == ''){
										$skip .= ltrim(rtrim($line, "/"), "!");
									}else{
										$skip .= ','.ltrim(rtrim($line, "/"), "!");
									}
								
								}else{
								
									/*
									 * Backup Folder
									 */
									if($backup == ''){
										$backup .= rtrim($line, "/");
									}else{
										$backup .= ','.rtrim($line, "/");
									}
								
								}
							
							}
					
					}
					
				}else{
				
					$backup = 'all';
				
				}
	
					/*
					 * Connections
					 */
					$con = ssh2_connect($node['node_ip'], 22);
					ssh2_auth_password($con, $node['username'], openssl_decrypt($node['password'], 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($node['encryption_iv'])));
					
					/*
					 * If server is running we need to do this differently.
					 */
					if($core->framework->rcon->online($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port')) === true){
					
						ssh2_exec($con, 'cd /srv/scripts; ./send_command.sh '.$core->framework->server->getData('name').' "save-all"');
						sleep(2);
						ssh2_exec($con, 'cd /srv/scripts; ./send_command.sh '.$core->framework->server->getData('name').' "save-off"');
					
					}
					
					
						/*
						 * Proceed with Backup
						 */
						$newBackup = $mysql->prepare("INSERT INTO `backups` VALUES(NULL, :hash, :token, :filename, 0, NULL, :time, NULL, 0, NULL, NULL)");
						
						if($newBackup->execute(array(
							':hash' => $core->framework->server->getData('hash'),
							':token' => $backupToken,
							':filename' => $filename,
							':time' => time()
						))){
						
								if($backup == 'all'){
								                                    
									$s = ssh2_exec($con, 'cd /srv/scripts; ./backup_server.sh '.$core->framework->server->getData('name').' "'.$filename.' '.$core->framework->server->getData('node').' '.$core->framework->server->getData('hash').' '.$backupToken.'" "*" "'.escapeshellcmd(str_replace(",", " ", $skip)).'"');

								}else{
													                                    
									$s = ssh2_exec($con, 'cd /srv/scripts; ./backup_server.sh '.$core->framework->server->getData('name').' "'.$filename.' '.$core->framework->server->getData('node').' '.$core->framework->server->getData('hash').' '.$backupToken.'" "'.escapeshellcmd(str_replace(",", " ", $backup)).'" "'.escapeshellcmd(str_replace(",", " ", $skip)).'"');
									
								}
								
								stream_set_blocking($s, true);
								
								$backupStatus .= '<div class="confirmation-box round">Your backup has been started. Backups will appear at the bottom of this page when they are finished.</div>';
						
						}else{
						
							ssh2_exec($con, 'cd /srv/scripts; ./send_command.sh '.$core->framework->server->getData('name').' "save-on"');
							
							$backupStatus .= '<div class="error-box round">A MySQL error was encountered and to ensure the integrity of your server we have stopped this backup process. Please try again or contact support if this error continues to occur.</div>';
						
						}
	
	}else{
	
		$backupStatus = "<div class='error-box round'>We were unable to preform the requested operation. This was due to the lack of any data sent to the function. Please retry your request.</div>";
	
	}

}

$query = $mysql->prepare("SELECT * FROM `backups` WHERE `server` = ? AND `timeend` IS NOT NULL ORDER BY `id` DESC");
$query->execute(array($core->framework->server->getData('hash')));

	if($query->rowCount() > 0){
			
			$returnBackups = '';
			while($row = $query->fetch()){
						
				$nodeSQLConnect = $mysql->prepare("SELECT * FROM `nodes` WHERE `node` = ? LIMIT 1");
				$nodeSQLConnect->execute(array($core->framework->server->getData('node')));
				
					$node = $nodeSQLConnect->fetch();
					$stat = stat($node['backup_dir'].$core->framework->server->getData('name').'/'.$row['file_name'].'.tar.gz');
				
				$md5 = (strlen($row['md5']) > 25) ? substr($row['md5'], 0, 22).'...' : $row['md5'];
				$sha1 = (strlen($row['sha1']) > 25) ? substr($row['sha1'], 0, 22).'...' : $row['sha1'];
				
				$returnBackups .= '<tr>
										<td style="text-align:center;"><a href="core/ajax/backup/download.php?id='.$row['id'].'"><i class="fa fa-download"></i></a>&nbsp;&nbsp;&nbsp;<i class="fa fa-times-circle-o"></i></td>
										<td>'.$row['file_name'].'</td>
										<td><abbr title="'.$row['md5'].'">'.$md5.'</abbr></td>
										<td><abbr title="'.$row['sha1'].'">'.$sha1.'</abbr></td>
										<td>'.date('M n, Y \a\t g:ia', $row['timeend']).'</td>
										<td>'.$core->framework->files->formatSize($stat['size']).'</td>
									</tr>';
			
			}
			
		$backupList = '<table>
							<thead>
								<tr>
									<th style="width:5%"></th>
									<th style="width:25%">Name</th>
									<th style="width:20%">MD5 Hash</th>
									<th style="width:20%">SHA1 Hash</th>
									<th style="width:20%">Date</th>
									<th style="width:5%">Size</th>
								</tr>
							</thead>
							<tbody>
								'.$returnBackups.'
							</tbody>
						</table>';
	
	}else{
	
		$backupList = '<div class="warning-box round">You don\'t have any server backups avaliable. You should create one right now in order to protect and ensure the integrity of your server.</div>';
	
	}
	
$selectParams = $mysql->prepare("SELECT `backup_pattern` FROM `backup_datastore` WHERE `server` = ?");
$selectParams->execute(array($core->framework->server->getData('hash')));

	if($selectParams->rowCount() == 1){
		$sRow = $selectParams->fetch();
		$backupParams = $sRow['backup_pattern'];
	}else{
		$backupParams = '';
	}

?>
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>PufferPanel - Server Backup Manager</title>
	
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
				<li><a href="#" class="round button dark"><i class="fa fa-user"></i>&nbsp;&nbsp; <strong><?php echo $core->framework->user->getData('username'); ?></strong></a></li>
				<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>servers.php" class="round button dark"><i class="fa fa-home"></i></a></li>
				<li><a class="round button dark"><i class="fa fa-hdd"></i>&nbsp;&nbsp; <?php echo $core->framework->server->getData('name'); ?></a></li>
			</ul>
			<ul id="nav" class="fr">
				<?php if($core->framework->user->getData('root_admin') == 1){ echo '<li><a href="'.$core->framework->settings->get('master_url').'admin/index.php" class="round button dark"><i class="fa fa-bar-chart-o"></i>&nbsp;&nbsp; Admin CP</a></li>'; } ?>
				<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>logout.php" class="round button dark"><i class="fa fa-power-off"></i></a></li>
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
					<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>account.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Edit Settings</a></li>
					<li><a href="<?php echo $core->framework->settings->get('master_url'); ?>servers.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> My Servers</a></li>
				</ul>
				<h3>Server Actions</h3>
				<ul>
					<li><a href="index.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Overview</a></li>
					<li><a href="console.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Live Console</a></li>
					<li><a href="settings.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Server Settings</a></li>
					<li><a href="plugins.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Server Plugins</a></li>
					<li><a href="files.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> File Manager</a></li>
					<li><a href="backup.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Backup Manager</a></li>
				</ul>
			</div>
			<div class="side-content fr">
				<div class="content-module">
					<div class="content-module-heading cf">
						<h3 class="fl">Create a Backup</h3>
					</div> <!-- end content-module-heading -->
					<div class="content-module-main" id="server_info">
						<?php echo $backupStatus; ?>
						<form action="backup.php?do=create" method="post">
							<fieldset>
								<p>
									<label for="backup_name">Backup Name</label>
									<input type="text" name="backup_name" class="round" style="width:425px;" value="backup-<?php echo date('d-n-Y_\a\t_H-i-s', time()); ?>"/>
								</p>
								<p>
									<label for="backup_params">Backup Parameters</label>
									<textarea name="backup_params" class="round" style="height:125px;width:425px;font-size: 1em;"><?php echo $backupParams; ?></textarea>
									<em>This is an advanced function. Please <a href="#example" id="v_ex">view the example</a> before making a backup.</em>
									<div id="v_ex_showcontent" class="information-box round" style="display:none;background: #e5f5f9;padding: 0.833em;">
										#Include these folders:<br />plugins<br />world<br />world_nether<br />world_the_end<br /><br />#Also include specific files not in those directories:<br />server.jar<br />server.properties<br />somefolder/t/test.txt<br /><br />#Exclude DynMap and Coreprotect:<br />!plugins/dynmap<br />!plugins/coreprotect<br /><br />#Exclue all plugin jars.<br />!plugins/*.jar
									</div>
								</p>
								<p>
									<label for="save_options">Save Backup Parameters</label>
									<input type="checkbox" name="save_options" value="1"/> Yes, please save the parameters I set above.
								</p>
								<div class="stripe-separator"><!--  --></div>
								<input type="submit" value="Generate Backup" name="dbu" class="round blue ic-right-arrow" />
							</fieldset>
						</form>
					</div> <!-- end content-module-main -->
				</div>
			</div>
			<div class="side-content fr">
				<div class="content-module">
					<div class="content-module-heading cf">
						<h3 class="fl">Current Backed Up Files</h3>
					</div> <!-- end content-module-heading -->
					<div class="content-module-main" id="server_info">
						<?php echo $backupList; ?>
					</div> <!-- end content-module-main -->
				</div>
			</div>
		</div>
	</div>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.3 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
	<script type="text/javascript">
		$("#v_ex").click(function(){$("#v_ex_showcontent").slideToggle();});
	</script>
</body>
</html>