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
require_once('../../../../../../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Page\components::redirect('../../../../index.php');
}

//Cookies :3
setcookie("__TMP_pp_admin_newserver", json_encode($_POST), time() + 30, '/', $core->settings->get('cookie_website'));

/*
 * Set Values
 */
$_POST['server_port'] = $_POST['server_port_'.str_replace('.', '_', $_POST['server_ip'])];

/*
 * Are they all Posted?
 */
if(!isset($_POST['server_name'], $_POST['node'], $_POST['modpack'], $_POST['email'], $_POST['server_ip'], $_POST['server_port'], $_POST['alloc_mem'], $_POST['alloc_disk'], $_POST['sftp_pass'], $_POST['sftp_pass_2'], $_POST['cpu_limit']))
	Page\components::redirect('../../add.php?disp=missing_args&error=na');

/*
 * Validate Server Name
 */
if(!preg_match('/^[\w-]{4,35}$/', $_POST['server_name']))
	Page\components::redirect('../../add.php?error=server_name&disp=s_fail');
	
/*
 * Determine if Node (IP & Port) is Avaliable
 */
$select = $mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :id");
$select->execute(array(
	':id' => $_POST['node']
));

if($select->rowCount() == 1)
	$node = $select->fetch();
else
	Page\components::redirect('../../add.php?error=node&disp=n_fail');

	/*
	 * Validate IP & Port
	 */
	$ips = json_decode($node['ips'], true);
	$ports = json_decode($node['ports'], true);

	if(!array_key_exists($_POST['server_ip'], $ips))
		Page\components::redirect('../../add.php?error=server_ip&disp=ip_fail');
		
	if(!array_key_exists($_POST['server_port'], $ports[$_POST['server_ip']]))
		Page\components::redirect('../../add.php?error=server_port&disp=port_fail');
		
	if($ports[$_POST['server_ip']][$_POST['server_port']] == 0)
		Page\components::redirect('../../add.php?error=server_port&disp=port_full');
	
/*
 * Validate Email
 */	
if(!filter_var($_POST['email'], FILTER_VALIDATE_EMAIL))
	Page\components::redirect('../../add.php?error=email&disp=e_fail');

$selectEmail = $mysql->prepare("SELECT `id` FROM `users` WHERE `email` = ?");
$selectEmail->execute(array($_POST['email']));

	if($selectEmail->rowCount() != 1)
		Page\components::redirect('../../add.php?error=email&disp=a_fail');
	else {
		$oid = $selectEmail->fetch();
		$oid = $oid['id'];
	}

/*
 * Validate Disk & Memory
 */	
if(!is_numeric($_POST['alloc_mem']) || !is_numeric($_POST['alloc_disk']))
	Page\components::redirect('../../add.php?error=alloc_mem|alloc_disk&disp=m_fail');

/*
 * Validate CPU Limit
 */	
if(!is_numeric($_POST['cpu_limit']))
	Page\components::redirect('../../add.php?error=cpu_limit&disp=cpu_limit');


/*
 * Validate SFTP Password
 */
if($_POST['sftp_pass'] != $_POST['sftp_pass_2'] || strlen($_POST['sftp_pass']) < 8)
	Page\components::redirect('../../add.php?error=sftp_pass|sftp_pass_2&disp=p_fail');				

$iv = base64_encode(mcrypt_create_iv(mcrypt_get_iv_size(MCRYPT_CAST_256, MCRYPT_MODE_CBC), MCRYPT_RAND));
$_POST['sftp_pass'] = openssl_encrypt($_POST['sftp_pass'], 'AES-256-CBC', file_get_contents(HASH), false, base64_decode($iv));

if($_POST['modpack'] != "none"){

	/*
	 * Validate Modpack
	 */
	$selectPack = $mysql->prepare("SELECT `min_ram`, `server_jar` FROM `modpacks` WHERE `hash` = :hash AND `deleted` = 0");
	$selectPack->execute(array(
		':hash' => $_POST['modpack']
	));
	
		if($selectPack->rowCount() != 1)
			Page\components::redirect('../../add.php?error=modpack&disp=no_modpack');
		else
			$pack = $selectPack->fetch();
			
	/*
	 * Modpack RAM Requirements
	 */
	if($pack['min_ram'] > $_POST['alloc_mem'])
		Page\components::redirect('../../add.php?error=modpack|alloc_mem&disp=modpack_ram&min_ram='.$pack['min_ram']);

}

/*
 * Add Server to Database
 */
$ftpUser = (strlen($_POST['server_name']) > 6) ? substr($_POST['server_name'], 0, 6).'_'.$core->auth->keygen(5) : $_POST['server_name'].'_'.$core->auth->keygen((11 - strlen($_POST['server_name'])));

$serverHash = $core->auth->keygen(42);
$modpack = (isset($pack) && is_array($pack)) ? $pack['server_jar'] : 'server.jar';

