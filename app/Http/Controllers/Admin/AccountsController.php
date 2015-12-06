<?php

namespace PufferPanel\Http\Controllers\Admin;

use Debugbar;
use PufferPanel\Models\User;

use PufferPanel\Http\Controllers\Controller;
use Illuminate\Http\Request;

class AccountsController extends Controller
{

    public function getIndex(Request $request)
    {
        return view('admin.accounts.index', [
            'users' => User::paginate(20)
        ]);
    }

    public function getNew(Request $request)
    {
        //
    }

    public function getView(Request $request, $id)
    {
        //
    }

}
