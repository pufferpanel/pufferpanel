<?php
session_start();
require_once('../../../core/framework/framework.core.php');
$filesIncluded = true;

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), $core->framework->auth->getCookie('pp_server_hash')) === true){

	echo $core->framework->files->last_lines($core->framework->server->getData('path').'server.log', 500);
	
}else{

	exit('Invalid Authentication.');

}
?>