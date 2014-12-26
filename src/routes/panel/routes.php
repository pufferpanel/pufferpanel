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
use \ORM;

$klein->respond('GET', '/account', function($request, $response, $service) use ($core) {

	$response->body($core->twig->render('panel/account.html', array(
		'xsrf' => $core->auth->XSRF(),
		'flash' => $service->flashes(),
		'notify_login_s' => $core->user->getData('notify_login_s'),
		'notify_login_f' => $core->user->getData('notify_login_f')
	)))->send();

});

$klein->respond('POST', '/account/update/[:action]', function($request, $response, $service) use ($core) {

	$core->routes = new Router\Router_Controller('Account', $core->user);
	$core->routes = $core->routes->loadClass();

	if(!$core->auth->XSRF($request->param('xsrf'))) {

		$service->flash('<div class="alert alert-warning"> The XSRF token recieved was not valid. Please make sure cookies are enabled and try your request again.</div>');
		$response->redirect('/account')->send();

	}

	// Update Email Address
	if($request->param('action') == "email") {

		if(!$request->param('newemail') || !$request->param('password')) {

			$service->flash('<div class="alert alert-danger">Not all variables were passed to the script.</div>');
			$response->redirect('/account')->send();

		}

		if(!$core->auth->verifyPassword($core->user->getData('email'), $request->param('password'))) {

			$service->flash('<div class="alert alert-danger">We were unable to verify your account password. Please try your request again.</div>');
			$response->redirect('/account')->send();

		}

		if($request->param('newemail') == $core->user->getData('email')) {

			$service->flash('<div class="alert alert-danger">Sorry, you can\'t change your email to the email address you are currently using for the account, that wouldn\'t make sense!</div>');
			$response->redirect('/account')->send();

		}

		if(!filter_var($request->param('newemail'), FILTER_VALIDATE_EMAIL)) {

			$service->flash('<div class="alert alert-danger">The email you provided is not valid.</div>');
			$response->redirect('/account')->send();

		}

		$account = ORM::forTable('users')->findOne($core->user->getData('id'));
		$account->email = $request->param('newemail');
		$account->save();

		$service->flash('<div class="alert alert-success">Your email address has been successfully updated to <em>'.$request->param('newemail').'</em>.</div>');
		$response->redirect('/account')->send();

	}

	// Update Account Password
	if($request->param('action') == "password") {

		if(!$core->auth->verifyPassword($core->user->getData('email'), $request->param('p_password'))){

			$service->flash('<div class="alert alert-danger">We were unable to verify your account password. Please try your request again.</div>');
			$response->redirect('/account')->send();

		}

		if($request->param('p_password_new') != $request->param('p_password_new_2')) {

			$service->flash('<div class="alert alert-danger">Your passwords did not match.</div>');
			$response->redirect('/account')->send();

		}

		if(!$core->routes->validatePasswordRequirements($request->param('p_password_new'))) {

			$service->flash('<div class="alert alert-danger">Your password is not complex enough. Please make sure to include at least one number, and some type of mixed case. Your new password must also be at least 8 characters long.</div>');
			$response->redirect('/account')->send();

		}

		if(!$core->routes->updatePassword($request->param('p_password'), $request->param('p_password_new'))) {

			$service->flash('<div class="alert alert-danger">An unhandled error was encountered when trying to update your password.</div>');
			$response->redirect('/account')->send();

		} else {

			$service->flash('<div class="alert alert-success">Your password has been sucessfully changed! Please login again using your new password.</div>');
			$response->redirect('/auth/login')->send();

		}

	}

	// Update Account Notitification Preferences
	if($request->param('action') == "notifications") {

		if(!$core->auth->verifyPassword($core->user->getData('email'), $request->param('password'))) {

			$service->flash('<div class="alert alert-danger">We were unable to verify your account password. Please try your request again.</div>');
			$response->redirect('/account')->send();

		}

		$account = ORM::forTable('users')->findOne($core->user->getData('id'));
		$account->notify_login_s = $request->param('e_s');
		$account->notify_login_f = $request->param('e_f');
		$account->save();

		$service->flash('<div class="alert alert-success">Your notification preferences have been updated.</div>');
		$response->redirect('/account')->send();

	}

	// Add User as a Subuser for a Server
	if($request->param('action') == "subuser") {

		$query = ORM::forTable('account_change')->select_many('id', 'verified', 'content')->where(array('key' => $request->param('token'), 'verified' => 0))->findOne();

		if(!$query) {

			$service->flash('<div class="alert alert-danger">The token you entered is invalid.</div>');
			$response->redirect('/account')->send();
			return;

		}

		$_perms = json_decode($query->content, true);
		$info = ORM::forTable('servers')
			->select_many('servers.*', 'users.permissions', 'nodes.ip', array('node_gsd_secret' => 'nodes.gsd_secret'), 'nodes.gsd_listen')
			->join('users', array('servers.owner_id', '=', 'users.id'))
			->join('nodes', array('servers.node', '=', 'nodes.id'))
			->where('hash', key($_perms))
			->findOne();

		$subusers = json_decode($info->subusers, true);
		if(!array_key_exists($core->user->getData('email'), $subusers)) {

			$service->flash('<div class="alert alert-danger">The token you entered is not valid for this email address.</div>');
			$response->redirect('/account')->send();

		}

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

			unset($subusers[$core->user->getData('email')]);
			$subusers[$core->user->getData('id')] = "verified";

			$permissions = @json_decode($info->permissions, true);
			$permissions = (is_array($permissions)) ? $permissions : array();
			$permissions[$info->hash] = $_perms[$info->hash];

			// set permissions for user
			$user = ORM::forTable('users')->findOne($core->user->getData('id'));
			$user->permissions = json_encode($permissions);

			//set server subusers
			$info->subusers = json_encode($subusers);

			// expire key
			$query->verified = 1;

			// save
			$info->save();
			$user->save();
			$query->save();

			$service->flash('<div class="alert alert-success">You have been added as a subuser for <em>'.$info->name.'</em>!</div>');
			$response->redirect('/account')->send();

		} catch(\Exception $e) {

			\Tracy\Debugger::log($e);
			$service->flash('<div class="alert alert-danger">The server management daemon is not responding, we were unable to add your permissions. Please try again later.</div>');
			$response->redirect('/account')->send();

		}

	}

});

