<?php

namespace PufferPanel\Http\Controllers\Base;

use Auth;
use Debugbar;
use Google2FA;
use Log;
use PufferPanel\Exceptions\AccountNotFoundException;
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

        try {
            $totpSecret = User::setTotpSecret(Auth::user()->id);
        } catch (\Exception $e) {
            if ($e instanceof AccountNotFoundException) {
                return response($e->getMessage(), 500);
            }
            throw $e;
        }

        return response()->json([
            'qrImage' => Google2FA::getQRCodeGoogleUrl(
                'PufferPanel',
                Auth::user()->email,
                $totpSecret
            ),
            'secret' => $totpSecret
        ]);

    }

    /**
     * Verifies that 2FA token recieved is valid and will work on the account.
     *
     * @param  \Illuminate\Http\Request $request
     * @return \Illuminate\Http\Response
     */
    public function postAccountTotp(Request $request)
    {

        if (!$request->has('token')) {
            return response('No input \'token\' defined.', 500);
        }

        try {
            if(User::toggleTotp(Auth::user()->id, $request->input('token'))) {
                return response('true');
            }
            return response('false');
        } catch (\Exception $e) {
            if ($e instanceof AccountNotFoundException) {
                return response($e->getMessage(), 500);
            }
            throw $e;
        }

    }

    /**
     * Disables TOTP on an account.
     *
     * @param  \Illuminate\Http\Request $request
     * @return \Illuminate\Http\Response
     */
    public function deleteAccountTotp(Request $request)
    {

        if (!$request->has('token')) {
            return redirect()->route('account.totp')->with('flash-error', 'Missing required `token` field in request.');
        }

        try {
            if(User::toggleTotp(Auth::user()->id, $request->input('token'))) {
                return redirect()->route('account.totp');
            }
            return redirect()->route('account.totp')->with('flash-error', 'Unable to disable TOTP on this account, was the token correct?');
        } catch (\Exception $e) {
            if ($e instanceof AccountNotFoundException) {
                return redirect()->route('account.totp')->with('flash-error', 'An error occured while attempting to perform this action.');
            }
            throw $e;
        }

    }

    /**
     * Display base account information page.
     *
     * @param  \Illuminate\Http\Request $request
     * @return \Illuminate\Contracts\View\View
     */
    public function getAccount(Request $request)
    {
        return view('base.account');
    }

    /**
     * Update an account email.
     *
     * @param  \Illuminate\Http\Request $request
     * @return \Illuminate\Http\Response
     */
    public function postAccountEmail(Request $request)
    {
        //
    }

    /**
     * Update an account password.
     *
     * @param  \Illuminate\Http\Request $request
     * @return \Illuminate\Http\Response
     */
    public function postAccountPassword(Request $request)
    {
        //
    }

}
