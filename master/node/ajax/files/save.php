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

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), $core->framework->auth->getCookie('pp_server_hash')) === false){

	$core->framework->page->redirect($core->framework->settings->get('master_url').'index.php');
	exit();
    
}

$canEdit = array('txt', 'yml', 'log', 'conf', 'html', 'json', 'properties', 'props', 'cfg', 'lang');

if(isset($_POST['file']))
    $_POST['file'] = str_replace('..', '', urldecode($_POST['file']));

$path = $core->framework->server->nodeData('server_dir').$core->framework->server->getData('ftp_user').'/server/';

if(isset($_POST['file'])){
        
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
                if($directory == './' || $directory == '.')
                    $directory = '';
        
                if(substr($directory, 0, 1) == '/')
                    $directory = substr($directory, 1);
                
                $fp = fopen($saveDir.'save.'.$file , 'w');
                fwrite($fp, $_POST['file_contents']);
                fclose($fp);
                
                    /*
                     * Upload Via SFTP
                     */
                    $connection = $core->framework->ssh->generateSSH2Connection(array(
                    	'ip' => $core->framework->server->nodeData('sftp_ip'),
                    	'user' => $core->framework->server->getData('ftp_user'),
                    	'pass' => $core->framework->server->getData('ftp_pass'),
                    	'iv' => $core->framework->server->getData('encryption_iv')
                    ), null, true);
                    
                        $FTPLocalFile = $saveDir.'save.'.$file;
                        $sftp = ssh2_sftp($connection);
                                                     
                            $stream = fopen("ssh2.sftp://$sftp/server/".$directory.$file, 'w+');
                            
                                if(!$stream){
                                
                                    exit('<div class="alert alert-danger">Unable to connect and upload file. This is usually due to a permissions error.</div>');
                                
                                }else{
                                    
                                    if(fwrite($stream, file_get_contents($FTPLocalFile))){
                                    
                                        fclose($stream);
                                        unlink($FTPLocalFile);
                                        exit('<div class="alert alert-success">File was sucessfully saved.</div>');
                                    
                                    }else{
                                    
                                        fclose($stream);
                                        exit('<div class="alert alert-danger">Unknown error. Unable to save file.</div>');
                                    
                                    }
                                
                                }
    
    }else{
    
        exit('<div class="alert alert-warning">This type of file cannot be edited via our online file manager. Please use a FTP client.</div>');
    
    }

}else{

    exit('<div class="alert alert-danger">The file specified could not be found on the server.</div>');

}