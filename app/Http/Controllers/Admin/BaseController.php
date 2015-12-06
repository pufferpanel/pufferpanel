<?php

namespace PufferPanel\Http\Controllers\Admin;

use Debugbar;

use PufferPanel\Http\Controllers\Controller;
use Illuminate\Http\Request;

class BaseController extends Controller
{

    public function getIndex(Request $request)
    {
        return view('admin.index');
    }

}
