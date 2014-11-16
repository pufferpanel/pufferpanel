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

$klein->respond('POST', '*', function($request, $response, $service) {
	$loggedIn = $core->auth->isLoggedIn($request->server()['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token')) !== true;
	$service->sharedData()['loggedIn'] = $loggedIn;
	if ($loggedIn) {
		$response->redirect('index.php?login', 302);
	}
});

$klein->respond('POST', '/subuser', function($request, $reponse, $service, $app) {
	if (!$service->sharedData()['loggedIn']) {
		return;
	}
	if ($request->param('token') === null) {
		$service->sharedData()['output'] = '<div class="alert alert-danger">The token you entered is invalid.</div>';
	} else {
		$core = $app->core;
		if ($core->auth->XSRF($request->param('xsrf_notify'), '_notify') !== true) {
			$service->sharedData()['output'] = '<div class="alert alert-danger">Unable to verify the token. Please reload the page and try again.</div>';
		} else {
			list($encrypted, $iv) = explode('.', $request->param('token'));

			if ($core->auth->decrypt($encrypted, $iv) != $core->user->getData('email')) {
				$service->sharedData()['output'] = '<div class="alert alert-danger">The token you entered is invalid.</div>';
			} else {
				$query = ORM::forTable('account_change')->select('content')->where(array('key' => $request->param('token'), 'verified' => 0))->findOne();

				if ($query === false) {
					$service->sharedData()['output'] = '<div class="alert alert-danger">The token you entered is invalid.</div>';
				} else {
					$_perms = json_decode($query->content, true);

					$info = ORM::forTable('servers')
							->select_many('servers.subusers', 'servers.hash', 'servers.name', 'users.permissions')
							->join('users', array('servers.owner_id', '=', 'users.id'))
							->where('hash', key(json_decode($row['content'], true)))
							->findOne();

					$subusers = json_decode($info->subusers, true);
					unset($subusers[$core->user->getData('email')]);
					$subusers[$core->user->getData('id')] = "verified";

					$perms = json_decode($info->permissions, true);
					$permissions = (is_array($perms)) ? $perms : array();
					$permissions[$info->hash] = $_perms[$info->hash];

					$info->permissions = json_encode($permissions);
					$info->subusers = json_encode($subusers);
					$query->verified = 1;

					$info->save();
					$query->save();

					$service->sharedData()['output'] = '<div class="alert alert-success">You have been added as a subuser for <em>' . $info->name . '</em>!</div>';
				}
			}
		}
	}
});

$klein->respond('POST', '/notifications', function($request, $reponse, $service, $app) {
	if (!$service->sharedData()['loggedIn']) {
		return;
	}
	$core = $app->core;
	if ($request->param('password') === null) {
		$service->sharedData()['output'] = '<div class="alert alert-danger">The password you entered is invalid.</div>';
	} else if ($core->auth->XSRF($request->param('xsrf_notify'), '_notify') !== true) {
		$service->sharedData()['output'] = '<div class="alert alert-danger">Unable to verify the token. Please reload the page and try again.</div>';
	} else {
		if ($core->auth->verifyPassword($core->user->getData('email'), $request->param('password')) === true) {
			$account = ORM::forTable('users')->findOne($core->user->getData('id'));
			$account->notify_login_s = $request->param('e_s');
			$account->notify_login_f = $request->param('e_f');
			$account->save();

			$core->log->getUrl()->addLog(0, 1, array('user.notifications_updated', 'The notification preferences for this account were updated.'));
			$service->sharedData()['output'] = '<div class="alert alert-success">Your notification preferences have been updated.</div>';
		} else {
			$core->log->getUrl()->addLog(1, 1, array('user.notifications_update_fail', 'The notification preferences for this account were unable to be updated because the supplied password was wrong.'));
			$service->sharedData()['output'] = '<div class="alert alert-danger">We were unable to verify your password. Please try again.</div>';
		}
	}
});

