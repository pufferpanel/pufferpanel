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

require_once("../src/captcha/recaptchalib.php");

$statusMessage = null;
$noShow = false;

if(isset($_GET['do']) && $_GET['do'] == 'recover'){

	/* XSRF Check */
	if($core->auth->XSRF(@$_POST['xsrf']) !== true)
		Components\Page::redirect('password.php?error=token');

	$resp = recaptcha_check_answer($core->settings->get('captcha_priv'), $_SERVER["REMOTE_ADDR"], @$_POST["recaptcha_challenge_field"], @$_POST["recaptcha_response_field"]);

	if($resp->is_valid){

		/*
		 * Find User
		 */
		$query = ORM::forTable('users')->where('email', $_POST['email'])->findOne();

			if($query !== false){

				$key = $core->auth->keygen('30');

				$account = ORM::forTable('account_change')->create();
				$account->set(array(
					'type' => 'password',
					'content' => $_POST['email'],
					'key' => $key,
					'time' => time() + 14400
				));
				$account->save();

					/*
					 * Send Email
					 */
					$core->email->buildEmail('password_reset', array(
                        'IP_ADDRESS' => $_SERVER['REMOTE_ADDR'],
                        'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($_SERVER['REMOTE_ADDR']),
                        'PKEY' => $key
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
	$query = ORM::forTable('account_change')->where(array('key' => $_GET['key'], 'verified' => 0))->where_gt('time', time())->findOne();

		if($query !== false){

			$password = $core->auth->keygen('12');
			$query->verified = 1;

			$user = ORM::forTable('users')->where('email', $query->content)->findOne();
			$user->password = $core->auth->hash($password);

			$user->save();
			$query->save();

            $core->log->getUrl()->addLog(0, 1, array('auth.password_reset', 'Your account password was successfull reset from the password reset form.'));

			$statusMessage = '<div class="alert alert-success">You should recieve an email within the next 5 minutes (usually instantly) with your new account password. We suggest changing this once you log in.</div>';
			$noShow = true;

				/*
				 * Send Email
				 */
				$core->email->buildEmail('new_password', array(
                    'NEW_PASS' => $password,
                    'EMAIL' => $query->content
                ))->dispatch($query->content, $core->settings->get('company_name').' - New Password');

		}else{

            $core->log->getUrl()->addLog(1, 0, array('auth.password_reset_fail', 'A password reset request was attempted but failed to be verified.'));
			$statusMessage = '<div class="alert alert-danger">Unable to verify password recovery request.<br />Did the key expire? Please contact support for more help or try again.</div>';

		}

}

echo $twig->render(
		'panel/password.html', array(
			'status' => $statusMessage,
			'noshow' => $noShow,
			'xsrf' => $core->auth->XSRF(),
			'footer' => array(
				'seconds' => number_format((microtime(true) - $pageStartTime), 4)
			)
	));

?>
