<?php
session_start();
require_once('../../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	exit('<div class="error-box round">Failed to Authenticate Account.</div>');
}

if(!isset($_POST['node'], $_POST['port'], $_POST['ip']))
	exit('POST Only');

/*
 * Verify port is Real & Not in Use
 */
$select = $mysql->prepare("SELECT `ips`, `ports` FROM `nodes` WHERE `id` = :nid");
$select->execute(array(
	':nid' => $_POST['node']
));

if($select->rowCount() != 1)
	exit('Invalid Node');


$data = $select->fetch();

$ips = json_decode($data['ips'], true);
$ports = json_decode($data['ports'], true);

if(array_key_exists($_POST['ip'], $ports) && array_key_exists($_POST['port'], $ports[$_POST['ip']]) && $ports[$_POST['ip']][$_POST['port']] == 1){
	
	unset($ports[$_POST['ip']][$_POST['port']]);
	$ips[$_POST['ip']]['ports_free'] = ($ips[$_POST['ip']]['ports_free'] - 1);
	
}else{

	exit('No Port/IP or Port in Use');

}
	
	
$update = $mysql->prepare("UPDATE `nodes` SET `ips` = :ips, `ports` = :ports WHERE `id` = :nid");
$update->execute(array(
	':nid' => $_POST['node'],
	':ips' => json_encode($ips),
	':ports' => json_encode($ports)
));

echo 'Done';

?>