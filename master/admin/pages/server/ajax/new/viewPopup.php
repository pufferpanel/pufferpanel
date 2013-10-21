<?php
session_start();
require_once('../../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	exit('<div class="error-box round">Failed to Authenticate Account.</div>');
}

if(!isset($_GET['node']) || empty($_GET['node']))
	exit('No node defined!');


$selectData = $mysql->prepare("SELECT * FROM `node_data` WHERE `node` = :name");
$selectData->execute(array(
	':name' => $_GET['node']
));
if($selectData->rowCount() == 1)
	$node = $selectData->fetch();
else
	exit('Unknown node.');

$ips = json_decode($node['ips'], true);
$ports = json_decode($node['ports'], true);

$i = 0;
foreach($ips as $ip => $internal){

	if($internal['ports_free'] > 0){
	
		if($i != 0) echo '<br /><br />';
		echo $ip.' has '.$internal['ports_free'].' avaliable port(s).<br />';
		
		foreach($ports[$ip] as $port => $avaliable){
			
			if($avaliable == 1)
				echo ' - - '.$port.' is <span style="color:#52964f;">avaliable</span>.<br />';
			else
				echo ' - - '.$port.' is <span style="color:#cf4425;">unavaliable</span>.<br />';
					
		}
		
	}else{
		
		if($i != 0) echo '<br /><br />';
		echo $ip.' has no avaliable port(s).<br />';
		
	}

	$i++;

}


?>