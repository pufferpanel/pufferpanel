<?php
session_start();
require_once('../../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	exit('<div class="error-box round">Failed to Authenticate Account.</div>');
}

/*
Array
(
    [server_name] => something
    [node] => demo1
    [email] => dane@daneeveritt.com
    [server_ip] => 63.143.53.10
    [server_port] => 25515
    [alloc_mem] => 1024
    [alloc_disk] => 10240
    [sftp_host] => 63.143.53.10
    [sftp_pass] => lWha9F6M6lj4
    [sftp_pass_2] => lWha9F6M6lj4
    [backup_disk] => 1024
    [backup_files] => 10
)
*/

foreach($_POST as $key => $val)
	{
	
		if(empty($val))
			exit();
			
		if($key == 'server_name' && !preg_match('/^[\w-]{4,35}$/', $val))
			exit();
			
		if($key == 'node')
			{
		
				/*
				 * Determine if Node (IP & Port) is Avaliable
				 */
		
			}
			
		if($key == 'email')
			{
			
				if(!filter_var($val, FILTER_VALIDATE_EMAIL))
					exit();
			
				$select = $mysql->prepare("SELECT `id` FROM `users` WHERE `email` = ?");
				$select->execute(array($val));
				
					if($select->rowCount() != 1)
						exit();
					else {
						$oid = $select->fetch();
						$oid = $oid['id'];
					}
					
			}
			
		if($key == 'alloc_mem' || $key == 'alloc_disk' && !is_numeric($val))
			exit();
			
		if($key == 'sftp_host' && !preg_match())
			exit();
			
		if($key == 'sftp_pass')
			{
			
				if($val != $_POST['sftp_pass_2'] || strlen($val) < 8)
					exit();				
				
				$iv = base64_encode(mcrypt_create_iv(mcrypt_get_iv_size(MCRYPT_CAST_256, MCRYPT_MODE_CBC), MCRYPT_RAND));
				$_POST['sftp_pass'] = openssl_encrypt($_POST['sftp_pass'], 'AES-256-CBC', file_get_contents(HASH), false, base64_decode($iv));
			
			}
			
		if($key == 'backup_disk' || $key == 'backup_space' && !is_numeric($val))
			exit();
	
	}

?>