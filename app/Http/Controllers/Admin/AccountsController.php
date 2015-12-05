<?php

namespace PufferPanel\Http\Controllers\Admin;

use Debugbar;
use PufferPanel\Models\User;

use PufferPanel\Http\Controllers\Controller;
use Illuminate\Http\Request;

class AccountsController extends Controller
{

    /**
     * Controller Constructor
     */
    public function __construct()
    {

        // All routes in this controller are protected by the authentication middleware.
        $this->middleware('auth');
        $this->middleware('admin');

    }

    public function getIndex()
    {
        return view('admin.accounts.index', [
            'users' => User::all()
        ]);
    }

}
