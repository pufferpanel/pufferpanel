<?php

namespace PufferPanel\Http\Routes;

use Illuminate\Routing\Router;

class AdminRoutes {

    public function map(Router $router) {
        $router->group(['prefix' => 'admin'], function ($server) use ($router) {
            $router->get('/', [ 'as' => 'admin.index', 'uses' => 'Admin\BaseController@getIndex' ]);

            // Account Routes
            $router->group(['prefix' => 'accounts'], function ($server) use ($router) {
                $router->get('/', [ 'as' => 'admin.accounts', 'uses' => 'Admin\AccountsController@getIndex' ]);
                $router->get('/new', [ 'as' => 'admin.accounts.new', 'uses' => 'Admin\AccountsController@getNew' ]);
            });
        });
    }

}
