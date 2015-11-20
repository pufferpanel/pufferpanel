<?php

namespace PufferPanel\Http\Controllers\Server;

use Gate;
use Auth;
use PufferPanel\Models\Server;
use PufferPanel\Models\Node;
use Debugbar;

use PufferPanel\Http\Controllers\Controller;
use Illuminate\Http\Request;

class ServerController extends Controller
{

    /**
     * Controller Constructor
     */
    public function __construct()
    {

        // All routes in this controller are protected by the authentication middleware.
        $this->middleware('auth');

        // Routes in this file are also checked aganist the server middleware. If the user
        // does not have permission to view the server it will not load.
        $this->middleware('server');

    }

    /**
     * Returns server index page for specified server.
     *
     * @param  \Illuminate\Http\Request $request
     * @return \Illuminate\Contracts\View\View
     */
    public function getIndex(Request $request)
    {

        $server = Server::getByUUID($request->route()->server);
        return view('server.index', [
            'server' => $server,
            'node' => Node::find($server->node)
        ]);
    }

}
