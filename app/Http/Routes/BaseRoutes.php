<?php

namespace PufferPanel\Http\Routes;

use Illuminate\Routing\Router;

class BaseRoutes {

	public function map(Router $router) {
		$router->get('/', [ 'as' => 'index', 'uses' => 'Base\IndexController@getIndex' ]);
	}

}
