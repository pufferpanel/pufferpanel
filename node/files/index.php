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

if(isset($_GET['file']))
    $_GET['file'] = str_replace('..', '', urldecode($_GET['file']));

if(isset($_GET['dir']))
    $_GET['dir'] = str_replace('..', '', urldecode($_GET['dir']));
	
if(isset($_GET['do']) && $_GET['do'] == 'download'){

    $path = $core->framework->server->nodeData('server_dir').$core->framework->server->getData('ftp_user').'/server/';
    if(file_exists($path.$_GET['file'])){
        
        /*
         * Download a File
         */
        header("Pragma: public");
        header("Expires: 0");
        header("Cache-Control: must-revalidate, post-check=0, pre-check=0");
        header("Content-Type: application/force-download");
        header("Content-Description: File Transfer");
        header('Content-Disposition: attachment; filename="'.$_GET['file'].'"');
        header("Content-Length: ".filesize($path.$_GET['file']));
            
        $core->framework->files->download($path.$_GET['file']);
        exit();
        
    }

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
					<li><a href="index.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> File Manager</a></li>
					<li><a href="../backup.php"><i class="fa fa-angle-double-right pull-right menu-arrows"></i> Backup Manager</a></li>
				</ul>
			</div>
			<div class="side-content fr">
				<div class="content-module">
					<div class="content-module-heading cf">
						<h3 class="fl">File Manager <i class="fa fa-cog fa-spin" id="loading_dir"></i></h3>
					</div> <!-- end content-module-heading -->
					<div class="content-module-main" id="load_files">
                        
					</div> <!-- end content-module-main -->
				</div>
			</div>
		</div>
	</div>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.4 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
    <script type="text/javascript">
    $(document).ready(function(){
       firstLoad(); 
    });
        
        $.urlParam = function(name, url){
			 var results = new RegExp('[\\?&]' + name + '=([^&#]*)').exec(decodeURIComponent(url));
			 if (results==null){
			     return null;
             }else{
			     return results[1] || 0;
			 }
        }
        
        function newLoad(){
            
            $("a.load_new").click(function(event){
                event.preventDefault();
                
                $("#loading_dir").fadeIn(200);
                if($.urlParam('dir', $(this).attr("href")) != null){
                
                    var dir = $.urlParam('dir', $(this).attr("href"));
                    $.ajax({
                        type: "POST",
                        url: '../core/ajax/files/list_dir.php',
                        data: {'dir': dir},
                        success: function(data) {
                            $("#load_files").slideUp(function(){
                                $("#load_files").html(data);
                                $("#load_files").slideDown();
                                $("#loading_dir").fadeOut(200);
                                newLoad();
                            });
                        }
                    });
                    
                }else{
                 
                    $.ajax({
                    type: "POST",
                    url: '../core/ajax/files/list_dir.php',
                    data: {},
                    success: function(data) {
                        $("#load_files").slideUp(function(){
                            $("#load_files").html(data);
                            $("#load_files").slideDown();
                            $("#loading_dir").fadeOut(200);
                            newLoad();
                        });
                    }
                });
                    
                }
                
            });
            
            $("a.edit_file").click(function(event){
             
                event.preventDefault();
                var file = $.urlParam('file', $(this).attr("href"));
                $().redirect('edit.php', {'file': file});
                
            });
            
        }
        
        function firstLoad() {
            
            $("#loading_dir").fadeIn(200);
            $.ajax({
                type: "POST",
                url: '../core/ajax/files/list_dir.php',
                data: {},
                success: function(data) {
                    $("#load_files").slideUp(function(){
                        $("#load_files").html(data);
                        $("#load_files").slideDown();
                        $("#loading_dir").fadeOut(200);
                        newLoad();
                    });
                }
            });
            
        }
    </script>
</body>
</html>