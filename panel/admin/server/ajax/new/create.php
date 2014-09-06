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
require_once('../../../../../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Page\components::redirect('../../../index.php');
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
if(!isset($_POST['server_name'], $_POST['node'], $_POST['email'], $_POST['server_ip'], $_POST['server_port'], $_POST['alloc_mem'], $_POST['alloc_disk'], $_POST['sftp_pass'], $_POST['sftp_pass_2'], $_POST['cpu_limit']))
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

$iv = $core->auth->generate_iv();
$_POST['sftp_pass'] = $core->auth->encrypt($_POST['sftp_pass'], $iv);

/*
 * Add Server to Database
 */
$ftpUser = Functions\general::generateFTPUsername($_POST['server_name']);

$serverHash = $core->auth->gen_UUID();
$modpack = (isset($pack) && is_array($pack)) ? $pack['server_jar'] : 'server.jar';

$add = $mysql->prepare("INSERT INTO `servers` VALUES(NULL, NULL, NULL, :hash, :gsd_secret, :e_iv, :node, :sname, :modpack, :sjar, 1, :oid, :ram, :disk, :cpu, :date, :sip, :sport, :ftpuser, :ftppass)");
$add->execute(array(
	':hash' => $serverHash,
	':gsd_secret' => $node['gsd_secret'],
	':e_iv' => $iv,
	':node' => $_POST['node'],
	':sname' => $_POST['server_name'],
	':modpack' => '0000-0000-0000-0',
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

$mysql->prepare("UPDATE `nodes` SET `ips` = :ips, `ports` = :ports WHERE `id` = :id")->execute(array(
	':ips' => json_encode($ips),
	':ports' => json_encode($ports),
	':id' => $_POST['node']
));

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
        "path" => "/home/".$ftpUser,
        "variables" => array(
        	"-jar" => $modpack,
            "-Xmx" => $_POST['alloc_mem']."M"
        ),
        "gameport" => $_POST['server_port'],
        "gamehost" => "",
        "plugin" => "minecraft",
        "autoon" => false
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
	 * Send User Email
	 */
	$core->email->buildEmail('admin_new_server', array(
	        'NAME' => $_POST['server_name'],
	        'CONNECT' => $node['sftp_ip'].':21',
	        'USER' => $ftpUser.'-'.$content['id'],
	        'PASS' => $_POST['sftp_pass_2']
	))->dispatch($_POST['email'], $core->settings->get('company_name').' - New Server Added');

	Page\components::redirect('../../view.php?id='.$lastInsert);

?>