$klein->respond('POST', '/email', function($request, $reponse, $service, $app) {
	if (!$service->sharedData()['loggedIn']) {
		return;
	}
	$core = $app->core;
	if (!$core->auth->XSRF($request->param('xsrf_email'), '_email')) {
		$service->sharedData()['output'] = '<div class="alert alert-danger">Unable to verify the token. Please reload the page and try again.</div>';
	} else {
		if ($request->param('newemail') === null || $request->param('password') === null) {
			$service->sharedData()['output'] = '<div class="alert alert-danger">Not all variables were passed to the script.</div>';
		} else {
			if ($core->auth->verifyPassword($core->user->getData('email'), $request->param('password')) === true) {
				$account = ORM::forTable('users')->findOne($core->user->getData('id'));
				$account->email = $request->param('newemail');
				$account->save();

				$core->log->getUrl()->addLog(0, 1, array('user.email_updated', 'Your account email was updated.'));
				$service->sharedData()['output'] = '<div class="alert alert-success">Your email has been updated successfully.</div>';
			} else {
				$core->log->getUrl()->addLog(1, 1, array('user.email_update_fail', 'Your email was unable to be updated due to an incorrect password provided.'));
				$service->sharedData()['output'] = '<div class="alert alert-danger">We were unable to verify your password. Please try again.</div>';
			}
		}
	}
});

$klein->respond('POST', '/password', function($request, $reponse, $service, $app) {
	if (!$service->sharedData()['loggedIn']) {
		return;
	}
	$core = $app->core;
	if ($core->auth->XSRF($request->param('xsrf_pass'), '_pass') !== true) {
		$service->sharedData()['output'] = '<div class="alert alert-danger">Unable to verify the token. Please reload the page and try again.</div>';
	} else {

		if ($core->auth->verifyPassword($core->user->getData('email'), $request->param('p_password'))) {

			if (preg_match("#.*^(?=.{8,200})(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9]).*$#", $request->param('p_password_new'))) {

				if ($request->param('p_password_new') == $request->param('p_password_new_2')) {
					$account = ORM::forTable('users')->findOne($core->user->getData('id'));
					$account->password = $core->auth->hash($request->param('p_password_new'));
					$account->session_id = null;
					$account->session_ip = null;
					$account->save();

					$core->email->buildEmail('password_changed', array(
						'IP_ADDRESS' => $request->server()['REMOTE_ADDR'],
						'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($request->server()['REMOTE_ADDR'])
					))->dispatch($core->user->getData('email'), $core->settings->get('company_name') . ' - Password Change Notification');

					$core->log->getUrl()->addLog(0, 1, array('user.password_updated', 'Your account password was changed.'));
					$service->sharedData()['output'] = '<div class="alert alert-success">Your password has been sucessfully changed!</div>';
				} else {
					$service->sharedData()['output'] = '<div class="alert alert-danger">Your passwords did not match.</div>';
				}
			} else {
				$service->sharedData()['output'] = '<div class="alert alert-danger">Your password is not complex enough. Please make sure to include at least one number, and some type of mixed case.</div>';
			}
		} else {

			$core->log->getUrl()->addLog(1, 1, array('user.password_update_fail', 'Your password was unable to be changed because the current password was not entered correctly.'));
			$service->sharedData()['output'] = '<div class="alert alert-danger">Current account password is not correct.</div>';
		}
	}
});

$klein->respond('POST', '*', function($request, $response, $service, $app) {
	$core = $app->core;
	if ($core->user->getData('notify_login_s') == 1) {
		$ns1 = 'checked="checked"';
		$ns0 = '';
	} else {
		$ns0 = 'checked="checked"';
		$ns1 = '';
	}
	if ($core->user->getData('notify_login_f') == 1) {
		$nf1 = 'checked="checked"';
		$nf0 = '';
	} else {
		$nf0 = 'checked="checked"';
		$nf1 = '';
	}

	/*
	 * Get Notification Preferences
	 */
	if ($core->user->getData('notify_login_s') == 1) {
		$ns1 = 'checked="checked"';
		$ns0 = '';
	} else {
		$ns0 = 'checked="checked"';
		$ns1 = '';
	}
	if ($core->user->getData('notify_login_f') == 1) {
		$nf1 = 'checked="checked"';
		$nf0 = '';
	} else {
		$nf0 = 'checked="checked"';
		$nf1 = '';
	}

	return $app->twig->render(
					'panel/account.html', array(
				'output' => $service->sharedData()['output'],
				'xsrf' => array(
					'pass' => $core->auth->XSRF(null, '_pass'),
					'email' => $core->auth->XSRF(null, '_email'),
					'notify' => $core->auth->XSRF(null, '_notify')
				),
				'failed_login' => array(
					'e_f' => array("value" => 1, "checked" => $nf1),
					'e_f_2' => array("value" => 0, "checked" => $nf0)
				),
				'success_login' => array(
					'e_s' => array("value" => 1, "checked" => $ns1),
					'e_s_2' => array("value" => 0, "checked" => $ns0)
				),
				'footer' => array(
					'seconds' => number_format((microtime(true) - $app->startTime), 4)
				)
	));
});
