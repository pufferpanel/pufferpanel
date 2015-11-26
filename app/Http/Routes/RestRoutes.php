<?php

namespace PufferPanel\Http\Routes;

use Illuminate\Routing\Router;

class RestRoutes {

    public function map(Router $router) {
        $router->group(['prefix' => 'api'], function ($server) use ($router) {

            $router->group(['prefix' => 'users'], function ($server) use ($router) {

                $router->get('/', [ 'uses' => 'API\UserController@getAllUsers' ]);
                $router->get('/{id}/{fields?}', [ 'uses' => 'API\UserController@getUser' ])->where('id', '[0-9]+');

            });

        });
    }

}
