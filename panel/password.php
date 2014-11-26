<?php

/*
	PufferPanel - A Minecraft Server Management Panel
	Copyright (c) 2014 Dane Everitt

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

require_once("../src/captcha/recaptchalib.php");

$klein->respond(function($request, $response, $service, $app) {
	if($app->isLoggedIn) {
		$response->redirect('/servers')->send();
	}
});

$klein->respond('POST', '*', function ($request, $response, $service) use ($core) {
	$noShow = false;

	if($request->param('do') === null) {

		$core->log->getUrl()->addLog(1, 0, array('auth.password_reset_fail', 'A password reset request was attempted but failed to be verified.'));
		$service->flash('<div class="alert alert-danger">Unable to verify password recovery request.<br />Did the key expire? Please contact support for more help or try again.</div>');

	} else if($request->param('do') === 'recover') {
		
		/* XSRF Check */
		if(!$core->auth->XSRF($request->param('xsrf'))) {
			$response->redirect('/password?error=token');
			return;
		}

		$resp = recaptcha_check_answer($core->settings->get('captcha_priv'), $request->server()["REMOTE_ADDR"], $request->param("recaptcha_challenge_field"), $request->param("recaptcha_response_field"));

		if($resp->is_valid) {

			/*
			 * Find User
			 */
			$query = ORM::forTable('users')->where('email', $request->param('email'))->findOne();

			if($query !== false) {

				$key = $core->auth->keygen('30');

				$account = ORM::forTable('account_change')->create();
				$account->set(array(
					'type' => 'password',
					'content' => $request->param('email'),
					'key' => $key,
					'time' => time() + 14400
				));
				$account->save();

				/*
				 * Send Email
				 */
				$core->email->buildEmail('password_reset', array(
					'IP_ADDRESS' => $request->server()["REMOTE_ADDR"],
					'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($request->server()["REMOTE_ADDR"]),
					'PKEY' => $key
				))->dispatch($request->param('email'), $core->settings->get('company_name') . ' - Reset Your Password');

				$core->log->getUrl()->addLog(0, 1, array('auth.password_reset_email', 'A password reset was requested and confimation emailed to your account email.'));

				$statusMessage = '<div class="alert alert-success">We have sent an email to the address you provided in the previous step. Please follow the instructions included in that email to continue. The verification key will expire in 4 hours.</div>';
				$noShow = true;

			} else {

				$core->log->getUrl()->addLog(1, 0, array('auth.password_reset_email_fail', 'A password reset request was attempted but the email used was not found in the database. The email attempted was `' . $request->param('email') . '`.'));
				$service->flash('<div class="alert alert-danger">We couldn\'t find that email in our database.</div>');
				
			}

		} else {

			$statusMessage = '<div class="alert alert-danger">The spam prevention was not filled out correctly. Please try it again.</div>';

		}

	} else if($request->param('key') !== null) {

		/*
		 * Change Password
		 */
		$key = $request->param('key');
		$query = ORM::forTable('account_change')->where(array('key' => $key, 'verified' => 0))->where_gt('time', time())->findOne();

		if($query !== false) {

			$password = $core->auth->keygen('12');
			$query->verified = 1;

			$user = ORM::forTable('users')->where('email', $query->content)->findOne();
			$user->password = $core->auth->hash($password);

			$user->save();
			$query->save();

			$core->log->getUrl()->addLog(0, 1, array('auth.password_reset', 'Your account password was successfull reset from the password reset form.'));

			$service->flash('<div class="alert alert-success">You should recieve an email within the next 5 minutes (usually instantly) with your new account password. We suggest changing this once you log in.</div>');
			$noShow = true;

			/*
			 * Send Email
			 */
			$core->email->buildEmail('new_password', array(
				'NEW_PASS' => $password,
				'EMAIL' => $query->content
			))->dispatch($query->content, $core->settings->get('company_name') . ' - New Password');

		} else {

			$core->log->getUrl()->addLog(1, 0, array('auth.password_reset_fail', 'A password reset request was attempted but failed to be verified.'));
			$service->flash('<div class="alert alert-danger">Unable to verify password recovery request.<br />Did the key expire? Please contact support for more help or try again.</div>');
			
		}
		
		$response->redirect('/password?success=' . $noShow, 302)->send();
		
	}
});

$klein->respond('GET', '*', function($request, $response, $service) use ($core, $twig, $pageStartTime) {
	echo $twig->render('panel/password.html', array(
				'status' => $service->flash(),
				'noshow' => $request->param('success', false),
				'xsrf' => $core->auth->XSRF(),
				'footer' => array(
					'seconds' => number_format((microtime(true) - $pageStartTime), 4)
				)
	));
});