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
require_once('../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), $core->framework->auth->getCookie('pp_server_hash')) === false){

	$core->framework->page->redirect($core->framework->settings->get('master_url').'index.php');
	exit();
    
}

$canEdit = array('txt', 'yml', 'log', 'conf', 'html', 'json', 'properties', 'props', 'cfg', 'lang');

if(isset($_POST['file']))
    $_POST['file'] = str_replace('..', '', urldecode($_POST['file']));

if(isset($_POST['dir']))
    $_POST['dir'] = str_replace('..', '', urldecode($_POST['dir']));

$parName = '(Editing: /'.$_POST['file'].')';
$path = $core->framework->server->nodeData('server_dir').$core->framework->server->getData('ftp_user').'/server/';

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
    <script type="text/javascript" src="<?php echo $core->framework->settings->get('master_url'); ?>assets/javascript/jquery.redirect.min.js"></script>
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
					<li><a href="../index.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Overview</a></li>
					<li><a href="../console.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Live Console</a></li>
					<li><a href="../settings.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Server Settings</a></li>
					<li><a href="../plugins.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Server Plugins</a></li>
					<li><a href="../files/index.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> File Manager</a></li>
					<li><a href="../backup.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Backup Manager</a></li>
				</ul>
			</div>
			<div class="side-content fr">
				<div class="content-module">
					<div class="content-module-heading cf">
						<h3 class="fl"><?php echo $parName; ?></h3>
					</div> <!-- end content-module-heading -->
					<div class="content-module-main" id="server_info">
                        <p id="save_inp" style="text-align:center;display:none;"><i class="fa fa-cog fa-spin" id="loading_dir"></i> Saving File...</p>
                        <span id="save_status" style="display:none;"></span>
						<?php                        
                        if(!isset($_GET['error'])){
                            
                            if(isset($_POST['file']) && !is_dir($_POST['file']) && file_exists($path.$_POST['file'])){
                            
                                if(in_array(pathinfo($path.$_POST['file'], PATHINFO_EXTENSION), $canEdit)){
                                            
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
                                    $file = pathinfo($path.$_POST['file'], PATHINFO_BASENAME);
                                    $directory = dirname($_POST['file']).'/';
                                    
                                    /*
                                     * Directory Cleaning
                                     */
                                    if($directory == './' || $directory == '.')
                                        $directory = '';
                                    
                                    if(substr($directory, 0, 1) == '/')
                                        $directory = substr($directory, 1);
                                                                
                                        /*
                                         * Download Via SFTP
                                         */
                                        $SFTPConnection = ssh2_connect($core->framework->server->getData('ftp_host'), 22);
                                        ssh2_auth_password($SFTPConnection, $core->framework->server->getData('ftp_user'), openssl_decrypt($core->framework->server->getData('ftp_pass'), 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($core->framework->server->getData('encryption_iv'))));
                                        
                                            $FTPLocalFile = $saveDir.'save.'.$file;
                                            $sftp = ssh2_sftp($SFTPConnection);
                                                                                
                                                $stream = fopen("ssh2.sftp://".$sftp."/server/".$directory.$file, 'r');
                            
                                                    if(!$stream){
                                                    
                                                        $core->framework->page->redirect('edit.php?error='.base64_encode('<div class="error-box round">Unable to download file for editing.</div>'));
                                                        exit();
                                                    
                                                    }else{
                                                        
                                                        $contents = stream_get_contents($stream);
                                                        fclose($stream);
                                                        
                                                    }
                                        
                                        echo '<form method="post" id="editing_file">
                                                        <textarea name="file_contents" id="live_console">'.$contents.'</textarea>
                                                        <input type="hidden" name="file" value="'.$_POST['file'].'" />
                                                        <center><a href="index.php" class="round button blue text-upper" style="padding:0.721em !important;">Back to File Manager</a>&nbsp;&nbsp;<input type="submit" style="font-weight:normal;" value="Save File" id="save_file" class="round blue ic-edit text-upper"/></center>
                                            </form>';
                                    
                                }else{
                                
                                    echo '<div class="error-box round">This type of file cannot be edited via our control panel.</div>';
                                
                                }
                            
                            }else{
                            
                                echo '<div class="error-box round">No file was specified that can be edited.</div>';
                            
                            }
                            
                        }else{
                            
                            echo (isset($_GET['error'])) ? '<div class="error-box round">'.base64_decode($_GET['error']).' You can return to the file manager by <a href="index.php">clicking here</a>.</div>' : '';
                            
                        }
                        ?>
					</div> <!-- end content-module-main -->
				</div>
			</div>
		</div>
	</div>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.4 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
    <script type="text/javascript">
        $("input#save_file").click(function(event){
            event.preventDefault();
            var file = $("input[name='file']").val();
            var content = $("#live_console").val();
            $("#save_inp").slideDown(200);
            $.ajax({
                type: "POST",
                url: '../core/ajax/files/save.php',
                data: {'file': file, 'file_contents': content},
                    success: function(data) {
                        $("#save_status").hide();
                        $("#save_inp").slideUp(200);
                        $("#save_status").html(data);
                        $("#save_status").slideDown().delay(2500).slideUp();
                     }
            });
       });
    </script>
</body>
</html>