<?php

namespace PufferPanel\Http\Controllers\Base;

use Auth;
use Debugbar;
use Google2FA;
use Log;
use PufferPanel\Exceptions\AccountNotFoundException;
use PufferPanel\Exceptions\DisplayException;
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

        $this->validate($request, [
            'new_email' => 'required|email',
            'password' => 'required'
        ]);

        if (!password_verify($request->input('password'), Auth::user()->password)) {
            return redirect()->route('account')->with('flash-error', 'The password provided was not valid for this account.');
        }

        // Met Validation, lets roll out.
        try {
            User::setEmail(Auth::user()->id, $request->input('new_email'));
            return redirect()->route('account')->with('flash-success', 'Your email address has successfully been updated.');
        } catch (\Exception $e) {
            if ($e instanceof AccountNotFoundException || $e instanceof DisplayException) {
                return redirect()->route('account')->with('flash-error', $e->getMessage());
            }
            throw $e;
        }
    }

    /**
     * Update an account password.
     *
     * @param  \Illuminate\Http\Request $request
     * @return \Illuminate\Http\Response
     */
    public function postAccountPassword(Request $request)
    {

        $this->validate($request, [
            'current_password' => 'required',
            'new_password' => 'required|confirmed|different:current_password|regex:((?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,})',
            'new_password_confirmation' => 'required'
        ]);

        if (!password_verify($request->input('current_password'), Auth::user()->password)) {
            return redirect()->route('account')->with('flash-error', 'The password provided was not valid for this account.');
        }

        // Met Validation, lets roll out.
        try {
            User::setPassword(Auth::user()->id, $request->input('new_password'));
            return redirect()->route('account')->with('flash-success', 'Your password has successfully been updated.');
        } catch (\Exception $e) {
            if ($e instanceof AccountNotFoundException || $e instanceof DisplayException) {
                return redirect()->route('account')->with('flash-error', $e->getMessage());
            }
            throw $e;
        }

    }

}
