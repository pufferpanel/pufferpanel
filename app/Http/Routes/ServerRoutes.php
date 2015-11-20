<?php

namespace PufferPanel\Http\Routes;

use Illuminate\Routing\Router;

class ServerRoutes {

	public function map(Router $router) {
		$router->group(['prefix' => 'server/{server}'], function ($server) use ($router) {
			$router->get('/', 'Server\ServerController@getIndex');

			// Ajax Routes
			$router->get('/ajax/status', 'Server\AjaxController@getStatus');
		});
	}

}