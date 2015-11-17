<?php

namespace PufferPanel\Http\Controllers\Base;

use PufferPanel\User;
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
        // $request->user() is instance of User::find(id)
        Debugbar::info($request->user()->toJson());
        return view('base.index', [
            'ip' => $request->ip(),
        ]);
    }

}
