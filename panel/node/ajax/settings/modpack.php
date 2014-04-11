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

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === true){

	if(isset($_POST['new_pack'])){
	
		if(!isset($_POST['new_pack']) || empty($_POST['new_pack']))
			exit('<div class="alert alert-error">No new Modpack was submitted with the request.</div>');
			
		/*
		 * Does the Pack Exist?
		 */
		$packs = $mysql->prepare("SELECT * FROM `modpacks` WHERE `hash` = :hash AND `deleted` = 0");
		$packs->execute(array(
			':hash' => $_POST['new_pack']
		));
		
			if($packs->rowCount() != 1)
				exit('<div class="alert alert-error">That Modpack does not exist in our system.</div>');
			else
				$pack = $packs->fetch();
		
		/*
		 * Minimum Requirements Met?
		 */
		if($pack['min_ram'] > $core->server->getData('max_ram'))
			exit('<div class="alert alert-error">Your server does not have enough RAM allocated to run this Modpack. This Modpack requires '.$pack['min_ram'].'MB of RAM, you have '.$core->server->getData('max_ram').'MB allocated.</div>');
		
		/*
		 * Generate URL
		 */
		$iv = $core->auth->generate_iv();
		$encryptedHash = $core->auth->encrypt($pack['download_hash'], $iv);
				
		$modpack_request = $core->settings->get('master_url').'modpacks/get.php?pack='.rawurlencode($encryptedHash.'.'.$iv);

		/*
		 * Connect and Run Function
		 */
		$callbackData = $core->ssh->generateSSH2Connection(array(
			'ip' => $core->server->nodeData('sftp_ip'),
			'user' => $core->server->nodeData('username')
		), array(
			'pub' => $core->server->nodeData('ssh_pub'),
			'priv' => $core->server->nodeData('ssh_priv'),
			'secret' => $core->server->nodeData('ssh_secret'),
			'secret_iv' => $core->server->nodeData('ssh_secret_iv')
		))->executeSSH2Command('cd /srv/scripts; sudo ./install_modpack.sh '.$core->server->getData('ftp_user').' "'.$modpack_request.'" "'.$pack['hash'].'.zip"', true);
		
		if(!empty($callbackData))
			echo $callbackData;
		
		
        $core->log->getUrl()->addLog(0, 1, array('user.modpack_install', 'A new modpack was installed for this server. The modpack installed was '.$pack['name'].' ('.$pack['version'].').'));
        
        /*
         * Update SQL
         */
        $mysql->prepare("UPDATE `servers` SET `modpack` = :pack WHERE `id` = :sid")->execute(array(
        	':pack' => $pack['hash'],
        	':sid' => $core->server->getData('id')
        ));
        
        echo '<div class="alert alert-success">Modpack successfully installed.</div>';
			
	}

}
?>