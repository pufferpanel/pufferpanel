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

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token')) === true){
	Page\components::redirect('servers.php');
}

if(isset($_GET['do']) && $_GET['do'] == 'login'){
	
        if(isset($_POST['redirect']) && !empty($_POST['redirect']))
            $postLoginURL = $_GET['redirect'];
        else
            $postLoginURL = 'servers.php';
    
			if($core->auth->verifyPassword($_POST['email'], $_POST['password']) === true){
			
				/*
				 * Account Exists
				 * Set Cookies and List Servers
				 */
				$token = $core->auth->keygen('12');
				$expires = (isset($_POST['remember_me'])) ? time() + 604800 : time() + 1800;
				$cookieExpire = (isset($_POST['remember_me'])) ? time() + 604800 : 0;
				
					setcookie("pp_auth_token", $token, $cookieExpire, '/', $core->settings->get('cookie_website'));
				
					$updateUsers = $mysql->prepare("UPDATE `users` SET `session_id` = :token, `session_ip` = :ipaddr, `session_expires` = :expires WHERE `email` = :email");
					$updateUsers->execute(array(
						':token' => $token,
						':ipaddr' => $_SERVER['REMOTE_ADDR'],
						':expires' => $expires,
						':email' => $_POST['email']
					));
				
					/*
					 * Send Email if Set
					 */
					$selectAccount = $mysql->prepare("SELECT * FROM `users` WHERE `email` = ?");
					$selectAccount->execute(array($_POST['email']));
					
					$row = $selectAccount->fetch();
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
	
}
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('assets/include/header.php'); ?>
	<title><?php echo $core->settings->get('company_name'); ?> - Login</title>
</head>
<body>
	<div class="container">
		<div class="pull-right" style="margin-top: -26px;">
			<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=ar"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/ar.png" /></a>
			<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=da"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/da.png" /></a>
			<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=de"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/de.png" /></a>
			<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=en"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/en.png" /></a>
			<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=es"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/es.png" /></a>
			<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=fr"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/fr.png" /></a>
			<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=ja"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/ja.png" /></a>
			<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=nl"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/nl.png" /></a>
			<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=pt"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/pt.png" /></a>
			<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=se"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/se.png" /></a>
			<a href="<?php echo $core->settings->get('master_url'); ?>core/ajax/set_language.php?language=zh"><img src="<?php echo $core->settings->get('assets_url'); ?>flags/zh.png" /></a>
		</div>
		<div class="navbar navbar-default">
			<div class="navbar-header">
				<a class="navbar-brand" href="#"><?php echo $core->settings->get('company_name'); ?></a>
			</div>
		</div>
		<div class="row">
			<div class="col-3"></div>
			<div class="col-6">
				<form action="index.php?do=login" method="POST" id="login-form">
					<legend><?php echo $_l->tpl('login.login_h1'); ?></legend>
					<fieldset>
						<?php 
							if(isset($_GET['error'])){
								echo '<div class="alert alert-danger"><strong>Oh snap!</strong> The username or password you submitted was incorrect.</div>';
							}
						?>
						<div class="form-group">
							<label for="email" class="control-label"><?php echo $_l->tpl('string.email'); ?></label>
							<div>
								<input type="text" class="form-control" name="email" placeholder="<?php echo $_l->tpl('string.email'); ?>" />
							</div>
						</div>
						<div class="form-group">
							<label for="login-password" class="control-label"><?php echo $_l->tpl('string.password'); ?></label>
							<div>
								<input type="password" class="form-control" name="password" placeholder="<?php echo $_l->tpl('string.password'); ?>" />
							</div>
						</div>
						<div class="form-group">
							<div>
								<input type="submit" class="btn btn-primary" value="Login" />
								<button class="btn btn-default" onclick="window.location='password.php';return false;"><?php echo $_l->tpl('string.reset_password'); ?></button>
							</div>
						</div>
					</fieldset>
				</form>
			</div>
			<div class="col-3"></div>
		</div>
		<div class="footer">
			<?php include('assets/include/footer.php'); ?>
		</div>
	</div>
</body>
</html>