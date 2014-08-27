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
require_once('../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token')) === true)
	Page\components::redirect('servers.php');

/*
 * jQuery Call for TOTP
 */
if(isset($_POST['totp']) && isset($_POST['check'])){

	if(empty($_POST['totp']) || empty($_POST['check']))
		echo false;
	else{
	
		$checkTOTP = $mysql->prepare("SELECT `use_totp` FROM `users` WHERE `email` = :email");
		$checkTOTP->execute(array(
			'email' => $_POST['check']
		));
		
			if($checkTOTP->rowCount() != 1)
				echo false;
			else{
			
				$row = $checkTOTP->fetch();
				
				echo ($row['use_totp'] == 1) ? true : false;
			
			}
			
	}

}else if(isset($_GET['do']) && $_GET['do'] == 'login'){
	
        if(isset($_POST['redirect']) && !empty($_POST['redirect']))
            $postLoginURL = $_GET['redirect'];
        else
            $postLoginURL = 'servers.php';
    
			if($core->auth->verifyPassword($_POST['email'], $_POST['password']) === true){
				
				/*
				 * Get the Account Details
				 */
				$selectAccount = $mysql->prepare("SELECT * FROM `users` WHERE `email` = ?");
				$selectAccount->execute(array($_POST['email']));
				
				$row = $selectAccount->fetch();
				
				/*
				 * Validate TOTP Key
				 */
				if($row['use_totp'] == 1){
				
					if($core->auth->validateTOTP($_POST['totp_token'], $row['totp_secret']) !== true){
					
						$core->log->getUrl()->addLog(0, 1, array('auth.account_login_fail_totp', 'A failed attempt to login to the account was made from '.$_SERVER['REMOTE_ADDR'].'. The login failed due to TOTP 2FA mis-match.'));
						
						Page\components::redirect('index.php?totp=error');
					
					}
				
				}
				
				/*
				 * Account Exists
				 * Set Cookies and List Servers
				 */
				$token = $core->auth->keygen('12');
				$expires = (isset($_POST['remember_me'])) ? (time() + 604800) : 0;
				
					setcookie("pp_auth_token", $token, $expires, '/', $core->settings->get('cookie_website'));
				
					$updateUsers = $mysql->prepare("UPDATE `users` SET `session_id` = :token, `session_ip` = :ipaddr WHERE `email` = :email");
					$updateUsers->execute(array(
						':token' => $token,
						':ipaddr' => $_SERVER['REMOTE_ADDR'],
						':email' => $_POST['email']
					));
				
					/*
					 * Send Email if Set
					 */
					if($row['notify_login_s'] == 1){
						
						/*
						 * Send Email
						 */
						$core->email->generateLoginNotification('success', array(
                            'IP_ADDRESS' => $_SERVER['REMOTE_ADDR'],
                            'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($_SERVER['REMOTE_ADDR'])
                        ))->dispatch($_POST['email'], $core->settings->get('company_name').' - Account Login  Notification');
						    
                    }
                
                    $core->log->getUrl()->addLog(0, 1, array('auth.account_login', 'Account was logged in from '.$_SERVER['REMOTE_ADDR'].'.', $row['id']));
                
					Page\components::redirect($postLoginURL);
			
			}else{
				
				/*
				 * Spam Prevention
				 * This makes sure that the email is even in our system so we don't randomly email a user who doesn't have an account.
				 */
				$selectUser = $mysql->prepare("SELECT * FROM `users` WHERE `email` = ?");
				$selectUser->execute(array($_POST['email']));
				
					if($selectUser->rowCount() == 1){
			
						/*
						 * Send Email if Set
						 */
						$row = $selectUser->fetch();
						if($row['notify_login_f'] == 1){
																			
							/*
							 * Send Email
							 */								 
							$core->email->generateLoginNotification('failed', array(
                                'IP_ADDRESS' => $_SERVER['REMOTE_ADDR'],
                                'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($_SERVER['REMOTE_ADDR'])
                            ))->dispatch($_POST['email'], $core->settings->get('company_name').' - Account Login Failure Notification');															 
								    
						}
							  
					}
                
                $core->log->getUrl()->addLog(0, 1, array('auth.account_login_fail', 'A failed attempt to login to the account was made from '.$_SERVER['REMOTE_ADDR'].'.'));
                
				Page\components::redirect('index.php?error=true');
			
			}
	
}else{

	echo $twig->render(
			'panel/index.html', array(
				'footer' => array(
					'queries' => Database\databaseInit::getCount(),
					'seconds' => number_format((microtime(true) - $pageStartTime), 4)
				)
		));
	
}
?>