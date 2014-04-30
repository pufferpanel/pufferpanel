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
require_once('../../../core/framework/framework.core.php');

$filesIncluded = true;

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false){

	exit('Invalid authentication.');
    
}

$canEdit = array('txt', 'yml', 'log', 'conf', 'html', 'json', 'properties', 'props', 'cfg', 'lang');

if(isset($_POST['dir'])){
	rtrim($_POST['dir'], '/');
    $_POST['dir'] = str_replace('..', '', urldecode($_POST['dir']));
}

/*
 * Display File Manager Overview Page
 */
if(isset($_POST['dir']) && !empty($_POST['dir'])){

    /*
     * Check First Character
     */
    if(substr($_POST['dir'], 0, 1) == '/')
        $_POST['dir'] = substr($_POST['dir'], 1);
                        
    if($_POST['dir'] == '/')
        $displayFolders = '';
        
    $_POST['dir'] = $_POST['dir'].'/';
	
	$connection = $core->ssh->generateSSH2Connection($core->server->getData('id'), false, true);
	
	$sftp = ssh2_sftp($connection);
	
	/*
	 * Validate Directory
	 */
	if(!is_dir("ssh2.sftp://$sftp/server/".$_POST['dir'])){
	    
	    $core->log->getUrl()->addLog(4, 0, array('system.path_missing', 'The directory for `'.$core->server->getData('name').'` could not be found by the file manager.'));
	    exit('<div class="error-box round">'.$_l->tpl('node.files.ajax.no_dir').'</a>');
	
	}
	
	$displayFolders = '';
	$displayFiles = '';
	
	$handle = opendir("ssh2.sftp://$sftp/server/".$_POST['dir']);
	
	$entries = array();
	while (false != ($entries[] = readdir($handle)));
	sort($entries);
	closedir($handle);
		
	/*
	 * Inside a Directory
	 */
	$goBack = explode('/', $_POST['dir']);
	
	    $displayFolders = ' <tr>
	                            <td><i class="fa fa-folder-open">&nbsp</i></td>
	                            <td><a href="index.php" class="load_new">&larr;</a></a></td>
	                            <td></td>
	                            <td></td>
	                            <td></td>
	                         </tr>';
	
	    if(array_key_exists(1, $goBack) && !empty($goBack[1])){
	    
	        unset($goBack[count($goBack) - 2]);
	        $previousDir = rtrim(implode('/', $goBack), '/');
	        
	        $displayFolders .= '<tr>
	                    <td><i class="fa fa-folder-open">&nbsp</i></td>
	                    <td><a href="index.php?dir='.$previousDir.'" class="load_new">&larr; '.$previousDir.'</a></a></td>
	                    <td></td>
	                    <td></td>
	                    <td></td>
	                </tr>';
	        
	    }
	
	foreach($entries as $entry){
		    
	    if($entry == '.' || $entry == '..' || $entry == '')
	    	echo null;
	    else {
	    
	        /*
	         * Get Stats on File
	         */
	        $stat = ssh2_sftp_stat($sftp, '/server/'.$_POST['dir'].$entry);
	        
	        /*
	         * Iterate into HTML Variable
	         */
	        if(is_dir("ssh2.sftp://$sftp/server/".$_POST['dir'].$entry)){
	        
	            $displayFolders .= 	'<tr>
	                                    <td><i class="fa fa-folder-open">&nbsp</i></td>
	                                    <td><a href="index.php?dir='.urlencode($_POST['dir'].$entry).'" class="load_new">'.$entry.'</a></td>
	                                    <td>'.$core->files->formatSize($stat['size']).'</td>
	                                    <td>'.date('m/d/y H:i:s', $stat['mtime']).'</td>
	                                    <td class="center"></td>
	                                </tr>';
	        
	        }else{
	        
	            $url = (in_array(pathinfo($entry, PATHINFO_EXTENSION), $canEdit)) ? '<a href="edit.php?file='.urlencode($_POST['dir'].$entry).'" class="edit_file">'.$entry.'</a>' : $url = $entry;
	               
	            $displayFiles .= 	'<tr>
	                                    <td><i class="fa fa-file-text"></i></td>
	                                    <td>'.$url.'</td>
	                                    <td>'.$core->files->formatSize($stat['size']).'</td>
	                                    <td>'.date('m/d/y H:i:s', $stat['mtime']).'</td>
	                                    <td style="text-align:center;"><a href="index.php?do=download&file='.urlencode($_POST['dir'].$entry).'"><i class="fa fa-download"></i></a></td>
	                                </tr>';
	        
	        }

		}
	            
    }
    
}else{
	
	$connection = $core->ssh->generateSSH2Connection($core->server->getData('id'), false, true);

	$sftp = ssh2_sftp($connection);
	
		if(!$sftp)
			exit('<div class="alert alert-danger">'.$_l->tpl('node.files.ajax.no_dl').'</div>');
	
	$displayFolders = '';
	$displayFiles = '';
	
	$handle = opendir("ssh2.sftp://$sftp/server/");
	
	$entries = array();
	while (false != ($entries[] = readdir($handle)));
	sort($entries);
	closedir($handle);
	
	foreach($entries as $entry){
	    
	    if($entry == '.' || $entry == '..' || $entry == '')
	    	echo null;
	    else {
	    
	        /*
	         * Get Stats on File
	         */
	        $stat = ssh2_sftp_stat($sftp, '/server/'.$entry);
	        
	        /*
	         * Iterate into HTML Variable
	         */
	        if(is_dir("ssh2.sftp://$sftp/server/".$entry)){
	        
	            $displayFolders .= 	'<tr>
	                                    <td><i class="fa fa-folder-open">&nbsp</i></td>
	                                    <td><a href="index.php?dir='.urlencode($entry).'" class="load_new">'.$entry.'</a></td>
	                                    <td>'.$core->files->formatSize($stat['size']).'</td>
	                                    <td>'.date('m/d/y H:i:s', $stat['mtime']).'</td>
	                                    <td class="center"></td>
	                                </tr>';
	        
	        }else{
	        
	            $url = (in_array(pathinfo($entry, PATHINFO_EXTENSION), $canEdit)) ? '<a href="edit.php?file='.urlencode($entry).'" class="edit_file">'.$entry.'</a>' : $url = $entry;
	               
	            $displayFiles .= 	'<tr>
	                                    <td><i class="fa fa-file-text"></i></td>
	                                    <td>'.$url.'</td>
	                                    <td>'.$core->files->formatSize($stat['size']).'</td>
	                                    <td>'.date('m/d/y H:i:s', $stat['mtime']).'</td>
	                                    <td style="text-align:center;"><a href="index.php?do=download&file='.urlencode($entry).'"><i class="fa fa-download"></i></a></td>
	                                </tr>';
	        
	        }

		}
	            
    }

}

/*
 * Setup Basic Display
 */
echo '<table class="table table-striped table-bordered table-hover">
            <thead>
                <tr>
                    <th style="width:5%;text-align:center;"></th>
                    <th style="width:45%">'.$_l->tpl('node.files.name').'</th>
                    <th style="width:20%">'.$_l->tpl('node.files.size').'</th>
                    <th style="width:20%">'.$_l->tpl('node.files.modified').'</th>
                    <th style="width:10%;text-align:center;">'.$_l->tpl('string.options').'</th>
                </tr>
            </thead>
            <tbody>
                '.$displayFolders.$displayFiles.'
            </tbody>
        </table>';

