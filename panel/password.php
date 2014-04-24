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

require_once("core/captcha/recaptchalib.php");
$statusMessage = ''; $noShow = false;

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

?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('assets/include/header.php'); ?>
	<title><?php echo $core->settings->get('company_name'); ?> - Reset Password</title>
	<script type="text/javascript">
		var RecaptchaOptions = {
			theme : 'custom',
			custom_theme_widget: 'recaptcha_widget'
		};
	</script>
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
				<a class="navbar-brand" href="<?php echo $core->settings->get('master_url'); ?>"><?php echo $core->settings->get('company_name'); ?></a>
			</div>
		</div>
		<div class="row">
			<div class="col-3">&nbsp;</div>
			<div class="col-6">
				<form action="password.php?do=recover" method="POST" id="login-form">
					<legend><?php echo $_l->tpl('reset.reset_h1'); ?></legend>
					<fieldset>
						<?php 
							echo $statusMessage;
							if($noShow === false){
						?>
						<div class="form-group">
							<label for="email" class="control-label"><?php echo $_l->tpl('string.email'); ?></label>
							<div>
								<input type="text" class="form-control" name="email" placeholder="<?php echo $_l->tpl('string.email'); ?>" />
							</div>
						</div>
						<div class="form-group">
							<label for="recaptcha_response_field" class="control-label"><?php echo $_l->tpl('string.spam_protection'); ?> <a href="javascript:Recaptcha.reload()"><?php echo $_l->tpl('string.refresh'); ?></a></label>
							<div>
								<div class="col-4" style="padding-left: 0;">
									<input type="text" class="form-control" id="recaptcha_response_field" name="recaptcha_response_field"/>
								</div>
								<div class="col-8">
									<div id="recaptcha_image"></div>
								</div>
								<script type="text/javascript" src="http://www.google.com/recaptcha/api/challenge?k=<?php echo $core->settings->get('captcha_pub'); ?>"></script>
								<noscript>
									<iframe src="http://www.google.com/recaptcha/api/noscript?k=<?php echo $core->settings->get('captcha_pub'); ?>"
								height="300" width="500" frameborder="0"></iframe><br>
									<textarea name="recaptcha_challenge_field" rows="3" cols="40"></textarea>
									<input type="hidden" name="recaptcha_response_field" value="manual_challenge">
								</noscript>
							</div>
						</div>
						<div class="form-group">
							<div>
								<input type="submit" class="btn btn-primary" value="<?php echo $_l->tpl('string.reset_password'); ?>" />
								<button class="btn btn-default" onclick="window.location='index.php';return false;"><?php echo $_l->tpl('string.login'); ?></button>
							</div>
						</div>
						<?php } ?>
					</fieldset>
				</form>
			</div>
			<div class="col-3">&nbsp;</div>
		</div>
		<div class="footer">
			<?php include('assets/include/footer.php'); ?>
		</div>
	</div>
</body>
</html>