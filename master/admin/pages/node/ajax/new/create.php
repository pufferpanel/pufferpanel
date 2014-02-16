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
error_reporting('E_ALL');
require_once('../../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), null, true) !== true){
	$core->framework->page->redirect('../../../../index.php');
}
	
//Cookies :3
setcookie("__TMP_pp_admin_newnode", json_encode($_POST), time() + 10, '/', $core->framework->settings->get('cookie_website'));

/*
 * Agree Warning
 */
if(!isset($_POST['read_warning']))
	$core->framework->page->redirect('../../add.php?disp=agree_warn');

/*
 * Are they all Posted?
 */
if(!isset($_POST['node_name'], $_POST['node_url'], $_POST['node_ip'], $_POST['node_sftp_ip'], $_POST['s_dir'], $_POST['s_dir_backup'], $_POST['ssh_user'], $_POST['ssh_pass'], $_POST['ip_port']))
	$core->framework->page->redirect('../../add.php?disp=missing_args');

/*
 * Validate Node Name
 */
if(!preg_match('/^[\w.-]{1,15}$/', $_POST['node_name']))
	$core->framework->page->redirect('../../add.php?error=node_name&disp=n_fail');
		
/*
 * Validate node_ip & node_sftp_ip
 */
if(!filter_var($_POST['node_ip'] , FILTER_VALIDATE_IP, FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE) || !filter_var($_POST['node_sftp_ip'] , FILTER_VALIDATE_IP, FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE))
	$core->framework->page->redirect('../../add.php?error=node_ip|node_sftp_ip&disp=ip_fail');

if(!preg_match('/^[a-zA-Z0-9_\.\/-]+[^\/]\/$/', $_POST['s_dir']) || !preg_match('/^[a-zA-Z0-9_\.\/-]+[^\/]\/$/', $_POST['s_dir_backup']))
	$core->framework->page->redirect('../../add.php?error=s_dir|s_dir_backup&disp=dir_fail');
	
if($_POST['s_dir'] == $_POST['s_dir_backup'])
	$core->framework->page->redirect('../../add.php?error=s_dir|s_dir_backup&disp=dir_match_fail');
	
if(strlen($_POST['ssh_user']) < 1 || $_POST['ssh_user'] == 'root')
	$core->framework->page->redirect('../../add.php?error=ssh_user&disp=user_fail');

if(strlen($_POST['ssh_pass']) < 12)
	$core->framework->page->redirect('../../add.php?error=ssh_pass&disp=pass_fail');

/*
 * Process IPs and Ports
 */
$IPP = array();
$IPA = array();
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
				if(!array_key_exists($l - 1, $IPP[$ip]) && !array_key_exists($l + 1, $IPP[$ip])
					$IPP[$ip][$ports[$l]] = 1;
			
			}
			
		/*
		 * Make sure Ports are in the array
		 */
		if(count($IPP[$ip]) > 0)
			$IPA[$ip] = array_merge($IPA[$ip], array("ports_free" => count($IPP[$ip])));
		else
			$core->framework->page->redirect('../../add.php?error=ip_port&disp=ip_port_space');
			
	}

$IPA = json_encode($IPA);
$IPP = json_encode($IPP);

$iv = base64_encode(mcrypt_create_iv(mcrypt_get_iv_size(MCRYPT_CAST_256, MCRYPT_MODE_CBC), MCRYPT_RAND));
$_POST['ssh_pass'] = openssl_encrypt($_POST['ssh_pass'], 'AES-256-CBC', file_get_contents(HASH), false, base64_decode($iv));

$create = $mysql->prepare("INSERT INTO `nodes` VALUES(NULL, :name, :ip, :sftp_ip, :sdir, :bdir, :suser, :iv, :spass, :ips, :ports)");
$create->execute(array(
	':name' => $_POST['node_name'],
	':ip' => $_POST['node_ip'],
	':sftp_ip' => $_POST['node_sftp_ip'],
	':sdir' => $_POST['s_dir'],
	'bdir' => $_POST['s_dir_backup'],
	':suser' => $_POST['ssh_user'],
	':iv' => $iv,
	':spass' => $_POST['ssh_pass'],
	':ips' => $IPA,
	':ports' => $IPP
));

$core->framework->page->redirect('../../view.php?id='.$mysql->lastInsertId());

?>