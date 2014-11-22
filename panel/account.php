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

$this->respond('POST', '/subuser', function($request, $response, $service, $app) {
	if($response->isSent()) {
		return;
	}
	
	if($request->param('token') === null) {
		$service->flash('<div class="alert alert-danger">The token you entered is invalid.</div>');
		return;
	}

	$core = $app->core;

	if(!$core->auth->XSRF($request->param('xsrf_notify'), '_notify')) {
		$service->flash('<div class="alert alert-danger">Unable to verify the token. Please reload the page and try again.</div>');
		return;
	}

	list($encrypted, $iv) = explode('.', $request->param('token'));

	if($core->auth->decrypt($encrypted, $iv) != $core->user->getData('email')) {
		$service->flash('<div class="alert alert-danger">The token you entered is invalid.</div>');
		return;
	}

	$query = ORM::forTable('account_change')
			->select('content')
			->where(array('key' => $request->param('token'), 'verified' => 0))
			->findOne();

	if(!$query) {
		$service->flash('<div class="alert alert-danger">The token you entered is invalid.</div>');
		return;
	}

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

	$service->flash('<div class="alert alert-success">You have been added as a subuser for <em>' . $info->name . '</em>!</div>');
});

$this->respond('POST', '/notifications', function($request, $response, $service, $app) {
	if($response->isSent()) {
		return;
	}

	$core = $app->core;

	if($request->param('password') === null) {
		$service->flash('<div class="alert alert-danger">The password you entered is invalid.</div>');
		return;
	} else if($core->auth->XSRF($request->param('xsrf_notify'), '_notify') !== true) {
		$service->flash('<div class="alert alert-danger">Unable to veriy the token. Please reload the page and try again.</div>');
		return;
	}

	if($core->auth->verifyPassword($core->user->getData('email'), $request->param('password')) === true) {

		$account = ORM::forTable('users')->findOne($core->user->getData('id'));
		$account->notify_login_s = $request->param('e_s');
		$account->notify_login_f = $request->param('e_f');
		$account->save();

		$core->log->getUrl()->addLog(0, 1, array('user.notifications_updated', 'The notification preferences for this account were updated.'));
		$service->flash('<div class="alert alert-success">Your notification preferences have been updated.</div>');

	} else {

		$core->log->getUrl()->addLog(1, 1, array('user.notifications_update_fail', 'The notification preferences for this account were unable to be updated because the supplied password was wrong.'));
		$service->flash('<div class="alert alert-danger">We were unable to verify your password. Please try again.</div>');

	}
});

$this->respond('POST', '/email', function($request, $response, $service, $app) {
	if($response->isSent()) {
		return;
	}

	$core = $app->core;
	
	if(!$core->auth->XSRF($request->param('xsrf_email'), '_email')) {
		$service->flash('<div class="alert alert-danger">Unable to verify the token. Please reload the page and try again.</div>');
		return;
	}

	if($request->param('newemail') === null || $request->param('password') === null) {
		$service->flash('<div class="alert alert-danger">Not all variables were passed to the script.</div>');
		return;
	}

	if($core->auth->verifyPassword($core->user->getData('email'), $request->param('password')) === true) {

		$account = ORM::forTable('users')->findOne($core->user->getData('id'));
		$account->email = $request->param('newemail');
		$account->save();

		$core->log->getUrl()->addLog(0, 1, array('user.email_updated', 'Your account email was updated.'));
		$service->flash('<div class="alert alert-success">Your email has been updated successfully.</div>');
		
	} else {

		$core->log->getUrl()->addLog(1, 1, array('user.email_update_fail', 'Your email was unable to be updated due to an incorrect password provided.'));
		$service->flash('<div class="alert alert-danger">We were unable to verify your password. Please try again.</div>');
		
	}
});

$this->respond('POST', '/password', function($request, $response, $service, $app) {
	if($response->isSent()) {
		return;
	}

	$core = $app->core;

	if($core->auth->XSRF($request->param('xsrf_pass'), '_pass') !== true) {
		$service->flash('<div class="alert alert-danger">Unable to verify the token. Please reload the page and try again.</div>');
		return;
	}

	if(!$core->auth->verifyPassword($core->user->getData('email'), $request->param('p_password'))) {

		$core->log->getUrl()->addLog(1, 1, array('user.password_update_fail', 'Your password was unable to be changed because the current password was not entered correctly.'));
		$service->flash('<div class="alert alert-danger">Current account password is not correct.</div>');
		return;
		
	}

	if(!preg_match("#.*^(?=.{8,200})(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9]).*$#", $request->param('p_password_new'))) {
		$service->flash('<div class="alert alert-danger">Your password is not complex enough. Please make sure to include at least one number, and some type of mixed case.</div>');
		return;
	}

	if($request->param('p_password_new') != $request->param('p_password_new_2')) {
		$service->flash('<div class="alert alert-danger">Your passwords did not match.</div>');
		return;
	}

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
	$service->flash('<div class="alert alert-success">Your password has been sucessfully changed!</div>');
});

$this->respond('POST', '*', function($request, $response) {
	if($response->isSent()) {
		return;
	}
	
	$response->redirect('/account', 302)->send();
});

$this->respond('GET', '*', function($request, $response, $service, $app) {
	if($response->isSent()) {
		return;
	}

	$core = $app->core;
	$twig = $app->twig;
	$startTime = $app->pageStartTime;

	echo $twig->render('panel/account.html', array(
				'output' => $service->flashes(),
				'xsrf' => array(
					'pass' => $core->auth->XSRF(null, '_pass'),
					'email' => $core->auth->XSRF(null, '_email'),
					'notify' => $core->auth->XSRF(null, '_notify')
				),
				'notify_login_s' => $core->user->getData('notify_login_s'),
				'notify_login_f' => $core->user->getData('notify_login_f'),
				'footer' => array(
					'seconds' => number_format((microtime(true) - $startTime), 4)
				)
	));
});