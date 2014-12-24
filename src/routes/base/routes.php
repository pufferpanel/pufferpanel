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
use PufferPanel\Core, \ORM;

$klein->respond('GET', '/?', function($request, $response, $service, $app) {

	if($response->isSent()) {
		return;
	}

	if($app->isLoggedIn) {
		$response->redirect('/servers', 302)->send();
	} else {
		$response->redirect('/index')->send();
	}

});

$klein->respond('GET', '/account', function() use ($core) {
});

$klein->respond('POST', '/account', function() use ($core) {
});

$klein->with('/index', function() use ($klein, $core) {

	$klein->respond('GET', '', function($request, $response, $service) use ($core) {

		$response->body($core->twig->render('panel/index.html', array(
			'xsrf' => $core->auth->XSRF(),
			'flash' => $service->flashes()
		)));

	});

	$klein->respond('POST', '', function($request, $response, $service) use ($core) {

		if($core->auth->XSRF($request->param('xsrf')) !== true) {

			$service->flash('<div class="alert alert-warning">The XSRF token recieved was not valid. Please make sure cookies are enabled and try your request again.</div>');
			$response->redirect('/index')->send();

		}

		$account = ORM::forTable('users')->where('email', $request->param('email'))->findOne();

		if(!$core->auth->verifyPassword($request->param('email'), $request->param('password'))){

			if($account && $account->notify_login_f == 1) {

				$core->email->generateLoginNotification('failed', array(
					'IP_ADDRESS' => $request->ip(),
					'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($request->ip())
				))->dispatch($request->param('email'), $core->settings->get('company_name').' - Account Login Failure Notification');

			}

			$service->flash('<div class="alert alert-danger"><strong>Oh snap!</strong> The username or password you submitted was incorrect.</div>');
			$response->redirect('/index')->send();

		} else {

			if($account->use_totp == 1 && !$core->auth->validateTOTP($request->param('totp_token'), $account->totp_secret)){
				$service->flash('<div class="alert alert-danger"><strong>Oh snap!</strong> Your Two-Factor Authentication token was missing or incorrect.</div>');
				$response->redirect('/index')->send();
			}

			$cookie = (object) array('token' => $core->auth->keygen('12'), 'expires' => ($request->param('remember_me')) ? (time() + 604800) : null);

			$account->set(array(
				'session_id' => $cookie->token,
				'session_ip' => $request->ip()
			));
			$account->save();

			if($account->notify_login_s == 1){

				$core->email->generateLoginNotification('success', array(
					'IP_ADDRESS' => $request->ip(),
					'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($request->ip())
				))->dispatch($request->param('email'), $core->settings->get('company_name').' - Account Login Notification');

			}

			$response->cookie('pp_auth_token', $cookie->token, $cookie->expires);
			$response->redirect('/servers')->send();

		}

	});

	$klein->respond('POST', '/totp', function($request, $response) use ($core) {

		if(!$request->param('totp') || !$request->param('check')) {
			$response->body(false);
		} else {

			$totp = ORM::forTable('users')->select('use_totp')->where('email', $request->param('check'))->findOne();

			if(!$totp) {
				$response->body(false);
			} else {
				$response->body(($totp->use_totp == 1) ? true : false);
			}

		}

	});

});

$klein->respond('GET', '/logout', function() use ($core) {
	return 'Called /logout';
});

$klein->respond('POST', '/logout', function() use ($core) {
});

$klein->respond('GET', '/password', function() use ($core) {
});

$klein->respond('POST', '/password', function() use ($core) {
});

$klein->respond('GET', '/register', function() use ($core) {
});

$klein->respond('POST', '/register', function() use ($core) {
});

$klein->respond('GET', '/servers', function() use ($core) {
});

$klein->respond('POST', '/servers', function() use ($core) {
});