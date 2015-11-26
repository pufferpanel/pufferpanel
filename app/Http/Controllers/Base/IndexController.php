<?php

namespace PufferPanel\Http\Controllers\Base;

use Auth;
use Debugbar;
use Google2FA;
use PufferPanel\Models\User;
use PufferPanel\Models\Server;

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

    /**
     * Returns TOTP Management Page.
     *
     * @param  \Illuminate\Http\Request $request
     * @return \Illuminate\Contracts\View\View
     */
    public function getAccountTotp(Request $request)
    {
        return view('base.totp');
    }

    /**
     * Generates TOTP Secret and returns popup data for user to verify
     * that they can generate a valid response.
     *
     * @param  \Illuminate\Http\Request $request
     * @return \Illuminate\Contracts\View\View
     */
    public function putAccountTotp(Request $request)
    {

        $totpSecretKey = Google2FA::generateSecretKey();

        $user = User::find(Auth::user()->id);
        $user->totp_secret = $totpSecretKey;
        $user->save();

        return response()->json([
            'qrImage' => Google2FA::getQRCodeGoogleUrl(
                'PufferPanel',
                $user->email,
                $user->totp_secret
            ),
            'secret' => $totpSecretKey
        ]);

    }

    /**
     * Verifies that 2FA token recieved is valid and will work on the account.
     *
     * @param  \Illuminate\Http\Request $request
     * @return string
     */
    public function postAccountTotp(Request $request)
    {

        if (!$request->has('token')) {
            return response('No input \'token\' defined.', 500);
        }

        $user = User::find(Auth::user()->id);
        if (!Google2FA::verifyKey($user->totp_secret, $request->input('token'))) {
            return response('false');
        }

        $user->use_totp = 1;
        $user->save();

        return response('true');

    }

    public function deleteAccountTotp(Request $request)
    {

        if (!$request->has('token')) {
            return redirect()->route('account.totp')->with('flash-error', 'Missing required `token` field in request.');
        }

        $user = User::find(Auth::user()->id);
        if (!Google2FA::verifyKey($user->totp_secret, $request->input('token'))) {
            return response('false');
        }

        $user->totp_secret = null;
        $user->use_totp = 0;
        $user->save();

        return redirect()->route('account.totp');

    }

}
