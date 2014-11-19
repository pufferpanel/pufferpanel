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

$this->respond('GET', '/index', function($request, $response, $service, $app)  {
	if($response->isSent()) {
		return;
	}	

	$core = $app->core;
	$twig = $app->twig;
	$pageStartTime = $app->pageStartTime;

	echo $twig->render('panel/index.html', array(
		'xsrf' => $core->auth->XSRF(),
		'footer' => array(
			'seconds' => number_format((microtime(true) - $pageStartTime), 4)
		)
	));
});

$this->respond('POST', '/index', function($request, $response, $service, $app) {
	if($response->isSent()) {
		return;
	}

	$core = $app->core;

	if($request->param('totp') !== null && $request->param('check') !== null) {

		if(empty($request->param('totp')) || empty($request->param('check'))) {
			echo false;
		} else {

			$totp = ORM::forTable('users')->select('use_totp')->where('email', $request->param('check'))->findOne();

			if($totp === false) {
				echo false;
			} else {
				echo ($totp->use_totp == 1) ? true : false;
			}

		}

	} else if($request->param('do') == 'login') {

		/* XSRF Check */
		if($core->auth->XSRF($request->param('xsrf')) !== true) {
			$response->redirect('/index?error=token', 302)->send();
			return;
		}
			

		/*
		 * Get the Account Details
		 */
		$account = ORM::forTable('users')->where('email', $request->param('email'))->findOne();

		if($core->auth->verifyPassword($request->param('email'), $request->param('password'))) {

			/*
			 * Validate TOTP Key
			 */
			if($account->use_totp == 1) {

				if(!$core->auth->validateTOTP($request->param('totp_token'), $account->totp_secret)) {
					$core->log->getUrl()->addLog(0, 1, array('auth.account_login_fail_totp', 'A failed attempt to login to the account was made from ' . $request->server()['REMOTE_ADDR'] . '. The login failed due to TOTP 2FA mismatch.'));
					$response->redirect('/index?totp=error', 302)->send();
					return;
				}

			}

			/*
			 * Account Exists
			 * Set Cookies and List Servers
			 */
			$token = $core->auth->keygen('12');
			$expires = ($request->param('remember_me') === null) ? (time() + 604800) : null;

			setcookie("pp_auth_token", $token, $expires, '/');

			$account->set(array(
				'session_id' => $token,
				'session_ip' => $request->server()['REMOTE_ADDR']
			));
			$account->save();

			if($account->notify_login_s == 1) {

				$core->email->generateLoginNotification('success', array(
					'IP_ADDRESS' => $request->server()['REMOTE_ADDR'],
					'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($request->server()['REMOTE_ADDR'])
				))->dispatch($request->param('email'), $core->settings->get('company_name') . ' - Account Login Notification');

			}

			$core->log->getUrl()->addLog(0, 1, array('auth.account_login', 'Account was logged in from ' . $request->server()['REMOTE_ADDR'] . '.', $account->id));
			$response->redirect('/servers', 302)->send();

		} else {

			if($account !== false) {

				if($account->notify_login_f == 1) {

					/*
					 * Send Email
					 */
					$core->email->generateLoginNotification('failed', array(
						'IP_ADDRESS' => $request->server()['REMOTE_ADDR'],
						'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($request->server()['REMOTE_ADDR'])
					))->dispatch($request->param('email'), $core->settings->get('company_name') . ' - Account Login Failure Notification');

				}
				
			}

			$core->log->getUrl()->addLog(0, 1, array('auth.account_login_fail', 'A failed attempt to login to the account was made from ' . $request->server()['REMOTE_ADDR'] . '.'));
			$request->param('/index?error=true', 302)->send();
			
		}
	}
});
