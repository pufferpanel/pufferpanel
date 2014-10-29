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
use \ORM as ORM;

require_once('../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token')) === true)
	Components\Page::redirect('servers.php');

/*
 * jQuery Call for TOTP
 */
if(isset($_POST['totp']) && isset($_POST['check'])){

	if(empty($_POST['totp']) || empty($_POST['check']))
		echo false;
	else{

		$totp = ORM::forTable('users')->select('use_totp')->where('email', $_POST['check'])->findOne();

		if($totp === false)
			echo false;
		else
			echo ($totp->use_totp == 1) ? true : false;

	}

}else if(isset($_GET['do']) && $_GET['do'] == 'login'){

	/* XSRF Check */
	if($core->auth->XSRF(@$_POST['xsrf']) !== true)
		Components\Page::redirect('index.php?error=token');

	/*
	* Get the Account Details
	*/
	$account = ORM::forTable('users')->where('email', $_POST['email'])->findOne();

	if($core->auth->verifyPassword($_POST['email'], $_POST['password']) === true){

		/*
		 * Validate TOTP Key
		 */
		if($account->use_totp == 1){

			if($core->auth->validateTOTP($_POST['totp_token'], $account->totp_secret) !== true){

				$core->log->getUrl()->addLog(0, 1, array('auth.account_login_fail_totp', 'A failed attempt to login to the account was made from '.$_SERVER['REMOTE_ADDR'].'. The login failed due to TOTP 2FA mismatch.'));
				Components\Page::redirect('index.php?totp=error');

			}

		}

		/*
		 * Account Exists
		 * Set Cookies and List Servers
		 */
		$token = $core->auth->keygen('12');
		$expires = (isset($_POST['remember_me'])) ? (time() + 604800) : null;

			setcookie("pp_auth_token", $token, $expires, '/');

			$account->set(array(
				'session_id' => $token,
				'session_ip' => $_SERVER['REMOTE_ADDR']
			));
			$account->save();

			if($account->notify_login_s == 1){

				$core->email->generateLoginNotification('success', array(
                    'IP_ADDRESS' => $_SERVER['REMOTE_ADDR'],
                    'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($_SERVER['REMOTE_ADDR'])
                ))->dispatch($_POST['email'], $core->settings->get('company_name').' - Account Login Notification');

            }

            $core->log->getUrl()->addLog(0, 1, array('auth.account_login', 'Account was logged in from '.$_SERVER['REMOTE_ADDR'].'.', $account->id));
			Components\Page::redirect('servers.php');

	}else{

		if($account !== false){

			if($account->notify_login_f == 1){

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
		Components\Page::redirect('index.php?error=true');

	}

}else{

	echo $twig->render(
			'panel/index.html', array(
				'xsrf' => $core->auth->XSRF(),
				'footer' => array(
					'seconds' => number_format((microtime(true) - $pageStartTime), 4)
				)
		));

}
?>
