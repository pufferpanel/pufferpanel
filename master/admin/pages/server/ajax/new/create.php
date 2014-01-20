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
require_once('../../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	$core->framework->page->redirect('../../../../index.php');
}

//Cookies :3
setcookie("__TMP_pp_admin_newserver", json_encode($_POST), time() + 30, '/', $core->framework->settings->get('cookie_website'));

/*
 * Set Values
 */
$_POST['server_port'] = $_POST['server_port_'.str_replace('.', '_', $_POST['server_ip'])];

/*
 * Are they all Posted?
 */
if(!isset($_POST['server_name'], $_POST['node'], $_POST['modpack'], $_POST['email'], $_POST['server_ip'], $_POST['server_port'], $_POST['alloc_mem'], $_POST['alloc_disk'], $_POST['sftp_pass'], $_POST['sftp_pass_2'], $_POST['backup_disk'], $_POST['backup_files']))
	$core->framework->page->redirect('../../add.php?disp=missing_args&error=na');

/*
 * Validate Server Name
 */
if(!preg_match('/^[\w-]{4,35}$/', $_POST['server_name']))
	$core->framework->page->redirect('../../add.php?error=server_name&disp=s_fail');
	
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
	$core->framework->page->redirect('../../add.php?error=node&disp=n_fail');

	/*
	 * Validate IP & Port
	 */
	$ips = json_decode($node['ips'], true);
	$ports = json_decode($node['ports'], true);

	if(!array_key_exists($_POST['server_ip'], $ips))
		$core->framework->page->redirect('../../add.php?error=server_ip&disp=ip_fail');
		
	if(!array_key_exists($_POST['server_port'], $ports[$_POST['server_ip']]))
		$core->framework->page->redirect('../../add.php?error=server_port&disp=port_fail');
		
	if($ports[$_POST['server_ip']][$_POST['server_port']] == 0)
		$core->framework->page->redirect('../../add.php?error=server_port&disp=port_full');
	
/*
 * Validate Email
 */	
if(!filter_var($_POST['email'], FILTER_VALIDATE_EMAIL))
	$core->framework->page->redirect('../../add.php?error=email&disp=e_fail');

$selectEmail = $mysql->prepare("SELECT `id` FROM `users` WHERE `email` = ?");
$selectEmail->execute(array($_POST['email']));

	if($selectEmail->rowCount() != 1)
		$core->framework->page->redirect('../../add.php?error=email&disp=a_fail');
	else {
		$oid = $selectEmail->fetch();
		$oid = $oid['id'];
	}

/*
 * Validate Disk & Memory
 */	
if(!is_numeric($_POST['alloc_mem']) || !is_numeric($_POST['alloc_disk']))
	$core->framework->page->redirect('../../add.php?error=alloc_mem|alloc_disk&disp=m_fail');

/*
 * Validate SFTP Password
 */
if($_POST['sftp_pass'] != $_POST['sftp_pass_2'] || strlen($_POST['sftp_pass']) < 8)
	$core->framework->page->redirect('../../add.php?error=sftp_pass|sftp_pass_2&disp=p_fail');				

$iv = base64_encode(mcrypt_create_iv(mcrypt_get_iv_size(MCRYPT_CAST_256, MCRYPT_MODE_CBC), MCRYPT_RAND));
$_POST['sftp_pass'] = openssl_encrypt($_POST['sftp_pass'], 'AES-256-CBC', file_get_contents(HASH), false, base64_decode($iv));

/*
 * Validate Backup Disk & Files
 */
if(!is_numeric($_POST['backup_disk']) || !is_numeric($_POST['backup_files']))
	$core->framework->page->redirect('../../add.php?error=backup_disk|backup_space&disp=b_fail');

/*
 * Validate Modpack
 */
$selectPack = $mysql->prepare("SELECT `min_ram` FROM `modpacks` WHERE `hash` = :hash AND `deleted` = 0");
$selectPack->execute(array(
	':hash' => $_POST['modpack']
));

	if($selectPack->rowCount() != 1)
		$core->framework->page->redirect('../../add.php?error=modpack&disp=no_modpack');
	else
		$pack = $selectPack->fetch();
		
