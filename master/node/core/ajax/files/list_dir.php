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

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), $core->framework->auth->getCookie('pp_server_hash')) === false){

	exit('Invalid authentication.');
    
}

$canEdit = array('txt', 'yml', 'log', 'conf', 'html', 'json', 'properties', 'props', 'cfg', 'lang');
$parName = '';

if(isset($_POST['dir']))
    $_POST['dir'] = str_replace('..', '', urldecode($_POST['dir']));

/*
 * Display File Manager Overview Page
 */
$path = $core->framework->server->nodeData('server_dir').$core->framework->server->getData('ftp_user').'/server/';
if(isset($_POST['dir']) && !empty($_POST['dir'])){

    /*
     * Check First Character
     */
    if(substr($_POST['dir'], 0, 1) == '/')
        $_POST['dir'] = substr($_POST['dir'], 1);

    /*
     * Validate Directory
     */
    if(!is_dir($path.$_POST['dir'])){
        
        $core->framework->log->getUrl()->addLog(4, 0, array('system.path_missing', 'The directory for `'.$core->framework->server->getData('name').'` could not be found by the file manager.'));
        exit('<div class="error-box round">Unable to locate request directory. Error logged.</a>');
    
    }
    
    $parName = '(Viewing: /'.$_POST['dir'].')';
    
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
    
        if(count($goBack) > 2){
        
            unset($goBack[count($goBack) - 2]);
            $previousDir = implode('/', $goBack);
            
            $displayFolders .= '<tr>
                        <td><i class="fa fa-folder-open">&nbsp</i></td>
                        <td><a href="index.php?dir='.$previousDir.'" class="load_new">&larr; '.$previousDir.'</a></a></td>
                        <td></td>
                        <td></td>
                        <td></td>
                    </tr>';
            
        }
                        
    if($_POST['dir'] == '/')
        $displayFolders = '';
    
    $displayFiles = '';
    $files = glob($path.$_POST['dir']."*", GLOB_MARK);
    
    /*
     * Iterate through Files & Directories
     */	
    foreach($files as $file){
        
        /*
         * Get Stats on File
         */
        $stat = stat($file);
        
        /*
         * Iterate into HTML Variable
         */
        if(is_dir($file)){
        
            $displayFolders .= 	'<tr>
                                    <td><i class="fa fa-folder-open">&nbsp</i></td>
                                    <td><a href="index.php?dir='.urlencode(str_replace($path, '', $file)).'" class="load_new">'.str_replace($path.$_POST['dir'], '', $file).'</a></td>
                                    <td></td>
                                    <td>'.date('m/d/y H:i:s', $stat['mtime']).'</td>
                                    <td class="center"></td>
                                </tr>';
        
        }else{
        
            $url = (in_array(pathinfo($file, PATHINFO_EXTENSION), $canEdit)) ? '<a href="edit.php&file='.urlencode(str_replace($path, '', $file)).'" class="edit_file">'.str_replace($path.$_POST['dir'], '', $file).'</a>' : str_replace($path.$_POST['dir'], '', $file);

            $displayFiles .= 	'<tr>
                                    <td><i class="fa fa-file-text"></i></td>
                                    <td>'.$url.'</td>
                                    <td>'.$core->framework->files->formatSize($stat['size']).'</td>
                                    <td>'.date('m/d/y H:i:s', $stat['mtime']).'</td>
                                    <td style="text-align:center;"><a href="index.php?do=download&file='.urlencode(str_replace($path, '', $file)).'"><i class="fa fa-download"></i></a></td>
                                </tr>';
        
        }
    
    }
    
}else{

    /*
     * Not in a Directory
     */
    $displayFolders = '';
    $displayFiles = '';
    $files = glob($path."*", GLOB_MARK);
    
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
                                    <td><i class="fa fa-folder-open">&nbsp</i></td>
                                    <td><a href="index.php?dir='.urlencode(str_replace($path, '', $file)).'" class="load_new">'.str_replace($path, '', $file).'</a></td>
                                    <td></td>
                                    <td>'.date('m/d/y H:i:s', $stat['mtime']).'</td>
                                    <td class="center"></td>
                                </tr>';
        
        }else{
        
            $url = (in_array(pathinfo($file, PATHINFO_EXTENSION), $canEdit)) ? '<a href="edit.php?file='.urlencode(str_replace($path, '', $file)).'" class="edit_file">'.str_replace($path, '', $file).'</a>' : $url = str_replace($path, '', $file);
                
            $displayFiles .= 	'<tr>
                                    <td><i class="fa fa-file-text"></i></td>
                                    <td>'.$url.'</td>
                                    <td>'.$core->framework->files->formatSize($stat['size']).'</td>
                                    <td>'.date('m/d/y H:i:s', $stat['mtime']).'</td>
                                    <td style="text-align:center;"><a href="index.php?do=download&file='.urlencode(str_replace($path, '', $file)).'"><i class="fa fa-download"></i></a></td>
                                </tr>';
        
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
                    <th style="width:45%">File Name</th>
                    <th style="width:20%">File Size</th>
                    <th style="width:20%">Last Modified</th>
                    <th style="width:10%;text-align:center;">Options</th>
                </tr>
            </thead>
            <tbody>
                '.$displayFolders.$displayFiles.'
            </tbody>
        </table>';

