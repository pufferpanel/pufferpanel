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
require_once('core/framework/framework.core.php');
$error = '';

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token')) !== true){
	Page\components::redirect('index.php?login');
	exit();
}

/*
 * Lah-de-dah
 */
$outputMessage = null;

/*
 * Changing Account Details
 */
if(isset($_GET['action'])){

	if($_GET['action'] == 'notifications' && isset($_POST['password'])){
			
		if($core->auth->verifyPassword($core->user->getData('email'), $_POST['password']) === true){
		
			$updateUsers = $mysql->prepare("UPDATE `users` SET `notify_login_s` = :e_s, `notify_login_f` = :e_f WHERE `id` = :uid");
			$updateUsers->execute(array(
				':e_s' => $_POST['e_s'],
				':e_f' => $_POST['e_f'],
				':uid' => $core->user->getData('id')
			));
			
            $core->log->getUrl()->addLog(0, 1, array('user.notifications_updated', 'The notification preferences for this account were updated.'));
            
			$outputMessage = '<div class="alert alert-success">Your notification preferences have been updated.</div>';
		
		}else{
		
            $core->log->getUrl()->addLog(1, 1, array('user.notifications_update_fail', 'The notification preferences for this account were unable to be updated because the supplied password was wrong.'));
            
			$outputMessage = '<div class="alert alert-danger">We were unable to verify your password. Please try again.</div>';
		
		}

	}else if($_GET['action'] == 'email'){
	
		/*
		 * Update Email Address
		 */
		$emailKey = $core->auth->keygen('30');
		$expire = time() + 14400;
		
		if($_POST['newemail'] == $core->user->getData('email')){
		
			$outputMessage = '<div class="alert alert-danger">Sorry, you can\'t change your email to the email address you are currently using for the account, that wouldn\'t make sense!</div>';
		
		}else{
		
			if($core->auth->verifyPassword($core->user->getData('email'), $_POST['password']) === true){
					
				$updateEmail = $mysql->prepare("UPDATE `users` SET `email` = :email WHERE `id` = :id");
				$updateEmail->execute(array(
					':email' => $_POST['newemail'],
					':id' => $core->user->getData('id')
				));
				
                $core->log->getUrl()->addLog(0, 1, array('user.email_updated', 'Your account email was updated.'));
                
				$outputMessage = '<div class="alert alert-success">Your email has been updated successfully.</div>';
				
			}else{
			
                $core->log->getUrl()->addLog(1, 1, array('user.email_update_fail', 'Your email was unable to be updated due to an incorrect password provided.'));
                
				$outputMessage = '<div class="alert alert-danger">We were unable to verify your password. Please try again.</div>';
			
			}
				
		}
	
	}else if($_GET['action'] == 'password'){
	
		if($core->auth->verifyPassword($core->user->getData('email'), $_POST['p_password']) === true){
		
			if(preg_match("#.*^(?=.{8,200})(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9]).*$#", $_POST['p_password_new'])){
			
				if($_POST['p_password_new'] == $_POST['p_password_new_2']){
				
					$newPassword = $core->auth->hash($_POST['p_password_new']);
					
						/*
						 * Change Password
						 */
						$updatePassword = $mysql->prepare("UPDATE `users` SET `password` = :password, `session_id` = NULL, `session_ip` = NULL, `session_expires` = NULL WHERE `id` = :uid");
						$updatePassword->execute(array(
							':password' => $core->auth->hash($_POST['p_password_new']),
							':uid' => $core->user->getData('id')
						));
						
							
						/*
						 * Send Email
						 */
						$message = $core->email->buildEmail('password_changed', array(
                            'IP_ADDRESS' => $_SERVER['REMOTE_ADDR'],
                            'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($_SERVER['REMOTE_ADDR'])
                        ))->dispatch($core->user->getData('email'), $core->settings->get('company_name').' - Password Change Notification');
						
                    $core->log->getUrl()->addLog(0, 1, array('user.password_updated', 'Your account password was changed.'));
                    
					$outputMessage = '<div class="alert alert-success">Your password has been sucessfully changed!</div>';
						
				
				}else{
                    
					$outputMessage = '<div class="alert alert-danger">Your passwords did not match.</div>';
				
				}
			
			}else{
			
				$outputMessage = '<div class="alert alert-danger">Your password is not complex enough. Please make sure to include at least one number, and some type of mixed case.</div>';
			
			}
		
		}else{
		
            $core->log->getUrl()->addLog(1, 1, array('user.password_update_fail', 'Your password was unable to be changed because the current password was not entered correctly.'));
            
			$outputMessage = '<div class="alert alert-danger">Current account password is not correct.</div>';
		
		}
	
	}else{
	
		$outputMessage = '<div class="alert alert-danger">Invalid parameters passed. Did you fill out all required fields?</div>';
	
	}

}

/*
 * Get Notification Preferences
 */
if($core->user->getData('notify_login_s') == 1){ $ns1 = 'checked="checked"'; $ns0 = ''; }else{ $ns0 = 'checked="checked"'; $ns1 = ''; }
if($core->user->getData('notify_login_f') == 1){ $nf1 = 'checked="checked"'; $nf0 = ''; }else{ $nf0 = 'checked="checked"'; $nf1 = ''; }

/*
 * Display Page
 */
echo $twig->render(
		'panel/account.html', array(
			'output' => $outputMessage,
			'failed_login' => array(
				'e_f' => array("value" => 1, "checked" => $ns1),
				'e_f_2' => array("value" => 0, "checked" => $ns0)
			),
			'success_login' => array(
				'e_s' => array("value" => 1, "checked" => $nf1),
				'e_s_2' => array("value" => 0, "checked" => $nf0)
			),
			'footer' => array(
				'queries' => Database\databaseInit::getCount(),
				'seconds' => number_format((microtime(true) - $pageStartTime), 4)
			)
	));
?>