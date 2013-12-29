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

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token')) !== true){
	$core->framework->page->redirect('index.php');
	exit();
}

/*
 * Lah-de-dah
 */
$outputMessage = '';

/*
 * Changing Account Details
 */
if(isset($_GET['action'])){

	if($_GET['action'] == 'notifications' && isset($_POST['password'])){

		/*
		 * Update Notification Settings
		 */
		$selectAccount = $mysql->prepare("SELECT * FROM `users` WHERE `password` = :password AND `email` = :email");
		$selectAccount->execute(array(
			':password' => $core->framework->auth->encrypt($_POST['password']),
			':email' => $core->framework->user->getData('email')
		));
		
			if($selectAccount->rowCount() == 1){
			
				$updateUsers = $mysql->prepare("UPDATE `users` SET `notify_login_s` = :e_s, `notify_login_f` = :e_f WHERE `id` = :uid AND `password` = :password");
				$updateUsers->execute(array(
					':e_s' => $_POST['e_s'],
					':e_f' => $_POST['e_f'],
					':uid' => $core->framework->user->getData('id'),
					':password' => $core->framework->auth->encrypt($_POST['password'])
				));
				
                $core->framework->log->getUrl()->addLog(0, 1, array('user.notifications_updated', 'The notification preferences for this account were updated.'));
                
				$outputMessage = '<div class="alert alert-success">Your notification preferences have been updated.</div>';
			
			}else{
			
                $core->framework->log->getUrl()->addLog(1, 1, array('user.notifications_update_fail', 'The notification preferences for this account were unable to be updated because the supplied password was wrong.'));
                
				$outputMessage = '<div class="alert alert-danger">We were unable to verify your password. Please try again.</div>';
			
			}


	}else if($_GET['action'] == 'email'){
	
		/*
		 * Update Email Address
		 */
		$emailKey = $core->framework->auth->keygen('30');
		$expire = time() + 14400;
		
		if($_POST['newemail'] == $core->framework->user->getData('email')){
		
			$outputMessage = '<div class="alert alert-danger">Sorry, you can\'t change your email to the email address you are currently using for the account, that wouldn\'t make sense!</div>';
		
		}else{
		
			$selectAccount = $mysql->prepare("SELECT * FROM `users` WHERE `password` = :password AND `email` = :email");
			$selectAccount->execute(array(
				':password' => $core->framework->auth->encrypt($_POST['password']),
				':email' => $core->framework->user->getData('email')
			));
			
				if($selectAccount->rowCount() == 1){
						
						$updateEmail = $mysql->prepare("UPDATE `users` SET `email` = :email WHERE `id` = :id");
						$updateEmail->execute(array(
							':email' => $_POST['newemail'],
							':id' => $core->framework->user->getData('id')
						));
					
                    $core->framework->log->getUrl()->addLog(0, 1, array('user.email_updated', 'Your account email was updated.'));
                    
					$outputMessage = '<div class="alert alert-success">Your email has been updated successfully.</div>';
					
				}else{
				
                    $core->framework->log->getUrl()->addLog(1, 1, array('user.email_update_fail', 'Your email was unable to be updated due to an incorrect password provided.'));
                    
					$outputMessage = '<div class="alert alert-danger">We were unable to verify your password. Please try again.</div>';
				
				}
				
		}
	
	}else if($_GET['action'] == 'password'){
	
		/*
		 * Update Account Password
		 */
		$oldPassword = $core->framework->auth->encrypt($_POST['p_password']);
		 
			$selectAccount = $mysql->prepare("SELECT * FROM `users` WHERE `password` = :oldpass AND `email` = :email");
			$selectAccount->execute(array(
				':oldpass' => $core->framework->auth->encrypt($_POST['p_password']),
				':email' => $core->framework->user->getData('email')
			));
			
				if($selectAccount->rowCount() == 1){
					if(preg_match("#.*^(?=.{8,200})(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9]).*$#", $_POST['p_password_new'])){
					
						if($_POST['p_password_new'] == $_POST['p_password_new_2']){
						
							$newPassword = $core->framework->auth->encrypt($_POST['p_password_new']);
							
								/*
								 * Change Password
								 */
								$updatePassword = $mysql->prepare("UPDATE `users` SET `password` = :password, `session_id` = NULL, `session_ip` = NULL, `session_expires` = NULL WHERE `id` = :uid");
								$updatePassword->execute(array(
									':password' => $core->framework->auth->encrypt($_POST['p_password_new']),
									':uid' => $core->framework->user->getData('id')
								));
								
									
								/*
								 * Send Email
								 */
								$message = $core->framework->email->buildEmail('password_changed', array(
                                    'IP_ADDRESS' => $_SERVER['REMOTE_ADDR'],
                                    'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($_SERVER['REMOTE_ADDR'])
                                ))->dispatch($_POST['email'], $core->framework->settings->get('company_name').' - Password Change Notification');
								
                            $core->framework->log->getUrl()->addLog(0, 1, array('user.password_updated', 'Your account password was changed.'));
                            
							$outputMessage = '<div class="alert alert-success">Your password has been sucessfully changed!</div>';
								
						
						}else{
                            
							$outputMessage = '<div class="alert alert-danger">Your passowrds did not match.</div>';
						
						}
					
					}else{
					
						$outputMessage = '<div class="alert alert-danger">Your password is not complex enough. Please make sure to include at least one number, and some type of mixed case.</div>';
					
					}
				
				}else{
				
                    $core->framework->log->getUrl()->addLog(1, 1, array('user.password_update_fail', 'Your password was unable to be changed because the current password was not entered correctly.'));
                    
					$outputMessage = '<div class="alert alert-danger">Current account password is not correct.</div>';
				
				}
	
	}else{
	
		$outputMessage = '<div class="alert alert-danger">Invalid parameters passed. Did you fill out all required fields?</div>';
	
	}

}