$add = $mysql->prepare("INSERT INTO `servers` VALUES(NULL, NULL, NULL, :hash, :gsd_secret, :e_iv, :node, :sname, :modpack, :sjar, 1, :oid, :ram, :disk, :cpu, :date, :sip, :sport, :ftpuser, :ftppass)");
$add->execute(array(
	':hash' => $serverHash,
	':gsd_secret' => $core->auth->keygen(16).$core->auth->keygen(16),
	':e_iv' => $iv,
	':node' => $_POST['node'],
	':sname' => $_POST['server_name'],
	':modpack' => $_POST['modpack'],
	':sjar' => $modpack,
	':oid' => $oid,
	':ram' => $_POST['alloc_mem'],
	':disk' => $_POST['alloc_disk'],
	':cpu' => $_POST['cpu_limit'],
	':date' => time(),
	':sip' => $_POST['server_ip'],
	':sport' => $_POST['server_port'],
	':ftpuser' => $ftpUser,
	':ftppass' => $_POST['sftp_pass']
));

$lastInsert = $mysql->lastInsertId();

/*
 * Update IP Count
 */
$ips[$_POST['server_ip']]['ports_free']--;
$ports[$_POST['server_ip']][$_POST['server_port']]--;

$mysql->prepare("UPDATE `nodes` SET `ips` = :ips")->execute(array(':ips' => json_encode($ips)));
$mysql->prepare("UPDATE `nodes` SET `ports` = :ports")->execute(array(':ports' => json_encode($ports)));

/*
 * Do Server Making Stuff Here 
 */

	/*
	 * Build Call
	 */
	$data = array(
		"name" => $serverHash,
        "user" => $ftpUser,
        "overide_command_line" => "",
        "path" => $node['server_dir'].$ftpUser."/server/",
        "variables" => array(
        	"-jar" => $pack['server_jar'],
            "-Xmx" => $_POST['alloc_mem']."M"
        ),
        "gameport" => $_POST['server_port'],
        "gamehost" => "",
        "plugin" => "minecraft"
	);

	$data = json_encode($data);

	$ch = curl_init();
	curl_setopt($ch, CURLOPT_URL, 'http://'.$node['sftp_ip'].':8003/gameservers');
	curl_setopt($ch, CURLOPT_HTTPHEADER, array(                                                                          
	    'X-Access-Token: '.$node['gsd_secret']
	));
	curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 5); 
	curl_setopt($ch, CURLOPT_TIMEOUT, 5);
	curl_setopt($ch, CURLOPT_POST, true);
	curl_setopt($ch, CURLOPT_POSTFIELDS, "settings=".$data);
	curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
	$content = json_decode(curl_exec($ch), true);
	curl_close($ch);

	/*
	 * Update MySQL
	 */
	$update = $mysql->prepare("UPDATE `servers` SET `gsd_id` = :gsdid WHERE `hash` = :hash");
	$update->execute(array(
		':gsdid' => $content['id'],
		':hash' => $serverHash
	));
	
	/*
	 * Create the User
	 */
	$softLimit = ($_POST['alloc_disk'] <= 512) ? 0 : ($_POST['alloc_disk'] - 512);
	
	$core->ssh->generateSSH2Connection($node['id'], true)->executeSSH2Command('cd /srv/scripts; sudo ./create_user.sh '.$ftpUser.' '.$_POST['sftp_pass_2'].' '.$softLimit.' '.$_POST['alloc_disk'], false);
	
	
	//if($_POST['modpack'] != "none") {
	//
	//	/*
	//	 * Install Modpack
	//	 */
	//	
	//		/*
	//		 * Generate URL
	//		 */
	//		$packiv = $core->auth->generate_iv();
	//		$packEncryptedHash = $core->auth->encrypt($pack['download_hash'], $packiv);
	//		
	//		$modpack_request = $core->settings->get('master_url').'modpacks/get.php?pack='.rawurlencode($packEncryptedHash.'.'.$iv);
	//	
	//		/*
	//		 * Execute Commands
	//		 */
	//		$core->ssh->generateSSH2Connection($node['id'], true)->executeSSH2Command('cd /srv/scripts; sudo ./install_modpack.sh "'.$ftpUser.'" "'.$modpack_request.'" "'.$pack['hash'].'.zip"', false);
	//
	//}

	/*
	 * Send User Email
	 */
	$core->email->buildEmail('admin_new_server', array(
	        'NAME' => $_POST['server_name'],
	        'CONNECT' => $_POST['server_ip'].':'.$_POST['server_port'],
	        'USER' => $ftpUser,
	        'PASS' => $_POST['sftp_pass_2']
	))->dispatch($_POST['email'], $core->settings->get('company_name').' - New Server Added');
	
	Page\components::redirect('../../view.php?id='.$lastInsert);

?>
