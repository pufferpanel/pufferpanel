<?php

namespace PufferPanel\Http\Routes;

use Illuminate\Routing\Router;

class AuthRoutes {

	public function map(Router $router) {
		$router->group(['prefix' => 'auth'], function () use ($router) {
			$router->get('logout', 'Auth\AuthController@getLogout');
		});
	}

}