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

$backupStatus = null;
$returnBackups = null;
$backupList = null;
if(isset($_GET['do']) && $_GET['do'] == 'create'){

	if(isset($_POST['dbu'])){

        /*
		 * Does User have Space!
		 */
		$files = glob($core->framework->server->nodeData('backup_dir').$core->framework->server->getData('name').'/*', GLOB_MARK);
		
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
					$backupStatus .= '<div class="alert alert-warning">You have reached your allotment of backup space. Please delete older backups before continuing.</div>';
					
				}
				
				if($totalFiles >= $core->framework->server->getData('backup_file_limit')){
				
					$continue = false;
					$backupStatus .= '<div class="alert alert-warning">You have reached your limit for stored backups. Please delete older ones before continuing.</div>';
				
				} 
				
		
		/*
		 * Do Backup
		 */
		if($continue === true) {
		
            $user = $core->framework->server->getData('name');
            $filename = str_replace(" ", "", stripslashes(trim($_POST['backup_name'])));
            $filename = ltrim($filename, '/');
            $filename = rtrim($filename, '/');
            $filename = rtrim($filename, '.');
            $filename = preg_replace('/[^A-Za-z0-9_-]/', '', $filename);
            $filename = preg_replace('/(\.){2,}/', '.', $filename);
            $filename = $filename.'_'.$core->framework->auth->keygen('10');
            $saveOptions = (isset($_POST['save_options'])) ? true : false;
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
                
                if($data == '')
                    $data = '*';
                
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
                    $con = ssh2_connect($core->framework->server->nodeData('sftp_ip'), 22);
                    ssh2_auth_password($con, $core->framework->server->nodeData('username'), openssl_decrypt($core->framework->server->nodeData('password'), 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($core->framework->server->nodeData('encryption_iv'))));
                    
                    /*
                     * If server is running we need to do this differently.
                     */
                    if($core->framework->rcon->online($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port')) === true){
                    
                        $stream = ssh2_exec($con, 'cd /srv/scripts; sudo ./send_command.sh '.$core->framework->server->getData('ftp_user').' "save-all"', true);
                        $errorStream = ssh2_fetch_stream($stream, SSH2_STREAM_STDERR);
                        
                        stream_set_blocking($errorStream, true);
                        stream_set_blocking($stream, true);
                        
                        $isError = stream_get_contents($errorStream);
                        if(!empty($isError))
                        	echo $isError;
                        
                        fclose($errorStream);
                        fclose($stream);
                        
                        sleep(2);
                        
                        $stream = ssh2_exec($con, 'cd /srv/scripts; sudo ./send_command.sh '.$core->framework->server->getData('ftp_user').' "save-off"', true);
                        $errorStream = ssh2_fetch_stream($stream, SSH2_STREAM_STDERR);
                        
                        stream_set_blocking($errorStream, true);
                        stream_set_blocking($stream, true);
                        
                        $isError = stream_get_contents($errorStream);
                        if(!empty($isError))
                        	echo $isError;
                        
                        fclose($errorStream);
                        fclose($stream);
                    
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
                                                                    
                                    $stream = ssh2_exec($con, 'cd /srv/scripts; sudo ./backup_server.sh '.$core->framework->server->getData('ftp_user').' "'.$filename.' '.$core->framework->server->nodeData('node').' '.$core->framework->server->getData('hash').' '.$backupToken.'" "*" ""', true);
                                    $errorStream = ssh2_fetch_stream($stream, SSH2_STREAM_STDERR);
                                    
                                    stream_set_blocking($errorStream, true);
                                    stream_set_blocking($stream, true);
                                    
                                    $isError = stream_get_contents($errorStream);
                                    if(!empty($isError))
                                    	echo $isError;
                                    
                                    fclose($errorStream);
                                    fclose($stream);
    
                                }else{
                                                                                        
                                    $stream = ssh2_exec($con, 'cd /srv/scripts; sudo ./backup_server.sh '.$core->framework->server->getData('ftp_user').' "'.$filename.' '.$core->framework->server->nodeData('node').' '.$core->framework->server->getData('hash').' '.$backupToken.'" "'.escapeshellcmd(str_replace(",", " ", $backup)).'" "'.escapeshellcmd(str_replace(",", " ", $skip)).'"', true);
                                    $errorStream = ssh2_fetch_stream($stream, SSH2_STREAM_STDERR);
                                    
                                    stream_set_blocking($errorStream, true);
                                    stream_set_blocking($stream, true);
                                    
                                    $isError = stream_get_contents($errorStream);
                                    if(!empty($isError))
                                    	echo $isError;
                                    
                                    fclose($errorStream);
                                    fclose($stream);
                                                                        
                                }
                                                                
                                $core->framework->log->getUrl()->addLog(0, 1, array('user.backup_started', 'A backup was started for `'.$core->framework->server->getData('name').'`.'));
                                $backupStatus .= '<div class="alert alert-success">Your backup has been started. Backups will appear at the bottom of this page when they are finished.</div>';
                        
                        }else{
                        
                            $stream = ssh2_exec($con, 'cd /srv/scripts; sudo ./send_command.sh '.$core->framework->server->getData('ftp_user').' "save-on"', true);
                            $errorStream = ssh2_fetch_stream($stream, SSH2_STREAM_STDERR);
                            
                            stream_set_blocking($errorStream, true);
                            stream_set_blocking($stream, true);
                            
                            $isError = stream_get_contents($errorStream);
                            if(!empty($isError))
                            	echo $isError;
                            
                            fclose($errorStream);
                            fclose($stream);
                            
                            $core->framework->log->getUrl()->addLog(2, 1, array('system.sql_error', 'A backup was started for `'.$core->framework->server->getData('name').'` but failed due to a MySQL error.'));
                            
                            $backupStatus .= '<div class="alert alert-danger">A MySQL error was encountered and to ensure the integrity of your server we have stopped this backup process. Please try again or contact support if this error continues to occur.</div>';
                        
                        }
            
        }
	
	}else{
	
		$backupStatus = "<div class='alert alert-danger'>We were unable to preform the requested operation. This was due to the lack of any data sent to the function. Please retry your request.</div>";
	
	}

}

