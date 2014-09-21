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
	Page\components::redirect('../../../../index.php');
}

//Cookies :3
setcookie("__TMP_pp_admin_newnode", json_encode($_POST), time() + 10, '/', $core->settings->get('cookie_website'));

/*
 * Agree Warning
 */
if(!isset($_POST['read_warning']))
	Page\components::redirect('../../add.php?disp=agree_warn');

/*
 * Are they all Posted?
 */
if(!isset($_POST['node_name'], $_POST['node_ip'], $_POST['node_sftp_ip'], $_POST['ip_port']))
	Page\components::redirect('../../add.php?disp=missing_args');

/*
 * Validate Node Name
 */
if(!preg_match('/^[\w.-]{1,15}$/', $_POST['node_name']))
	Page\components::redirect('../../add.php?error=node_name&disp=n_fail');

/*
 * Validate node_ip & node_sftp_ip
 */
if(!filter_var($_POST['node_ip'] , FILTER_VALIDATE_IP, FILTER_FLAG_NO_RES_RANGE) || !filter_var($_POST['node_sftp_ip'] , FILTER_VALIDATE_IP, FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE))
	Page\components::redirect('../../add.php?error=node_ip|node_sftp_ip&disp=ip_fail');

/*
 * Process IPs and Ports
 */
$IPP = array();
$IPA = array();

/*
 * Clean Inputs
 */
$_POST['ip_port'] = str_replace(" ", "", $_POST['ip_port']);

$lines = explode("\r\n", $_POST['ip_port']);
foreach($lines as $id => $values)
	{

		list($ip, $ports) = explode('|', $values);

		$IPA = array_merge($IPA, array($ip => array()));
		$IPP = array_merge($IPP, array($ip => array()));

		$ports = explode(',', $ports);

		for($l=0; $l<count($ports); $l++)
			{

				/*
				 * Validate Port Spacing
				 */
				if(!array_key_exists($l, $IPP[$ip]))
					$IPP[$ip][$ports[$l]] = 1;

			}

		/*
		 * Make sure Ports are in the array
		 */
		if(count($IPP[$ip]) > 0)
			$IPA[$ip] = array_merge($IPA[$ip], array("ports_free" => count($IPP[$ip])));
		else
			Page\components::redirect('../../add.php?error=ip_port&disp=ip_port_space');

	}

/*
 * Add Node to Database
 */
$create = $mysql->prepare("INSERT INTO `nodes` VALUES(NULL, :name, :ip, :sftp_ip, :gsd_secret, :ips, :ports)");
$create->execute(array(
	':name' => $_POST['node_name'],
	':ip' => $_POST['node_ip'],
	':sftp_ip' => $_POST['node_sftp_ip'],
	':gsd_secret' => $core->auth->gen_UUID(),
	':ips' => json_encode($IPA),
	':ports' => json_encode($IPP)
));

Page\components::redirect('../../view.php?id='.$mysql->lastInsertId());

?>
