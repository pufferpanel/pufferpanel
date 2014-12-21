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
namespace PufferPanel\Core;
use \ORM;

require_once('../../../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true) {
	Components\Page::redirect('../../../index.php');
}

setcookie("__TMP_pp_admin_newserver", json_encode($_POST), time() + 30, '/');

/*
* Set Values
*/
@$_POST['server_port'] = $_POST['server_port_'.str_replace('.', '_', $_POST['server_ip'])];

/*
 * Are they all Posted?
 */
if(!isset($_POST['server_name'], $_POST['node'], $_POST['email'], $_POST['server_ip'], $_POST['server_port'], $_POST['alloc_mem'], $_POST['alloc_disk'], $_POST['ftp_pass'], $_POST['ftp_pass_2'], $_POST['cpu_limit'])) {
	Components\Page::redirect('../../add.php?disp=missing_args&error=na');
}

/*
* Determine if Node (IP & Port) is Avaliable
*/
$node = ORM::forTable('nodes')->findOne($_POST['node']);

if(!$node) {
	Components\Page::redirect('../../add.php?error=node&disp=n_fail');
}

/*
* GSD Must Be Online!
*/
if(!@fsockopen($node->ip, $node->gsd_listen, $num, $error, 3)) {
	Components\Page::redirect('../../add.php?disp=gsd_offline&error=na');
}

/*
 * Validate Server Name
 */
if(!preg_match('/^[\w-]{4,35}$/', $_POST['server_name'])) {
	Components\Page::redirect('../../add.php?error=server_name&disp=s_fail');
}

	/*
	 * Validate IP & Port
	 */
	$ips = json_decode($node->ips, true);
	$ports = json_decode($node->ports, true);

	if(!array_key_exists($_POST['server_ip'], $ips)) {
		Components\Page::redirect('../../add.php?error=server_ip&disp=ip_fail');
	}

	if(!array_key_exists($_POST['server_port'], $ports[$_POST['server_ip']])) {
		Components\Page::redirect('../../add.php?error=server_port&disp=port_fail');
	}

	if($ports[$_POST['server_ip']][$_POST['server_port']] == 0) {
		Components\Page::redirect('../../add.php?error=server_port&disp=port_full');
	}

/*
 * Validate Email
 */
if(!filter_var($_POST['email'], FILTER_VALIDATE_EMAIL)) {
	Components\Page::redirect('../../add.php?error=email&disp=e_fail');
}

$user = ORM::forTable('users')->select('id')->where('email', $_POST['email'])->findOne();

	if(!$user) {
		Components\Page::redirect('../../add.php?error=email&disp=a_fail');
	}

/*
 * Validate Disk & Memory
 */
if(!is_numeric($_POST['alloc_mem']) || !is_numeric($_POST['alloc_disk'])) {
	Components\Page::redirect('../../add.php?error=alloc_mem|alloc_disk&disp=m_fail');
}

/*
 * Validate CPU Limit
 */
if(!is_numeric($_POST['cpu_limit'])) {
	Components\Page::redirect('../../add.php?error=cpu_limit&disp=cpu_limit');
}


/*
 * Validate ftp Password
 */
if($_POST['ftp_pass'] != $_POST['ftp_pass_2'] || strlen($_POST['ftp_pass']) < 8) {
	Components\Page::redirect('../../add.php?error=ftp_pass|ftp_pass_2&disp=p_fail');
}

$iv = $core->auth->generate_iv();
$_POST['ftp_pass'] = $core->auth->encrypt($_POST['ftp_pass'], $iv);

/*
 * Add Server to Database
 */
$ftpUser = Components\Functions::generateFTPUsername($_POST['server_name']);
$modpack = (isset($pack) && is_array($pack)) ? $pack['server_jar'] : 'server.jar';

/*
 * Create Unique Values
 */
$serverHash = $core->auth->generateUniqueUUID('servers', 'hash');
$gsdSecret = $core->auth->generateUniqueUUID('servers', 'gsd_secret');

/*
 * Build Call
 */
$data = array(
	"name" => $serverHash,
	"user" => $ftpUser,
	"overide_command_line" => "",
	"path" => $node->gsd_server_dir.$ftpUser,
	"build" => array(
		"install_dir" => '/mnt/MC/CraftBukkit/',
		"disk" => array(
			"hard" => ($_POST['alloc_disk'] < 32) ? 32 : (int) $_POST['alloc_disk'],
			"soft" => ($_POST['alloc_disk'] > 2048) ? (int) $_POST['alloc_disk'] - 1024 : 32
		),
		"cpu" => (int) $_POST['cpu_limit']
	),
	"variables" => array(
		"-jar" => $modpack,
		"-Xmx" => $_POST['alloc_mem']."M"
	),
	"keys" => array(
		$gsdSecret => array("s:ftp", "s:get", "s:power", "s:files", "s:files:get", "s:files:put", "s:files:zip", "s:query", "s:console", "s:console:send")
	),
	"gameport" => (int) $_POST['server_port'],
	"gamehost" => $_POST['server_ip'],
	"plugin" => "minecraft",
	"autoon" => false
);

$data = json_encode($data);

$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, 'http://'.$node->ip.':'.$node->gsd_listen.'/gameservers');
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
$server = ORM::forTable('servers')->create();
$server->set(array(
	'gsd_id' => $content['id'],
	'hash' => $serverHash,
	'gsd_secret' => $gsdSecret,
	'encryption_iv' => $iv,
	'node' => $_POST['node'],
	'name' => $_POST['server_name'],
	'modpack' => '0000-0000-0000-0',
	'server_jar' => $modpack,
	'owner_id' => $user->id,
	'max_ram' => $_POST['alloc_mem'],
	'disk_space' => $_POST['alloc_disk'],
	'cpu_limit' => $_POST['cpu_limit'],
	'date_added' => time(),
	'server_ip' => $_POST['server_ip'],
	'server_port' => $_POST['server_port'],
	'ftp_user' => $ftpUser,
	'ftp_pass' => $_POST['ftp_pass']
));
$server->save();

/*
* Update IP Count
*/
$ips[$_POST['server_ip']]['ports_free']--;
$ports[$_POST['server_ip']][$_POST['server_port']]--;

$node->ips = json_encode($ips);
$node->ports = json_encode($ports);
$node->save();

/*
 * Send User Email
 */
$core->email->buildEmail('admin_new_server', array(
		'NAME' => $_POST['server_name'],
		'FTP' => $node->fqdn.':21',
		'MINECRAFT' => $node->fqdn.':'.$_POST['server_port'],
		'USER' => $ftpUser.'-'.$content['id'],
		'PASS' => $_POST['ftp_pass_2']
))->dispatch($_POST['email'], $core->settings->get('company_name').' - New Server Added');

Components\Page::redirect('../../view.php?id='.$server->id());

?>