$query = $mysql->prepare("SELECT * FROM `backups` WHERE `server` = ? AND `timeend` IS NOT NULL ORDER BY `id` DESC");
$query->execute(array($core->framework->server->getData('hash')));

	if($query->rowCount() > 0){
			
			$returnBackups = '';
			while($row = $query->fetch()){
						
				$stat = stat($core->framework->server->nodeData('backup_dir').$core->framework->server->getData('ftp_user').'/'.$row['file_name'].'.tar.gz');
				
				$md5 = (strlen($row['md5']) > 10) ? substr($row['md5'], 0, 7).'...' : $row['md5'];
				$sha1 = (strlen($row['sha1']) > 10) ? substr($row['sha1'], 0, 7).'...' : $row['sha1'];
				$fileName = (strlen($row['file_name']) > 30) ? substr($row['file_name'], 0, 27).'...' : $row['file_name'];
				
				$returnBackups .= '<tr>
										<td style="text-align:center;"><a href="core/ajax/backup/download.php?id='.$row['id'].'"><i class="fa fa-download"></i></a></td>
										<td><a href="#" onclick="return false;" data-toggle="popover" data-content="'.$row['file_name'].'" data-original-title="File Name">'.$fileName.'</a></td>
										<td><a href="#" onclick="return false;" data-toggle="popover" data-content="'.$row['md5'].'" data-original-title="MD5 Checksum">'.$md5.'</a></td>
										<td><a href="#" onclick="return false;" data-toggle="popover" data-content="'.$row['sha1'].'" data-original-title="SHA1 Checksum">'.$sha1.'</a></td>
										<td>'.date('M n, Y \a\t g:ia', $row['timeend']).'</td>
										<td>'.$core->framework->files->formatSize($stat['size']).'</td>
									</tr>';
			
			}
				
	}else{
	
		$backupList = '<div class="alert alert-info">You don\'t have any server backups avaliable. You should create one right now in order to protect and ensure the integrity of your server.</div>';
	
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
	<?php include('assets/include/header.php'); ?>
	<title>PufferPanel - Manage Your Server</title>
</head>
<body>
	<div class="container">
		<?php include('assets/include/navbar.php'); ?>
		<div class="row">
			<div class="col-3">
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Account Actions</strong></a>
					<a href="<?php echo $core->framework->settings->get('master_url'); ?>account.php" class="list-group-item">Settings</a>
					<a href="<?php echo $core->framework->settings->get('master_url'); ?>servers.php" class="list-group-item">My Servers</a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Server Actions</strong></a>
					<a href="index.php" class="list-group-item">Overview</a>
					<a href="console.php" class="list-group-item">Live Console</a>
					<a href="files/index.php" class="list-group-item">File Manager</a>
					<a href="backup.php" class="list-group-item active">Backup Manager</a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Server Settings</strong></a>
					<a href="properties.php" class="list-group-item">Server Properties</a>
					<a href="settings.php" class="list-group-item">Modpack Management</a>
					<a href="plugins/index.php" class="list-group-item">Server Plugins</a>
				</div>
			</div>
			<div class="col-9">
				<?php echo $backupStatus; ?>
				<h4 class="nopad">Generate New Backup</h4><hr />
				<form action="backup.php?do=create" method="post">
					<fieldset>
						<div class="form-group">
							<label for="backup_name" class="control-label">Backup Name</label>
							<div>
								<input type="text" class="form-control" name="backup_name" value="backup-<?php echo date('d-n-Y_\a\t_H-i-s', time()); ?>" />
							</div>
						</div>
						<div class="form-group">
							<label for="backup_name" class="control-label">Backup Parameters</label>
							<div>
								<textarea class="form-control" name="backup_params" rows="10"><?php echo $backupParams; ?></textarea>
							</div>
						</div>
						<div class="form-group">
							<div class="checkbox">
								<label>
									<input type="checkbox" name="save_options" value="1"> Save Backup Parameters
								</label>
							</div>
						</div>
						<div class="form-group">
							<div>
								<input type="submit" name="dbu" value="Generate Backup" class="btn btn-primary"/>
								<button class="btn btn-default btn-sm" data-toggle="modal" data-target="#bExample">Example Parameters</button>
							</div>
						</div>
					</fieldset>
				</form><br />
				<h4>Current Stored Backups</h4><hr />
					<table class="table table-striped table-bordered table-hover">
						<thead>
							<tr>
								<th style="width:5%"></th>
								<th style="width:35%">Name</th>
								<th style="width:10%">MD5</th>
								<th style="width:10%">SHA1</th>
								<th style="width:25%">Date</th>
								<th style="width:15%">Size</th>
							</tr>
						</thead>
						<tbody>
							<?php echo $returnBackups; ?>
						</tbody>
					</table>
					<?php echo $backupList; ?>
			</div>
			<div class="modal fade" id="bExample" tabindex="-1" role="dialog" aria-labelledby="BackupExample" aria-hidden="true">
				<div class="modal-dialog">
					<div class="modal-content">
						<div class="modal-header">
							<button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
							<h4 class="modal-title" id="BackupExample">Backup Example</h4>
						</div>
						<div class="modal-body">
							<pre>#Include these folders:<br />plugins<br />world<br />world_nether<br />world_the_end<br /><br />#Also include specific files not in those directories:<br />server.jar<br />server.properties<br />somefolder/t/test.txt<br /><br />#Exclude DynMap and Coreprotect:<br />!plugins/dynmap<br />!plugins/coreprotect<br /><br />#Exclue all plugin jars.<br />!plugins/*.jar</pre>
						</div>
						<div class="modal-footer">
							<button type="button" class="btn btn-primary" data-dismiss="modal">Close</button>
						</div>
					</div>
				</div>
			</div>
		</div>
		<div class="footer">
			<?php include('assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$('a[data-toggle="popover"]').popover({'placement' : 'top'});
	</script>
</body>
</html>