<?php

namespace PufferPanel\Http\Controllers\API;

use Log;
use Debugbar;
use PufferPanel\Models\User;

use PufferPanel\Http\Controllers\Controller;
use Illuminate\Http\Request;

class UserController extends Controller
{

    /**
     * Constructor
     */
    public function __construct()
    {
        // $this->middleware('api');
    }


    public function getAllUsers(Request $request)
    {
        return response()->json(User::all());
    }

    /**
     * Returns JSON response about a user given their ID.
     * If fields are provided only those fields are returned.
     *
     * Does not return protected fields (i.e. password & totp_secret)
     *
     * @param  Request $request
     * @param  int     $id
     * @param  string  $fields
     * @return Response
     */
    public function getUser(Request $request, $id, $fields = null)
    {
        if (is_null($fields)) {
            return response()->json(User::find($id));
        }

        $query = User::where('id', $id);
        $explode = explode(',', $fields);

        foreach($explode as &$exploded) {
            if(!empty($exploded)) {
                $query->addSelect($exploded);
            }
        }

        try {
            return response()->json($query->get());
        } catch (\Exception $e) {
            if ($e instanceof \Illuminate\Database\QueryException) {
                return response()->json([
                    'error' => 'One of the fields provided in your argument list is invalid.'
                ], 500);
            }
            throw $e;
        }

    }

}