$klein->respond('GET', '/index', function($request, $response, $service) use ($core) {

	if($core->user->getData('root_admin') == '1') {

		$servers = ORM::forTable('servers')
			->select('servers.*')->select('nodes.node', 'node_name')->select('locations.long', 'location')
			->join('nodes', array('servers.node', '=', 'nodes.id'))
			->join('locations', array('nodes.location', '=', 'locations.short'))
			->orderByDesc('active')
			->findArray();

	} else {

		$servers = ORM::forTable('servers')
			->select('servers.*')->select('nodes.node', 'node_name')->select('locations.long', 'location')
			->join('nodes', array('servers.node', '=', 'nodes.id'))
			->join('locations', array('nodes.location', '=', 'locations.short'))
			->where(array('servers.owner_id' => $core->user->getData('id'), 'servers.active' => 1))
			->where_raw('servers.owner_id = ? OR servers.hash IN(?)', array($core->user->getData('id'), join(',', $core->user->listServerPermissions())))
			->findArray();

	}

	/*
	* List Servers
	*/
	$response->body($core->twig->render('panel/index.html', array(
		'servers' => $servers,
		'flash' => $service->flashes()
	)))->send();


});

$klein->respond('GET', '/index/[:goto]', function($request, $response, $service) use ($core) {

	if(!$core->server->nodeRedirect($request->param('goto'), $core->user->getData('id'), $core->user->getData('root_admin'))) {

		$service->flash('<div class="alert alert-danger">The requested server or function does not exist, or you do not have permission to access that server or function.</div>');
		$response->redirect('/index')->send();

	} else {

		$response->cookie('pp_server_hash', $request->param('goto'), 0);
		$response->redirect('/node/index')->send();

	}

});

$klein->respond('GET', '/language/[:language]', function($request, $response, $service, $app) use ($core) {

	if(file_exists(SRC_DIR.'lang/'.$request->param('language').'.json')) {

		if($app->isLoggedIn) {

			$account = ORM::forTable('users')->findOne($core->user->getData('id'));
			$account->set(array(
				'language' => $request->param('language')
			));
			$account->save();

		}

		$response->cookie("pp_language", $request->param('language'), time() + 2678400);
		$response->redirect(($request->server()["HTTP_REFERER"]) ? $request->server()["HTTP_REFERER"] : '/servers')->send();

	} else {

		$response->redirect(($request->server()["HTTP_REFERER"]) ? $request->server()["HTTP_REFERER"] : '/servers')->send();

	}

});