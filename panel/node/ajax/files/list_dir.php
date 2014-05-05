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
require_once('../../../../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false)
	exit('Invalid authentication passed to script.');

/*
 * Set Defaults
 */
$displayFolders = array();
$displayFiles = array();
$entries = array();
$previousDir = array();

if(isset($_POST['dir']))
    $_POST['dir'] = str_replace('..', '', urldecode(rtrim($_POST['dir'], '/')));
else
	$_POST['dir'] = null;

/*
 * Gather Files and Folders
 */
$connection = $core->ssh->generateSSH2Connection($core->server->getData('id'), false, true);
$sftp = ssh2_sftp($connection);

	if(!$sftp)
		exit('<div class="alert alert-danger">'.$_l->tpl('node_files_ajax_no_dl').'</div>');

/*
 * Parse Through Files
 */
$handle = opendir("ssh2.sftp://$sftp/server/".$_POST['dir']);
while (false != ($entries[] = readdir($handle)));
sort($entries);
closedir($handle);
 
/*
 * Handle Directory
 */
if(isset($_POST['dir']) && !empty($_POST['dir'])){

	/*
	 * In dir, show first arrow in display
	 */
	$previousDir['first'] = true;
	
    /*
     * Check First Character
     */
    if(substr($_POST['dir'], 0, 1) == '/')
        $_POST['dir'] = substr($_POST['dir'], 1);
	else
		$_POST['dir'] = $_POST['dir'].'/';

}
	
/*
 * Validate Directory
 */
if(!is_dir("ssh2.sftp://$sftp/server/".$_POST['dir']))
    exit('<div class="error-box round">'.$_l->tpl('node_files_ajax_no_dir').'</a>');
	
/*
 * Inside a Directory
 */
$goBack = explode('/', $_POST['dir']);
if(array_key_exists(1, $goBack) && !empty($goBack[1])){

	/*
	 * Do we show previous-dir arrow?
	 */
	if(!empty($goBack[1]))
		$previousDir['show'] = true;
		
    unset($goBack[count($goBack) - 2]);
    $previousDir['link'] = rtrim(implode('/', $goBack), '/');
    
}

/*
 * Setting More Variables
 */
if(array_key_exists('link', $previousDir) && strpos(rtrim($previousDir['link'], '/'), '/'))
	$previousDir['link_show'] = end(explode('/', $previousDir['link']));
elseif(array_key_exists('link', $previousDir))
	$previousDir['link_show'] = $previousDir['link'];

/*
 * Loop Through
 */	
foreach($entries as $entry) {
	    
    if($entry != '.' && $entry != '..' && $entry != ''){
    
        /*
         * Get Stats on File
         */
        $stat = ssh2_sftp_stat($sftp, '/server/'.$_POST['dir'].$entry);
        
		 /*
         * Iterate into Arrays
         */
        if(is_dir("ssh2.sftp://$sftp/server/".$_POST['dir'].$entry)){
        
        	$displayFolders = array_merge($displayFolders, array(array(
        		"entry" => $entry,
        		"directory" => $_POST['dir'],
        		"size" => $core->files->formatSize($stat['size']),
        		"date" => date('m/d/y H:i:s', $stat['mtime'])
        	)));

        }else{
        
        	$displayFiles = array_merge($displayFiles, array(array(
        		"entry" => $entry,
        		"directory" => $_POST['dir'],
        		"extension" => pathinfo($entry, PATHINFO_EXTENSION),
        		"size" => $core->files->formatSize($stat['size']),
        		"date" => date('m/d/y H:i:s', $stat['mtime'])
        	)));
        
        }
        
	}
            
}

/*
 * Render Page
 */
echo $twig->render(
		'node/files/ajax/list_dir.html', array(
			'files' => $displayFiles,
			'folders' => $displayFolders,
			'extensions' => array('txt', 'yml', 'log', 'conf', 'html', 'json', 'properties', 'props', 'cfg', 'lang'),
			'directory' => $previousDir
	));
?>