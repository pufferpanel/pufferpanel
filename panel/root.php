<?php

/*
  PufferPanel - A Minecraft Server Management Panel
  Copyright (c) 2014 PufferPanel

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

$this->respond('GET', '/logout', function($request, $response, $service, $app) {
	if($response->isSent()) {
		return;
	}

	if($app->isLoggedIn) {

		/*
		 * Expire Session and Clear Database Details
		 */
		setcookie("pp_auth_token", null, time() - 86400, '/');
		setcookie("pp_server_node", null, time() - 86400, '/');
		setcookie("pp_server_hash", null, time() - 86400, '/');

		$logout = ORM::forTable('users')->where(array('session_id' => $request->cookies()['pp_auth_token'], 'session_ip' => $request->server()['REMOTE_ADDR']))->findOne();
		$logout->session_id = null;
		$logout->session_ip = null;
		$logout->save();

		$app->core->log->getUrl()->addLog(0, 1, array('auth.user_logout', 'Account logged out from ' . $request->server()['REMOTE_ADDR'] . '.'));
	}

	$response->redirect('/index');
});

$this->respond('GET', '/[index]', function($request, $response, $service, $app) {
	if($response->isSent()) {
		return;
	}

	echo $app->twig->render('panel/index.html', array(
		'xsrf' => $app->core->auth->XSRF(),
		'footer' => array(
			'seconds' => number_format((microtime(true) - $app->pageStartTime), 4)
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
			$response->redirect('/index?error=true', 302)->send();
		}
	}
});

$this->respond('GET', '/servers', function($request, $response, $service, $app) {
	if($response->isSent()) {
		return;
	}

	$core = $app->core;

	/*
	 * Redirect
	 */
	if($request->param('goto') != null && !empty($request->param('goto'))) {
		$core->server->nodeRedirect($request->param('goto'), $core->user->getData('id'), $core->user->getData('root_admin'));
	}

	/*
	 * Get the Servers
	 */
	if($core->user->getData('root_admin') == '1') {

		$servers = ORM::forTable('servers')->select('servers.*')->select('nodes.node', 'node_name')
				->join('nodes', array('servers.node', '=', 'nodes.id'))
				->orderByDesc('active')
				->findArray();
	} else {

		$servers = ORM::forTable('servers')->select('servers.*')->select('nodes.node', 'node_name')
				->join('nodes', array('servers.node', '=', 'nodes.id'))
				->where(array('servers.owner_id' => $core->user->getData('id'), 'servers.active' => 1))
				->where_raw('servers.owner_id = ? OR servers.hash IN(?)', array($core->user->getData('id'), join(',', $core->user->listServerPermissions())))
				->findArray();
	}

	/*
	 * List Servers
	 */
	echo $app->twig->render(
			'panel/servers.html', array(
		'servers' => $servers,
		'footer' => array(
			'seconds' => number_format((microtime(true) - $app->pageStartTime), 4)
		)
	));
});

$this->respond('GET', '/totp', function($request, $response, $service, $app) {
	echo $app->twig->render(
			'panel/totp.html', array(
		'totp' => array(
			'enabled' => $app->core->user->getData('use_totp')
		),
		'footer' => array(
			'seconds' => number_format((microtime(true) - $app->pageStartTime), 4)
		)
	));
});

$this->respond('GET', '/register', function($request, $response, $service, $app) {
	echo $app->twig->render(
		'panel/register.html', array(
			'xsrf' => $app->core->auth->XSRF(),
			'footer' => array(
				'seconds' => number_format((microtime(true) - $app->pageStartTime), 4)
			)
	));
});

$this->respond('POST', '/register', function($request, $response, $service, $app) {
	if($request->param('do') !== 'register') {
		
		$response->redirect('/register', 302);
		$response->send();
		return;

	}

	$token = $request->param('token');

	if($token == null || empty($token)) {
		$response->redirect('register?error=token', 302)->send();
		return;
	}

	list($encrypted, $iv) = explode('.', $request->param('token'));

	$core = $app->core;

	/* XSRF Check */
	if($core->auth->XSRF($request->param('xsrf')) !== true) {
		$response->redirect('index?error=xsrf&token=' . urlencode($token), 302)->send();
	}

	if(!preg_match('/^[\w-]{4,35}$/', $request->param('username'))) {
		$response->redirect('register?error=u_fail&token=' . urlencode($token), 302)->send();
		return;
	}

	if(strlen($request->param('password')) < 8 || $request->param('password') != $request->param('password_2')) {
		$response->redirect('register?error=p_fail&token=' . urlencode($token))->send();
		return;
	}

	$user = ORM::forTable('users')
			->where_any_is(array(
				array('username' => $request->param('username')),
				array('email' => $core->auth->decrypt($encrypted, $iv))))
			->findOne();

	if($user !== false) {
		$response->redirect('register?error=a_fail&token=' . $token, 302)->send();
		return;
	}

	$query = ORM::forTable('account_change')
			->where(array(
				'type' => 'subuser',
				'key' => $token,
				'verified' => 0
			))
			->findOne();

	if($query === false) {
		$response->redirect('register?error=t_fail&token=' . $token)->send();
		return;
	}

	$user = ORM::forTable('users')->create();
	$user->set(array(
		'uuid' => $core->auth->gen_UUID(),
		'username' => $request->param('username'),
		'email' => $core->auth->decrypt($encrypted, $iv),
		'password' => $core->auth->hash($request->param('password')),
		'permissions' => $row['content'],
		'language' => $core->settings->get('default_language'),
		'time' => time()
	));
	$user->save();

	$server = ORM::forTable('servers')
			->selectMany('subusers', 'hash')
			->where('hash', key(json_decode($row['content'], true)))
			->findOne();
	$subusers = json_decode($server->subusers, true);
	unset($subusers[$core->auth->decrypt($encrypted, $iv)]);
	$subusers[$user->id()] = "verified";

	$server->subusers = json_encode($subusers);
	$query->verified = 1;

	$server->save();
	$query->save();

	$response->redirect('index?registered', 302)->send();
});