/*
 * Get Notification Preferences
 */
$core->framework->user = new user($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'));
if($core->framework->user->getData('notify_login_s') == 1){ $ns1 = 'checked="checked"'; $ns0 = ''; }else{ $ns0 = 'checked="checked"'; $ns1 = ''; }
if($core->framework->user->getData('notify_login_f') == 1){ $nf1 = 'checked="checked"'; $nf0 = ''; }else{ $nf0 = 'checked="checked"'; $nf1 = ''; }

?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('assets/include/header.php'); ?>
	<title>PufferPanel - Your Settings</title>
</head>
<body>
	<div class="container">
		<div class="navbar navbar-default">
			<div class="navbar-header">
				<a class="navbar-brand" href="#"><?php echo $core->framework->settings->get('company_name'); ?></a>
			</div>
			<div class="navbar-collapse navbar-responsive-collapse collapse" style="height: 1px;">
				<ul class="nav navbar-nav navbar-right">
					<li class="dropdown">
						<a href="#" class="dropdown-toggle" data-toggle="dropdown">Account <b class="caret"></b></a>
							<ul class="dropdown-menu">
								<li><a href="logout.php">Logout</a></li>
								<?php if($core->framework->user->getData('root_admin') == 1){ echo '<li><a href="admin/index.php">Admin CP</a></li>'; } ?>
							</ul>
					</li>
				</ul>
			</div>
		</div>
		<div class="row">
			<div class="col-3">
				<div class="list-group">
					<a href="account.php" class="list-group-item active">Settings</a>
					<a href="servers.php" class="list-group-item">My Servers</a>
				</div>
			</div>
			<div class="col-9">
				<?php echo $outputMessage; ?>
				<div class="row">
					<div class="col-6">
						<h3 style="margin-top:0;">Change Password</h3><hr />
							<form action="account.php?action=password" method="post">
								<div class="form-group">
									<label for="email" class="control-label">Current Password</label>
									<div>
										<input type="password" class="form-control" name="p_password" />
									</div>
								</div>
								<div class="form-group">
									<label for="p_password_new" class="control-label">New Password</label>
									<div>
										<input type="password" class="form-control" name="p_password_new" />
									</div>
								</div>
								<div class="form-group">
									<label for="p_password_new_2" class="control-label">New Password (Again)</label>
									<div>
										<input type="password" class="form-control" name="p_password_new_2" />
									</div>
								</div>
								<div class="form-group">
									<div>
										<input type="submit" class="btn btn-primary" value="Change Password" />
									</div>
								</div>
							</form>
					</div>
					<div class="col-6">
						<h3 style="margin-top:0;">Update Account Email</h3><hr />
							<form action="account.php?action=email" method="post">
								<div class="form-group">
									<label for="newemail" class="control-label">New Email</label>
									<div>
										<input type="text" class="form-control" name="newemail" />
									</div>
								</div>
								<div class="form-group">
									<label for="password" class="control-label">Current Password</label>
									<div>
										<input type="password" class="form-control" name="password" />
									</div>
								</div>
								<div class="form-group">
									<div>
										<input type="submit" class="btn btn-primary" value="Update Email" />
									</div>
								</div>
							</form>
					</div>
				</div>
				<h3>Update Notification Preferences</h3><hr />
				<form action="account.php?action=notifications" method="post">
					<div class="col-6 nopad">
						<div class="form-group">
							<h4>Successful Login</h4>
							<div class="radio">
		                        <label for="e_s" class="alt-label"><input type="radio" id="e_s" name="e_s" value="1" <?php echo $ns1; ?>/>Please Email Me</label>
							</div>
							<div class="radio">
							    <label for="e_s_2" class="alt-label"><input type="radio" id="e_s_2" name="e_s" value="0" <?php echo $ns0; ?>/>Don't Email Me</label>
							</div>
						</div>
						<div class="form-group">
							<h4>Failed Login</h4>
							<div class="radio">
						        <label for="e_f" class="alt-label"><input type="radio" id="e_f" name="e_f" value="1" <?php echo $nf1; ?>/>Please Email Me</label>
							</div>
							<div class="radio">
							    <label for="e_f_2" class="alt-label"><input type="radio" id="e_f_2" name="e_f" value="0" <?php echo $nf0; ?>/>Don't Email Me</label>
							</div>
						</div>
						<div class="form-group">
							<label for="password" class="control-label">Current Password</label>
							<div>
								<input type="password" class="form-control" name="password" />
							</div>
						</div>
						<div class="form-group">
							<div>
								<input type="submit" class="btn btn-primary" value="Update Preferences" />
							</div>
						</div>
					</div>
				</form>	
			</div>
		</div>
		<div class="footer">
			<?php include('assets/include/footer.php'); ?>
		</div>
	</div>
</body>
</html>