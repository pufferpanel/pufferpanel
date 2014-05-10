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
		
require_once("core/captcha/recaptchalib.php");

$statusMessage = null;
$noShow = false;

if(isset($_GET['do']) && $_GET['do'] == 'recover'){

	$resp = recaptcha_check_answer($core->settings->get('captcha_priv'), $_SERVER["REMOTE_ADDR"], @$_POST["recaptcha_challenge_field"], @$_POST["recaptcha_response_field"]);
		
	if($resp->is_valid){
	
		/*
		 * Find User
		 */
		$query = $mysql->prepare("SELECT * FROM `users` WHERE `email` = ?");
		$query->execute(array($_POST['email']));
		
			if($query->rowCount() == 1){
	
				$pKey = $core->auth->keygen('30');
				
				$accountChangeInsert = $mysql->prepare("INSERT INTO `account_change` VALUES(NULL, NULL, 'password', :email, :pkey, :expires, 0)");
				$accountChangeInsert->execute(array(
					':email' => $_POST['email'],
					':pkey' => $pKey,
					':expires' => time() + 14400
				));
				
					/*
					 * Send Email
					 */
					$core->email->buildEmail('password_reset', array(
                        'IP_ADDRESS' => $_SERVER['REMOTE_ADDR'],
                        'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($_SERVER['REMOTE_ADDR']),
                        'PKEY' => $pKey
                    ))->dispatch($_POST['email'], $core->settings->get('company_name').' - Reset Your Password');
				
                $core->log->getUrl()->addLog(0, 1, array('auth.password_reset_email', 'A password reset was requested and confimation emailed to your account email.'));
                
				$statusMessage = '<div class="alert alert-success">We have sent an email to the address you provided in the previous step. Please follow the instructions included in that email to continue. The verification key will expire in 4 hours.</div>';
				$noShow = true;
			
			}else{
                
			     $core->log->getUrl()->addLog(1, 0, array('auth.password_reset_email_fail', 'A password reset request was attempted but the email used was not found in the database. The email attempted was `'.$_POST['email'].'`.'));
                
				$statusMessage = '<div class="alert alert-danger">We couldn\'t find that email in our database.</div>';
			
			}
	
	}else{
	
		$statusMessage = '<div class="alert alert-danger">The spam prevention was not filled out correctly. Please try it again.</div>';
	
	}

}else if(isset($_GET['key'])){

	/*
	 * Change Password
	 */
	$key = $_GET['key'];
	$query = $mysql->prepare("SELECT * FROM `account_change` WHERE `key` = :key AND `verified` = '0' AND `time` > :time");
	$query->execute(array(
		':key' => $_GET['key'],
		':time' => time()
	));
		
		if($query->rowCount() ==  1){
		
			$row = $query->fetch();
			
			$raw_newpassword = $core->auth->keygen('12');
			
			$updateAccountChange = $mysql->prepare("UPDATE `account_change` SET `verified` = 1 WHERE `key` = ?");
			$updateAccountChange->execute(array($key));
			
			$updateUsers = $mysql->prepare("UPDATE `users` SET `password` = :newpass WHERE `email` = :email");
			$updateUsers->execute(array(
				':newpass' => $core->auth->hash($raw_newpassword),
				':email' => $row['content']
			));
			
            $core->log->getUrl()->addLog(0, 1, array('auth.password_reset', 'Your account password was successfull reset from the password reset form.'));
            
			$statusMessage = '<div class="alert alert-success">You should recieve an email within the next 5 minutes (usually instantly) with your new account password. We suggest changing this once you log in.</div>';
			$noShow = true;
		
				/*
				 * Send Email
				 */
				$core->email->buildEmail('new_password', array(
                    'NEW_PASS' => $raw_newpassword,
                    'EMAIL' => $row['content']
                ))->dispatch($row['content'], $core->settings->get('company_name').' - New Password');
		
		}else{
		
            $core->log->getUrl()->addLog(1, 0, array('auth.password_reset_fail', 'A password reset request was attempted but failed to be verified.'));
            
			$statusMessage = '<div class="alert alert-danger">Unable to verify password recovery request.<br />Did the key expire? Please contact support for more help or try again.</div>';
		
		}
		
}

echo $twig->render(
		'panel/password.html', array(
			'status' => $statusMessage,
			'noshow' => $noShow,
			'footer' => array(
				'queries' => Database\databaseInit::getCount(),
				'seconds' => number_format((microtime(true) - $pageStartTime), 4)
			)
	));

?>