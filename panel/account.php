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
namespace PufferPanel\Core;
use \ORM, \Unirest;

require_once('../src/core/core.php');
$error = '';

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token')) !== true){
	Components\Page::redirect('index.php?login');
	exit();
}

/*
 * Lah-de-dah
 */
$outputMessage = null;
$exception = false;

/*
 * Changing Account Details
 */
if(isset($_GET['action'])){

	if($_GET['action'] == 'subuser' && isset($_POST['token'])){

		/* XSRF Check */
		if($core->auth->XSRF(@$_POST['xsrf_notify'], '_notify') !== true)
			$outputMessage = '<div class="alert alert-danger">Unable to verify the token. Please reload the page and try again.</div>';
		else {

			$query = ORM::forTable('account_change')->select_many('id', 'verified', 'content')->where(array('key' => $_POST['token'], 'verified' => 0))->findOne();

			if(!$query) {
				$outputMessage = '<div class="alert alert-danger">The token you entered is invalid.</div>';
			} else {

				$_perms = json_decode($query->content, true);
				$info = ORM::forTable('servers')
					->select_many('servers.*', 'users.permissions', array('uid' => 'users.id'), 'nodes.ip', array('node_gsd_secret' => 'nodes.gsd_secret'), 'nodes.gsd_listen')
					->join('users', array('servers.owner_id', '=', 'users.id'))
					->join('nodes', array('servers.node', '=', 'nodes.id'))
					->where('hash', key($_perms))
					->findOne();

				try {

					$request = Unirest::put(
						'http://' . $info->ip . ':' . $info->gsd_listen . '/gameservers/' . $info->gsd_id,
						array(
							"X-Access-Token" => $info->node_gsd_secret
						),
						array(
							"keys" => json_encode(array(
								$_perms[$info->hash]['key'] => $_perms[$info->hash]['perms_gsd']
							))
						)
					);

				} catch(\Exception $e) {

					\Tracy\Debugger::log($e);
					$exception = true;
					$outputMessage = '<div class="alert alert-danger">The server management daemon is not responding, we were unable to add your permissions. Please try again later.</div>';

				}

				if($exception === false) {

					$subusers = json_decode($info->subusers, true);
					unset($subusers[$core->user->getData('email')]);
					$subusers[$core->user->getData('id')] = "verified";

					$permissions = @json_decode($info->permissions, true);
					$permissions = (is_array($permissions)) ? $permissions : array();
					$permissions[$info->hash] = $_perms[$info->hash];

					// set permissions for user
					$user = ORM::forTable('users')->findOne($info->uid);
					$user->permissions = json_encode($permissions);

					//set server subusers
					$info->subusers = json_encode($subusers);

					// expire key
					$query->verified = 1;

					// save
					$info->save();
					$user->save();
					$query->save();

					$outputMessage = '<div class="alert alert-success">You have been added as a subuser for <em>'.$info->name.'</em>!</div>';

				}

			}

		}

	}else if($_GET['action'] == 'notifications' && isset($_POST['password'])){

		/* XSRF Check */
		if($core->auth->XSRF(@$_POST['xsrf_notify'], '_notify') !== true)
			$outputMessage = '<div class="alert alert-danger">Unable to verify the token. Please reload the page and try again.</div>';
		else {

			if($core->auth->verifyPassword($core->user->getData('email'), $_POST['password']) === true){

				$account = ORM::forTable('users')->findOne($core->user->getData('id'));
				$account->notify_login_s = $_POST['e_s'];
				$account->notify_login_f = $_POST['e_f'];
				$account->save();

				$core->log->getUrl()->addLog(0, 1, array('user.notifications_updated', 'The notification preferences for this account were updated.'));
				$outputMessage = '<div class="alert alert-success">Your notification preferences have been updated.</div>';

			}else{

				$core->log->getUrl()->addLog(1, 1, array('user.notifications_update_fail', 'The notification preferences for this account were unable to be updated because the supplied password was wrong.'));
				$outputMessage = '<div class="alert alert-danger">We were unable to verify your password. Please try again.</div>';

			}

		}

	}else if($_GET['action'] == 'email'){

		/* XSRF Check */
		if($core->auth->XSRF(@$_POST['xsrf_email'], '_email') !== true)
			$outputMessage = '<div class="alert alert-danger">Unable to verify the token. Please reload the page and try again.</div>';
		else {

			/*
			 * Update Email Address
			 */
			$emailKey = $core->auth->keygen('30');
			$expire = time() + 14400;

			if(!isset($_POST['newemail'], $_POST['password']))
				$outputMessage = '<div class="alert alert-danger">Not all variables were passed to the script.</div>';
			else{

				if($_POST['newemail'] == $core->user->getData('email'))
					$outputMessage = '<div class="alert alert-danger">Sorry, you can\'t change your email to the email address you are currently using for the account, that wouldn\'t make sense!</div>';
				else{

					if($core->auth->verifyPassword($core->user->getData('email'), $_POST['password']) === true){

						$account = ORM::forTable('users')->findOne($core->user->getData('id'));
						$account->email = $_POST['newemail'];
						$account->save();

						$core->log->getUrl()->addLog(0, 1, array('user.email_updated', 'Your account email was updated.'));
						$outputMessage = '<div class="alert alert-success">Your email has been updated successfully.</div>';

					}else{

						$core->log->getUrl()->addLog(1, 1, array('user.email_update_fail', 'Your email was unable to be updated due to an incorrect password provided.'));
						$outputMessage = '<div class="alert alert-danger">We were unable to verify your password. Please try again.</div>';

					}

				}

			}

		}

	}else if($_GET['action'] == 'password'){

		/* XSRF Check */
		if($core->auth->XSRF(@$_POST['xsrf_pass'], '_pass') !== true)
			$outputMessage = '<div class="alert alert-danger">Unable to verify the token. Please reload the page and try again.</div>';
		else {

			if($core->auth->verifyPassword($core->user->getData('email'), $_POST['p_password']) === true){

				if(preg_match("#.*^(?=.{8,200})(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9]).*$#", $_POST['p_password_new'])){

					if($_POST['p_password_new'] == $_POST['p_password_new_2']){

						$newPassword = $core->auth->hash($_POST['p_password_new']);

							/*
							 * Change Password
							 */
							$account = ORM::forTable('users')->findOne($core->user->getData('id'));
							$account->password = $core->auth->hash($_POST['p_password_new']);
							$account->session_id = null;
							$account->session_ip = null;
							$account->save();

							/*
							 * Send Email
							 */
							$message = $core->email->buildEmail('password_changed', array(
								'IP_ADDRESS' => $_SERVER['REMOTE_ADDR'],
								'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($_SERVER['REMOTE_ADDR'])
							))->dispatch($core->user->getData('email'), $core->settings->get('company_name').' - Password Change Notification');

						$core->log->getUrl()->addLog(0, 1, array('user.password_updated', 'Your account password was changed.'));
						$outputMessage = '<div class="alert alert-success">Your password has been sucessfully changed!</div>';

					}else
						$outputMessage = '<div class="alert alert-danger">Your passwords did not match.</div>';

				}else
					$outputMessage = '<div class="alert alert-danger">Your password is not complex enough. Please make sure to include at least one number, and some type of mixed case.</div>';

			}else{

				$core->log->getUrl()->addLog(1, 1, array('user.password_update_fail', 'Your password was unable to be changed because the current password was not entered correctly.'));
				$outputMessage = '<div class="alert alert-danger">Current account password is not correct.</div>';

			}

		}

	}else
		$outputMessage = '<div class="alert alert-danger">Invalid parameters passed. Did you fill out all required fields?</div>';

}

/*
 * Display Page
 */
echo $twig->render(
		'panel/account.html', array(
			'output' => $outputMessage,
			'xsrf' => array(
				'pass' => $core->auth->XSRF(null, '_pass'),
				'email' => $core->auth->XSRF(null, '_email'),
				'notify' => $core->auth->XSRF(null, '_notify')
			),
			'notify_login_s' => $core->user->getData('notify_login_s'),
			'notify_login_f' => $core->user->getData('notify_login_f'),
			'footer' => array(
				'seconds' => number_format((microtime(true) - $pageStartTime), 4)
			)
	));
?>