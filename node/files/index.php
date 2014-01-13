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
	<?php include('../assets/include/header.php'); ?>
	<script type="text/javascript" src="<?php echo $core->framework->settings->get('master_url'); ?>assets/javascript/jquery.redirect.min.js"></script>
	<title>PufferPanel - Manage Your Server</title>
</head>
<body>
	<div class="container">
		<?php include('../assets/include/navbar.php'); ?>
		<div class="row">
			<div class="col-3">
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Account Actions</strong></a>
					<a href="<?php echo $core->framework->settings->get('master_url'); ?>account.php" class="list-group-item">Settings</a>
					<a href="<?php echo $core->framework->settings->get('master_url'); ?>servers.php" class="list-group-item">My Servers</a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong>Server Actions</strong></a>
					<a href="../index.php" class="list-group-item">Overview</a>
					<a href="../console.php" class="list-group-item">Live Console</a>
					<a href="../settings.php" class="list-group-item">Server Settings</a>
					<a href="../plugins/index.php" class="list-group-item">Server Plugins</a>
					<a href="index.php" class="list-group-item active">File Manager <i class="fa fa-spinner fa fa-spin" id="loading_dir" style="display:none;"></i></a>
					<a href="../backup.php" class="list-group-item">Backup Manager</a>
				</div>
			</div>
			<div class="col-9" id="load_files">
			</div>
		</div>
		<div class="footer">
			<?php include('../assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
	$(document).ready(function(){ 
	    $.urlParam = function(name, url){
			 var results = new RegExp('[\\?&]' + name + '=([^&#]*)').exec(decodeURIComponent(url));
			 if (results==null){
			     return null;
	         }else{
			     return results[1] || 0;
			 }
	    }
	    var doneLoad = false;
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
	        if($.urlParam('dir', $(location).attr('href')) != null && doneLoad === false){
	        	var dir = $.urlParam('dir', $(location).attr('href'));
				$("#loading_dir").fadeIn(200);
	        	$.ajax({
	        	    type: "POST",
	        	    url: '../core/ajax/files/list_dir.php',
	        	    data: {'dir': dir},
	        	    success: function(data) {
	        	        $("#load_files").slideUp(function(){
	        	            $("#load_files").html(data);
	        	            $("#load_files").slideDown();
	        	            $("#loading_dir").fadeOut(200);
	        	            doneLoad = true;
	        	            newLoad();
	        	        });
	        	    }
	        	});
	        }else{
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
		                    doneLoad = true;
		                    newLoad();
		                });
		            }
		        });
	        }
	    }
	    firstLoad();
	});
	</script>
</body>
</html>