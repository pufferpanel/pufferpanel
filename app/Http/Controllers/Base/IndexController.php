<?php

namespace PufferPanel\Http\Controllers\Base;

use PufferPanel\User;
use PufferPanel\Server;
use Debugbar;

use PufferPanel\Http\Controllers\Controller;
use Illuminate\Http\Request;

class IndexController extends Controller
{

    /**
     * Controller Constructor
     */
    public function __construct()
    {

        // All routes in this controller are protected by the authentication middleware.
        $this->middleware('auth');
    }

    /**
     * Returns listing of user's servers.
     *
     * @param  \Illuminate\Http\Request $request
     * @return \Illuminate\Contracts\View\View
     */
    public function getIndex(Request $request)
    {
        return view('base.index', [
            'servers' => Server::getUserServers(),
        ]);
    }

}
