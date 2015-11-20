<?php

namespace PufferPanel\Http\Routes;

use Illuminate\Routing\Router;

class ServerRoutes {

	public function map(Router $router) {
		$router->group(['prefix' => 'server/{server}'], function ($server) use ($router) {
			$router->get('/', [ 'as' => 'server.index', 'uses' => 'Server\ServerController@getIndex' ]);

			// Ajax Routes
			$router->get('/ajax/status', [ 'as' => 'server.ajax.status', 'uses' => 'Server\AjaxController@getStatus' ]);
		});
	}

}