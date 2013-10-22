<?php
session_start();
require_once('../../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	exit('<div class="error-box round">Failed to Authenticate Account.</div>');
}

//Cookies :3
setcookie("__TMP_pp_admin_newserver", json_encode($_POST), time() + 60, '/', $core->framework->settings->get('cookie_website'));

/*
 * Are they all Posted?
 */
if(!isset($_POST['server_name'], $_POST['node'], $_POST['email'], $_POST['server_ip'], $_POST['server_port'], $_POST['alloc_mem'], $_POST['alloc_disk'], $_POST['sftp_pass'], $_POST['sftp_pass_2'], $_POST['backup_disk'], $_POST['backup_files']))
	$core->framework->page->redirect('../../add.php?disp=missing_args');

/*
 * Validate Server Name
 */
if(!preg_match('/^[\w-]{4,35}$/', $_POST['server_name']))
	$core->framework->page->redirect('../../add.php?error=server_name&disp=s_fail');
	
/*
 * Determine if Node (IP & Port) is Avaliable
 */
$select = $mysql->prepare("SELECT * FROM `nodes` WHERE `node` = :name");
$select->execute(array(
	':name' => $_POST['node']
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
 * Add Server to Database
 */
$add = $mysql->prepare("INSERT INTO `servers` VALUES(NULL, NULL, :hash, :e_iv, :node, :sname, 1, :oid, :ram, :disk, :path, :date, :sip, :sport, :ftphost, :ftpuser, :ftppass, :bfiles, :bdisk)");
$add->execute(array(
	':hash' => $core->framework->auth->keygen(42),
	':e_iv' => $iv,
	':node' => $_POST['node'],
	':sname' => $_POST['server_name'],
	':oid' => $oid,
	':ram' => $_POST['alloc_mem'],
	':disk' => $_POST['alloc_disk'],
	':path' => $node['server_dir'].$_POST['server_name'].'/',
	':date' => time(),
	':sip' => $_POST['server_ip'],
	':sport' => $_POST['server_port'],
	':ftphost' => $node['sftp_ip'],
	':ftpuser' => (strlen($_POST['server_name']) > 6) ? substr($_POST['server_name'], 0, 6).'_'.$core->framework->auth->keygen(5) : $_POST['server_name'].'_'.$core->framework->auth->keygen((11 - strlen($_POST['server_name']))),
	':ftppass' => $_POST['sftp_pass'],
	':bfiles' => $_POST['backup_files'],
	':bdisk' => $_POST['backup_disk']
));

/*
 * Do Server Making Stuff Here
 */

?>