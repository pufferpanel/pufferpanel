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

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), $core->framework->auth->getCookie('pp_server_hash')) === true){

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
		if($pack['min_ram'] > $core->framework->server->getData('max_ram'))
			exit('<div class="alert alert-error">Your server does not have enough RAM allocated to run this Modpack. This Modpack requires '.$pack['min_ram'].'MB of RAM, you have '.$core->framework->server->getData('max_ram').'MB allocated.</div>');
		
		/*
		 * Generate URL
		 */
		$iv = base64_encode(mcrypt_create_iv(mcrypt_get_iv_size(MCRYPT_CAST_256, MCRYPT_MODE_CBC), MCRYPT_RAND));
		$encryptedHash = openssl_encrypt($pack['download_hash'], 'AES-256-CBC', file_get_contents(HASH), false, base64_decode($iv));
		
		$modpack_request = $core->framework->settings->get('master_url').'modpacks/get.php?pack='.rawurlencode($encryptedHash.'.'.$iv);
		
		/*
		 * Get Node Info
		 */
		$nodeSQLConnect = $mysql->prepare("SELECT * FROM `nodes` WHERE `id` = ? LIMIT 1");
		$nodeSQLConnect->execute(array($core->framework->server->getData('node')));
		$node = $nodeSQLConnect->fetch();
		
		/*
		 * Connect and Run Function
		 */
//		$con = ssh2_connect($node['sftp_ip'], 22);
//		ssh2_auth_password($con, $node['username'], openssl_decrypt($node['password'], 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($node['encryption_iv'])));
		
		exit('cd /srv/scripts; sudo ./install_modpack.sh '.$core->framework->server->getData('ftp_user').' "'.$modpack_request.'" "'.$pack['hash'].'.zip"');
		
//		$stream = ssh2_exec($con, 'cd /srv/scripts; sudo ./install_modpack.sh '.$core->framework->server->getData('ftp_user').' "'.$modpack_request.'" "'.$pack['hash'].'.zip"', true);
//		$errorStream = ssh2_fetch_stream($stream, SSH2_STREAM_STDERR);
//		
//		stream_set_blocking($errorStream, true);
//		stream_set_blocking($stream, true);
//		
//		$isError = stream_get_contents($errorStream);
//		if(!empty($isError))
//			echo $isError;
//		
//		fclose($errorStream);
//		fclose($stream);
		
//        $core->framework->log->getUrl()->addLog(0, 1, array('user.modpack_install', 'A new modpack was installed for this server. The modpack installed was '.$pack['name'].' ('.$pack['version'].').'));
			
	}

}
?>