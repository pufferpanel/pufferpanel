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

$canEdit = array('txt', 'yml', 'log', 'conf', 'html', 'json', 'properties', 'props', 'cfg', 'lang');
$parName = '';

/*
 * Security Patch
 */
if(isset($_GET['file'])){

	$_GET['file'] = str_replace('..', '', urldecode($_GET['file']));

}

if(isset($_GET['dir'])){

	$_GET['dir'] = str_replace('..', '', urldecode($_GET['dir']));

}

if(isset($_POST['file'])){

	$_POST['file'] = str_replace('..', '', urldecode($_POST['file']));

}

/*
 * End Security Patch
 */

if(!isset($_GET['do'])){	
	
	/*
	 * Display File Manager Overview Page
	 */
	if(isset($_GET['dir']) && !empty($_GET['dir'])){
	
		/*
		 * Check First Character
		 */
		if(substr($_GET['dir'], 0, 1) == '/'){ $core->framework->page->redirect('files.php?dir='.substr($_GET['dir'], 1)); }
	
		/*
		 * Validate Directory
		 */
		if(!is_dir($core->framework->server->getData('path').$_GET['dir'])){ $core->framework->page->redirect('files.php'); }
		$parName = '(Viewing: /'.$_GET['dir'].')';
		
		/*
		 * Inside a Directory
		 */
		$goBack = explode('/', $_GET['dir']);
		$gbTotal = count($goBack);
		
			if(count($goBack) <= 2){
			
				$showUpOne = '';
			
			}else{
			
				unset($goBack[$gbTotal - 1]);
				unset($goBack[$gbTotal - 2]);
				$previousDir = implode('/', $goBack).'/';
				
				$showUpOne = '	<tr>
									<td><img src="'.$core->framework->settings->get('master_url').'assets/images/node/file_manager/directory.png" alt="DIR" title="Directory"/></td>
									<td><a href="files.php?dir='.$previousDir.'">&larr; '.$previousDir.'</a></a></td>
									<td></td>
									<td></td>
									<td class="center"></td>
								</tr>';
				
			}
		
		$displayFolders = 	'<tr>
								<td><img src="'.$core->framework->settings->get('master_url').'assets/images/node/file_manager/directory.png" alt="DIR" title="Directory"/></td>
								<td><a href="files.php">&larr;</a></a></td>
								<td></td>
								<td></td>
								<td class="center"></td>
							</tr>'.$showUpOne;
							
		if($_GET['dir'] == '/'){ $displayFolders = ''; }
		$displayFiles = '';
		$files = glob($core->framework->server->getData('path').$_GET['dir']."*", GLOB_MARK);
		
		/*
		 * Iterate through Files & Directories
		 */	
		foreach($files as $file){
			
			/*
			 * Get Stats on File
			 */
			$stat = stat($file);
			#$filesize = $core->framework->files->formatSize($stat['size']);
			
			/*
			 * Iterate into HTML Variable
			 */
			if(is_dir($file)){
			
				$displayFolders .= 	'<tr>
										<td><i class="icon-folder-open">&nbsp</i></td>
										<td><a href="files.php?dir='.urlencode(str_replace($core->framework->server->getData('path'), '', $file)).'">'.str_replace($core->framework->server->getData('path').$_GET['dir'], '', $file).'</a></td>
										<td></td>
										<td>'.date('m/d/y H:i:s', $stat['mtime']).'</td>
										<td class="center"><!--<img src="'.$core->framework->settings->get('master_url').'assets/images/node/file_manager/transparent.png"/>&nbsp;&nbsp;&nbsp;<a href="files.php?do=delete&file='.urlencode(str_replace($core->framework->server->getData('path'), '', $file)).'"><i class="icon-remove-circle">&nbsp</i></a>--></td>
									</tr>';
			
			}else{
			
				if(in_array(pathinfo($file, PATHINFO_EXTENSION), $canEdit)){
					$url = '<a href="files.php?do=edit&file='.urlencode(str_replace($core->framework->server->getData('path'), '', $file)).'">'.str_replace($core->framework->server->getData('path').$_GET['dir'], '', $file).'</a>';
				}else{
					$url = str_replace($core->framework->server->getData('path').$_GET['dir'], '', $file);
				}
				$displayFiles .= 	'<tr>
										<td><i class="icon-file-text"></i></td>
										<td>'.$url.'</td>
										<td>'.$core->framework->files->formatSize($stat['size']).'</td>
										<td>'.date('m/d/y H:i:s', $stat['mtime']).'</td>
										<td class="center"><a href="files.php?do=download&file='.urlencode(str_replace($core->framework->server->getData('path'), '', $file)).'"><i class="icon-download-alt">&nbsp</i></a>&nbsp;&nbsp;&nbsp;<!--<a href="files.php?do=delete&file='.urlencode(str_replace($core->framework->server->getData('path'), '', $file)).'"><i class="icon-remove-circle">&nbsp</i></a>--></td>
									</tr>';
			
			}
		
		}
		
	}else{
	
		/*
		 * Not in a Directory
		 */
		$displayFolders = '';
		$displayFiles = '';
		$files = glob($core->framework->server->getData('path')."*", GLOB_MARK);
		
		/*
		 * Iterate through Files & Directories
		 */	
		foreach($files as $file){
			
			/*
			 * Get Stats on File
			 */
			$stat = stat($file);
			#$filesize = $core->framework->files->formatSize($stat['size']);
			
			/*
			 * Iterate into HTML Variable
			 */
			if(is_dir($file)){
			
				$displayFolders .= 	'<tr>
										<td><i class="icon-folder-open">&nbsp</i></td>
										<td><a href="files.php?dir='.urlencode(str_replace($core->framework->server->getData('path'), '', $file)).'">'.str_replace($core->framework->server->getData('path'), '', $file).'</a></td>
										<td></td>
										<td>'.date('m/d/y H:i:s', $stat['mtime']).'</td>
										<td class="center"><!--<img src="'.$core->framework->settings->get('master_url').'assets/images/node/file_manager/transparent.png"/>&nbsp;&nbsp;&nbsp;<a href="files.php?do=delete&file='.urlencode(str_replace($core->framework->server->getData('path'), '', $file)).'"><i class="icon-remove-circle">&nbsp</i></a>--></td>
									</tr>';
			
			}else{
			
				if(in_array(pathinfo($file, PATHINFO_EXTENSION), $canEdit)){
					$url = '<a href="files.php?do=edit&file='.urlencode(str_replace($core->framework->server->getData('path'), '', $file)).'">'.str_replace($core->framework->server->getData('path'), '', $file).'</a>';
				}else{
					$url = str_replace($core->framework->server->getData('path'), '', $file);
				}
					
				$displayFiles .= 	'<tr>
										<td><i class="icon-file-text"></i></td>
										<td>'.$url.'</td>
										<td>'.$core->framework->files->formatSize($stat['size']).'</td>
										<td>'.date('m/d/y H:i:s', $stat['mtime']).'</td>
										<td class="center"><a href="files.php?do=download&file='.urlencode(str_replace($core->framework->server->getData('path'), '', $file)).'"><i class="icon-download-alt">&nbsp</i></a>&nbsp;&nbsp;&nbsp;<!--<a href="files.php?do=delete&file='.urlencode(str_replace($core->framework->server->getData('path'), '', $file)).'"><i class="icon-remove-circle">&nbsp</i></a>--></td>
									</tr>';
			
			}
		
		}
	
	}
	
	/*
	 * Setup Basic Display
	 */
	$HTML = '<table>
				<thead>
					<tr>
						<th style="width:5%"></th>
						<th style="width:45%">File Name</th>
						<th style="width:20%">File Size</th>
						<th style="width:20%">Last Modified</th>
						<th style="width:10%" class="center">Options</th>
					</tr>
				</thead>
				<tbody>
					'.$displayFolders.$displayFiles.'
				</tbody>
			</table>';

}else{

	if($_GET['do'] == 'edit'){
	
		if(!isset($_GET['action'])){
		
			/*
			 * Editing File
			 */
			if(isset($_GET['file']) && !is_dir($_GET['file']) && file_exists($core->framework->server->getData('path').$_GET['file'])){
			$_GET['file'] = $_GET['file'];
			
				if(in_array(pathinfo($core->framework->server->getData('path').$_GET['file'], PATHINFO_EXTENSION), $canEdit)){
				
					/*
					 * Name Error Messages
					 */
					(isset($_GET['error'])) ? $errormsg = '<div class="error-box round">'.base64_decode($_GET['error']).' You can return to the file manager by <a href="files.php">clicking here</a>.</div>' : $errormsg = '';
						
					/*
					 * Begin Advanced Saving
					 */
					$saveDir = '/tmp/'.$core->framework->server->getData('hash').'/';
					
						/*
						 * Check that Secure User Directory Exists
						 */
						if(!is_dir($saveDir)){
						
							/*
							 * Make Directory
							 */
							mkdir($saveDir);
						
						}
						
					/*
					 * SFTP Connect to Show File
					 */
					$file = pathinfo($core->framework->server->getData('path').$_GET['file'], PATHINFO_BASENAME);
					$directory = dirname($_GET['file']).'/';
					
					/*
					 * Directory Cleaning
					 */
					if($directory == './' || $directory == '.'){ $directory = ''; }
					if(substr($directory, 0, 1) == '/'){ $directory = substr($directory, 1); }
												
						/*
						 * Download Via SFTP
						 */
						$SFTPConnection = ssh2_connect($core->framework->server->getData('ftp_host'), 22);
						ssh2_auth_password($SFTPConnection, $core->framework->server->getData('ftp_user'), openssl_decrypt($core->framework->server->getData('ftp_pass'), 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($core->framework->server->getData('encryption_iv'))));
						
							$FTPLocalFile = $saveDir.'save.'.$file;
							$sftp = ssh2_sftp($SFTPConnection);
																
								$stream = fopen("ssh2.sftp://".$sftp."/server/".$directory.$file, 'r');

									if(!$stream){
									
										$errormsg = '<div class="error-box round">Unable to download file for editing.</div>';
									
									}else{
										
										$contents = stream_get_contents($stream);
										fclose($stream);
										
									}
						
						$parName = '(Editing: /'.$_GET['file'].')';
						
						if($errormsg != ''){
						
							$HTML = $errormsg;
						
						}else{
						
							$HTML = '<form action="files.php?do=edit&action=savefile" method="post">
										<textarea name="file_contents" id="live_console">'.$contents.'</textarea>
										<input type="hidden" name="file" value="'.$_GET['file'].'" />
										<center><a href="files.php?dir='.urlencode($directory).'" class="round button blue text-upper" style="padding:0.721em !important;">Back to File Manager</a>&nbsp;&nbsp;<input type="submit" style="font-weight:normal;" value="Save File" id="save_file" class="round blue ic-edit text-upper"/></center>
									</form>';
								
						}
					
				}else{
				
					$HTML = '<div class="error-box round">This type of file cannot be edited via our control panel.</div>';
				
				}
			
			}else{
			
				$HTML = '<div class="error-box round">No file was specified that can be edited.</div>';
			
			}
			
		}else{
		
			/*
			 * Save File
			 */
			if(isset($_GET['action']) && $_GET['action'] == 'savefile'){
			
				if(isset($_POST['file']) && !is_dir($_POST['file']) && file_exists($core->framework->server->getData('path').$_POST['file'])){
				
					if(in_array(pathinfo($_POST['file'], PATHINFO_EXTENSION), $canEdit)){
					
						/*
						 * Begin Advanced Saving
						 */
						$saveDir = '/tmp/'.$core->framework->server->getData('hash').'/';
						
							/*
							 * Check that Secure User DIrectory Exists
							 */
							if(!is_dir($saveDir)){
							
								/*
								 * Make Directory
								 */
								mkdir($saveDir);
							
							}
							
								/*
								 * Create Save File
								 */
								$file = pathinfo($_POST['file'], PATHINFO_BASENAME);
								$directory = dirname($_POST['file']).'/';
								
								/*
								 * Directory Cleaning
								 */
								if($directory == './' || $directory == '.'){ $directory = ''; }
								if(substr($directory, 0, 1) == '/'){ $directory = substr($directory, 1); }
								
								$fp = fopen($saveDir.'save.'.$file , 'w');
								fwrite($fp, $_POST['file_contents']);
								fclose($fp);
								
									/*
									 * Upload Via SFTP
									 */
									$SFTPConnection = ssh2_connect($core->framework->server->getData('ftp_host'), 22);
									ssh2_auth_password($SFTPConnection, $core->framework->server->getData('ftp_user'), openssl_decrypt($core->framework->server->getData('ftp_pass'), 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($core->framework->server->getData('encryption_iv'))));
									
										$FTPLocalFile = $saveDir.'save.'.$file;
										$sftp = ssh2_sftp($SFTPConnection);
																			
											$stream = fopen("ssh2.sftp://".$sftp."/server/".$directory.$file, 'w');
											
												if(!$stream){
												
													$core->framework->page->redirect('files.php?do=edit&file='.urlencode($_POST['file']).'&error='.base64_encode('Unable to connect and upload file.'));
												
												}else{
													
													if(fwrite($stream, file_get_contents($FTPLocalFile))){
													
														fclose($stream);
														unlink($FTPLocalFile);
														$HTML = '<div class="confirmation-box round">File was sucessfully saved. Return to <a href="files.php?do=edit&file='.urlencode($_POST['file']).'">file</a>, or return to <a href="files.php?dir='.urlencode($directory).'">file manager</a>.</div>';
													
													}else{
													
														fclose($stream);
														$core->framework->page->redirect('files.php?do=edit&file='.urlencode($_POST['file']).'&error='.base64_encode('An unknown error occured. Unable to save file.'));
													
													}
												
												}
					
					}else{
					
						$core->framework->page->redirect('files.php?do=edit&file='.urlencode($_POST['file']).'&error='.base64_encode('This type of file cannot be edited via our online file manager. Please use a FTP client.'));
					
					}
				
				}else{
				
					$core->framework->page->redirect('files.php?do=edit&file='.urlencode($_POST['file']).'&error='.base64_encode('The file specified could not be found on the server.'));
				
				}
			
			}
		
		}
	
	}else if($_GET['do'] == 'download'){
	
		if(file_exists($core->framework->server->getData('path').$_GET['file'])){
			
			/*
			 * Download a File
			 */
			header("Pragma: public");
			header("Expires: 0");
			header("Cache-Control: must-revalidate, post-check=0, pre-check=0");
			header("Content-Type: application/force-download");
			header("Content-Description: File Transfer");
			header('Content-Disposition: attachment; filename="'.$_GET['file'].'"');
			header("Content-Length: ".filesize($core->framework->server->getData('path').$_GET['file']));
				
			$core->framework->files->download($core->framework->server->getData('path').$_GET['file']);
			exit();
			
		}
	
	}
//	}else if($_GET['do'] == 'delete'){
//	
//		/*
//		 * Deleting File or Directory
//		 */
//		if($rcon->s->isOnline($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port')) === true){
//			
//			$HTML = '<div class="error-box round">Server must be off in order to delete files via the online file management interface.</div>';
//			
//		}else{
//		
//			$_GET['file'] = str_replace(array(';', '|', '\\', ' ', '..'), '', trim($_GET['file']));
//		
//			if(!is_dir($core->framework->server->getData('path').$_GET['file']) && file_exists($core->framework->server->getData('path').$_GET['file'])){
//			
//				$isFile = true;
//				
//				/*
//				 * File Info
//				 */
//				$file = pathinfo($_GET['file'], PATHINFO_BASENAME);
//				$directory = dirname($_GET['file']).'/';
//				
//				/*
//				 * Directory Cleaning
//				 */
//				if($directory == './' || $directory == '.'){ $directory = ''; }
//				if(substr($directory, 0, 1) == '/'){ $directory = substr($directory, 1); }
//				
//				$FTPRemoteFile = "/server/".$directory.$file;
//				
//			}else{
//			
//				$isFile = false;
//				$FTPRemoteFile = rtrim($_GET['file'], '/');
//			
//			}
//			
//			/*
//			 * Delete File Via FTP
//			 */				
//			$SFTPConnection = ssh2_connect($core->framework->server->getData('ftp_host'), 22);
//			ssh2_auth_password($SFTPConnection, $core->framework->server->getData('ftp_user'), openssl_decrypt($core->framework->server->getData('ftp_pass'), 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($core->framework->server->getData('encryption_iv'))));
//			
//			$sftp = ssh2_sftp($SFTPConnection);
//			
//				if($isFile === true){
//				
//					if(ssh2_sftp_unlink($sftp, $FTPRemoteFile)){
//					
//						$HTML = '<div class="confirmation-box round">File was sucessfully deleted. Return to <a href="files.php">file manager</a>.</div>';
//					
//					}else{
//					
//						$HTML = '<div class="error-box round">Unable to delete file.</div>';
//					
//					}
//					
//				}else{
//				
//					if(is_dir('/srv/servers/'.$core->framework->server->getData('name').'/server/'.$FTPRemoteFile)){
//					$nodeSQLConnect = $mysql->prepare("SELECT * FROM `nodes` WHERE `node_name` = ? LIMIT 1");
//					$nodeSQLConnect->execute(array($core->framework->server->getData('node')));
//					
//					$row = $nodeSQLConnect->fetch();
//					
//					$con = ssh2_connect($row['ip'], 22);
//					ssh2_auth_password($con, $row['user'], $row['password']);
//				
//					$deletePath = escapeshellarg($row['data_dir'].$core->framework->server->getData('name').'/server/'.$FTPRemoteFile);
//				
//					$s = ssh2_exec($con, 'echo "'.$row['password'].'" | sudo -S su - root -c "cd '.$row['data_dir'].$core->framework->server->getData('name').'/server; rm -rf '.$deletePath.'"');
//					stream_set_blocking($s, true);
//
//					}else{
//					
//						$HTML = '<div class="error-box round">Not a known directory.</div>';
//					
//					}
//				
//				}
//		
//		}
//		
//	}

}

?>
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title><?php echo $core->framework->settings->get('company_name'); ?> - Server Console</title>
	
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
			<div class="side-content fr">
				<div class="content-module">
					<div class="content-module-heading cf">
						<h3 class="fl">File Manager <?php echo $parName; ?></h3>
					</div> <!-- end content-module-heading -->
					<div class="content-module-main" id="server_info">
						<?php echo $HTML; ?>
					</div> <!-- end content-module-main -->
				</div>
			</div>
		</div>
	</div>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.3 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>