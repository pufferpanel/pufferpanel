<?php

namespace PufferPanel\Http\Routes;

use Illuminate\Routing\Router;

class BaseRoutes {

    public function map(Router $router) {
        $router->get('/', [ 'as' => 'index', 'uses' => 'Base\IndexController@getIndex' ]);

        $router->get('/account/totp', [ 'as' => 'account.totp', 'uses' => 'Base\IndexController@getAccountTotp' ]);
        $router->put('/account/totp', [ 'uses' => 'Base\IndexController@putAccountTotp' ]);
        $router->post('/account/totp', [ 'uses' => 'Base\IndexController@postAccountTotp' ]);
        $router->delete('/account/totp', [ 'uses' => 'Base\IndexController@deleteAccountTotp' ]);

    }

}
