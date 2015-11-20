<?php

namespace PufferPanel\Http\Routes;

use Illuminate\Routing\Router;

class BaseRoutes {

	public function map(Router $router) {
		$router->get('/', 'Base\IndexController@getIndex');

		$router->get('login', 'Auth\AuthController@getLogin');
		$router->post('login', 'Auth\AuthController@postLogin');
	}

}