/*
 * Modpack RAM Requirements
 */
if($pack['min_ram'] > $_POST['alloc_mem'])
	$core->framework->page->redirect('../../add.php?error=modpack|alloc_mem&disp=modpack_ram&min_ram='.$pack['min_ram']);

/*
 * Add Server to Database
 */
$ftpUser = (strlen($_POST['server_name']) > 6) ? substr($_POST['server_name'], 0, 6).'_'.$core->framework->auth->keygen(5) : $_POST['server_name'].'_'.$core->framework->auth->keygen((11 - strlen($_POST['server_name'])));

$add = $mysql->prepare("INSERT INTO `servers` VALUES(NULL, NULL, :hash, :e_iv, :node, :sname, :modpack, 1, :oid, :ram, :disk, :date, :sip, :sport, :ftphost, :ftpuser, :ftppass, :bfiles, :bdisk)");
$add->execute(array(
	':hash' => $core->framework->auth->keygen(42),
	':e_iv' => $iv,
	':node' => $_POST['node'],
	':sname' => $_POST['server_name'],
	':modpack' => $_POST['modpack'],
	':oid' => $oid,
	':ram' => $_POST['alloc_mem'],
	':disk' => $_POST['alloc_disk'],
	':date' => time(),
	':sip' => $_POST['server_ip'],
	':sport' => $_POST['server_port'],
	':ftphost' => $node['sftp_ip'],
	':ftpuser' => $ftpUser,
	':ftppass' => $_POST['sftp_pass'],
	':bfiles' => $_POST['backup_files'],
	':bdisk' => $_POST['backup_disk']
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
$con = ssh2_connect($node['sftp_ip'], 22);
ssh2_auth_password($con, $node['username'], openssl_decrypt($node['password'], 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($node['encryption_iv'])));

	/*
	 * Set the Soft Limit
	 */
	$softLimit = ($_POST['alloc_disk'] <= 512) ? 0 : ($_POST['alloc_disk'] - 512);

	/*
	 * Create User
	 */
	$stream = ssh2_exec($con, 'cd /srv/scripts; sudo ./create_user.sh '.$ftpUser.' '.$_POST['sftp_pass_2'].' '.$softLimit.' '.$_POST['alloc_disk'], true);
	$errorStream = ssh2_fetch_stream($stream, SSH2_STREAM_STDERR);
	
	stream_set_blocking($errorStream, true);
	stream_set_blocking($stream, true);
	
	$isError = stream_get_contents($errorStream);
	
	fclose($errorStream);
	fclose($stream);
	
	/*
	 * Install Modpack
	 */
	
		/*
		 * Generate URL
		 */
		$packiv = base64_encode(mcrypt_create_iv(mcrypt_get_iv_size(MCRYPT_CAST_256, MCRYPT_MODE_CBC), MCRYPT_RAND));
		$packEncryptedHash = openssl_encrypt($pack['download_hash'], 'AES-256-CBC', file_get_contents(HASH), false, base64_decode($packiv));
		
		$modpack_request = $core->framework->settings->get('master_url').'modpacks/get.php?pack='.rawurlencode($packEncryptedHash.'.'.$iv);
	
		/*
		 * Execute Commands
		 */
		$stream = ssh2_exec($con, 'cd /srv/scripts; sudo ./install_modpack.sh "'.$ftpUser.'" "'.$modpack_request.'" "'.$pack['hash'].'.zip"', true);
		$errorStream = ssh2_fetch_stream($stream, SSH2_STREAM_STDERR);
		
		stream_set_blocking($errorStream, true);
		stream_set_blocking($stream, true);
		
		$isError = stream_get_contents($errorStream);
		if(!empty($isError))
			echo $isError;
		
		fclose($errorStream);
		fclose($stream);

$core->framework->email->buildEmail('admin_new_server', array(
        'NAME' => $_POST['server_name'],
        'CONNECT' => $_POST['server_ip'].':'.$_POST['server_port'],
        'USER' => $ftpUser,
        'PASS' => $_POST['sftp_pass_2']
))->dispatch($_POST['email'], $core->framework->settings->get('company_name').' - New Server Added');

$core->framework->page->redirect('../../view.php?id='.$lastInsert);

?>