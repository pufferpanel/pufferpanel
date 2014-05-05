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

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false){

	Page\components::redirect($core->settings->get('master_url').'index.php?login');
	exit();
    
}

$canEdit = array('txt', 'yml', 'log', 'conf', 'html', 'json', 'properties', 'props', 'cfg', 'lang');

if(isset($_POST['file']))
    $_POST['file'] = str_replace('..', '', urldecode($_POST['file']));

if(isset($_POST['file'])){
        
    if(in_array(pathinfo($_POST['file'], PATHINFO_EXTENSION), $canEdit)){
    
    	/*
         * Create File Path
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
    	                    
		$url = "http://".$core->server->nodeData('sftp_ip').":8003/gameservers/".$core->server->getData('gsd_id')."/file/".$directory.$file;
		
		$data = array("contents" => $_POST['file_contents']);
		
		$curl = curl_init($url);
		curl_setopt($curl, CURLOPT_CUSTOMREQUEST, "PUT");
		curl_setopt($curl, CURLOPT_HTTPHEADER, array("X-Access-Token: ".$core->server->nodeData('gsd_secret')));
		curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
		curl_setopt($curl, CURLOPT_POSTFIELDS, http_build_query($data));
		$response = curl_exec($curl);

        if(empty($response)){
        
        	exit('<div class="alert alert-success">'.$_l->tpl('node_files_ajax_saved').'</div>');
        
        }else{
        
        	exit('<div class="alert alert-danger">'.$_l->tpl('node_files_ajax_no_save').'</div>');
        
        }
    
    }else{
    
        exit('<div class="alert alert-warning">'.$_l->tpl('node_files_ajax_no_edit').'</div>');
    
    }

}else{

    exit('<div class="alert alert-danger">'.$_l->tpl('node_files_ajax_file_404').'</div>');

}