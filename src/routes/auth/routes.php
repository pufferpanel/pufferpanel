<?php

/*
  PufferPanel - A Game Server Management Panel
  Copyright (c) 2015 Dane Everitt

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see http://www.gnu.org/licenses/.
 */

namespace PufferPanel\Core;

use \ORM,
    \Unirest;

$klein->respond('GET', '/auth/login', function($request, $response, $service) use ($core) {

    $response->body($core->twig->render('auth/login.html', array(
                'xsrf' => $core->auth->XSRF(),
                'flash' => $service->flashes()
    )))->send();
});

$klein->respond('POST', '/auth/login', function($request, $response, $service) use ($core) {

    if (!$core->auth->XSRF($request->param('xsrf'))) {

        $service->flash('<div class="alert alert-warning"> The XSRF token recieved was not valid. Please make sure cookies are enabled and try your request again.</div>');
        $response->redirect('/auth/login')->send();
        return;
    }

    $account = ORM::forTable('users')->where('email', $request->param('email'))->findOne();

    if (!$core->auth->verifyPassword($request->param('email'), $request->param('password'))) {

        if ($account && $account->notify_login_f == 1) {

            $core->email->generateLoginNotification('failed', array(
                'IP_ADDRESS' => $request->ip(),
                'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($request->ip())
            ))->dispatch($request->param('email'), Settings::config()->company_name . ' - Account Login Failure Notification');
        }

        $service->flash('<div class="alert alert-danger"><strong>Oh snap!</strong> The username or password you submitted was incorrect.</div>');
        $response->redirect('/auth/login')->send();
        return;
    } else {

        if ($account->use_totp == 1 && !$core->auth->validateTOTP($request->param('totp_token'), $account->totp_secret)) {
            $service->flash('<div class="alert alert-danger"><strong>Oh snap!</strong> Your Two-Factor Authentication token was missing or incorrect.</div>');
            $response->redirect('/auth/login')->send();
            return;
        }

        $cookie = (object) array('token' => $core->auth->keygen('12'), 'expires' => ($request->param('remember_me')) ? (time() + 604800) : null);

        $account->set(array(
            'session_id' => $cookie->token,
            'session_ip' => $request->ip()
        ));
        $account->save();

        if ($account->notify_login_s == 1) {

            $core->email->generateLoginNotification('success', array(
                'IP_ADDRESS' => $request->ip(),
                'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($request->ip())
            ))->dispatch($request->param('email'), Settings::config()->company_name . ' - Account Login Notification');
        }

        $response->cookie('pp_auth_token', $cookie->token, $cookie->expires);        
        
        $response->redirect('/index')->send();
        return;
    }
});

$klein->respond('POST', '/auth/login/totp', function($request, $response) {

    if (!$request->param('check')) {
        $response->body('false')->send();
        return;
    } else {

        $totp = ORM::forTable('users')->select('use_totp')->where('email', $request->param('check'))->findOne();

        if (!$totp) {
            $response->body('false')->send();
            return;
        } else {
            $response->body(($totp->use_totp == 1) ? 'true' : 'false')->send();
            return;
        }
    }
});

$klein->respond('GET', '/auth/logout', function($request, $response) use ($core) {

    if ($core->auth->isLoggedIn()) {

        $response->cookie("pp_auth_token", null, time() - 86400);
        $response->cookie("pp_server_node", null, time() - 86400);
        $response->cookie("pp_server_hash", null, time() - 86400);
        $response->cookie("accessToken", null, time() - 86400);

        $logout = ORM::forTable('users')->where(array('session_id' => $request->cookies()['pp_auth_token'], 'session_ip' => $request->ip()))->findOne();
        $logout->session_id = null;
        $logout->session_ip = null;
        $logout->save();
    }

    $response->redirect('/auth/login')->send();
});


$klein->respond('GET', '/auth/password', function($request, $response, $service) use ($core) {

    $response->body($core->twig->render('auth/password.html', array(
                'xsrf' => $core->auth->XSRF(),
                'flash' => $service->flashes()
    )))->send();
});

$klein->respond('GET', '/auth/password/[:action]', function($request, $response, $service) use ($core) {

    $response->body($core->twig->render('auth/password.html', array(
                'flash' => $service->flashes(),
                'noshow' => true
    )))->send();
});

$klein->respond('GET', '/auth/password/verify/[:key]', function($request, $response, $service) use ($core) {

    $query = ORM::forTable('account_change')->where(array('key' => $request->param('key'), 'verified' => 0))->where_gt('time', time())->findOne();

    if (!$query) {

        $service->flash('<div class="alert alert-danger">Unable to verify password recovery request.<br />Did the key expire? Please contact support for more help or try again.</div>');
        $response->redirect('/auth/password')->send();
    } else {

        $password = $core->auth->keygen('12');
        $query->verified = 1;

        $user = ORM::forTable('users')->where('email', $query->content)->findOne();
        $user->password = $core->auth->hash($password);
        $user->save();

        $query->save();

        /*
         * Send Email
         */
        $core->email->buildEmail('new_password', array(
            'NEW_PASS' => $password,
            'EMAIL' => $query->content
        ))->dispatch($query->content, Settings::config()->company_name . ' - New Password');

        $service->flash('<div class="alert alert-success">You should recieve an email within the next 5 minutes (usually instantly) with your new account password. We suggest changing this once you log in.</div>');
        $response->redirect('/auth/login')->send();
    }
});

$klein->respond('POST', '/auth/password', function($request, $response, $service) use ($core) {

    if ($core->auth->XSRF($request->param('xsrf')) !== true) {

        $service->flash('<div class="alert alert-warning">The XSRF token recieved was not valid. Please make sure cookies are enabled and try your request again.</div>');
        $response->redirect('/auth/password')->send();
    }

    try {

        $unirest = Unirest\Request::get(
                        "https://www.google.com/recaptcha/api/siteverify", array(), array(
                    'secret' => Settings::config('captcha_priv'),
                    'response' => $request->param('g-recaptcha-response'),
                    'remoteip' => $request->ip()
                        )
        );

        if (!isset($unirest->body->success) || !$unirest->body->success) {

            $service->flash('<div class="alert alert-danger">The spam prevention was not filled out correctly. Please try it again.</div>');
            $response->redirect('/auth/password')->send();
            return;
        }
    } catch (\Exception $e) {

        $service->flash('<div class="alert alert-danger">Unable to query the captcha validation servers. Please try it again.</div>');
        $response->redirect('/auth/password')->send();
        return;
    }

    $query = ORM::forTable('users')->where('email', $request->param('email'))->findOne();

    if ($query) {

        $key = $core->auth->keygen('30');

        ORM::forTable('account_change')->create()->set(array(
            'type' => 'password',
            'content' => $request->param('email'),
            'key' => $key,
            'time' => time() + 14400
        ))->save();

        $core->email->buildEmail('password_reset', array(
            'IP_ADDRESS' => $request->ip(),
            'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($request->ip()),
            'PKEY' => $key
        ))->dispatch($request->param('email'), Settings::config()->company_name . ' - Reset Your Password');

        $service->flash('<div class="alert alert-success">We have sent an email to the address you provided in the previous step. Please follow the instructions included in that email to continue. The verification key will expire in 4 hours.</div>');
        $response->redirect('/auth/password/pending')->send();
    } else {

        $service->flash('<div class="alert alert-danger">We couldn\'t find that email in our database.</div>');
        $response->redirect('/auth/password')->send();
    }
});

$klein->respond('GET', '/auth/register/[:token]?', function($request, $response, $service) use ($core) {

    $response->body($core->twig->render('auth/register.html', array(
                'xsrf' => $core->auth->XSRF(),
                'token' => $request->param('token'),
                'flash' => $service->flashes()
    )))->send();
});

$klein->respond('POST', '/auth/register', function($request, $response, $service) use ($core) {

    if (!$request->param('token')) {

        $service->flash('<div class="alert alert-danger"><i class="fa fa-warning"></i> No token was submitted with the request.</div>');
        $response->redirect('/auth/register')->send();
        return;
    }

    /* XSRF Check */
    if (!$core->auth->XSRF($request->param('xsrf'))) {

        $service->flash('<div class="alert alert-warning"><i class="fa fa-warning"></i> Invalid XSRF token submitted with the form.</div>');
        $response->redirect('/auth/register/' . $request->param('token'))->send();
        return;
    }

    $query = ORM::forTable('account_change')->where(array(
                'type' => 'user_register',
                'key' => $request->param('token'),
                'verified' => 0
            ))->findOne();

    if (!$query) {

        $service->flash('<div class="alert alert-danger"><i class="fa fa-warning"></i> The token you provided appears to be invalid.</div>');
        $response->redirect('/auth/register/' . $request->param('token'))->send();
        return;
    }

    if (!preg_match('/^[\w-]{4,35}$/', $request->param('username'))) {

        $service->flash('<div class="alert alert-danger"><i class="fa fa-warning"></i> The username you entered does not meet the requirements. Must be at least 4 characters, and no more than 35. Username can only contain the following characters: a-zA-Z0-9_-</div>');
        $response->redirect('/auth/register/' . $request->param('token'))->send();
        return;
    }

    if (!$core->auth->validatePasswordRequirements($request->param('password'))) {

        $service->flash('<div class="alert alert-danger">Your password is not complex enough. Please make sure to include at least one number, and some type of mixed case. Your new password must also be at least 8 characters long.</div>');
        $response->redirect('/auth/register/' . $request->param('token'))->send();
        return;
    }

    $user = ORM::forTable('users')->where_any_is(array(array('username' => $request->param('username')), array('email' => $request->param('email'))))->findOne();

    if ($user) {

        $service->flash('<div class="alert alert-danger">Account with that username or email already exists in the system.</div>');
        $response->redirect('/auth/register/' . $request->param('token'))->send();
        return;
    }

    ORM::forTable('users')->where('id', $query->user_id)->findOne()->set(array(
        'username' => $request->param('username'),
        'password' => $core->auth->hash($request->param('password')),
        'register_time' => time()
    ))->save();

    $query->delete();

    $service->flash('<div class="alert alert-success">Your account has been created successfully, you may now login.</div>');
    $response->redirect('/auth/login')->send();
});

$klein->respond('GET', '/auth/remote/deploy/[:key]', function($request, $response) use ($core) {

    $deploy = ORM::forTable('autodeploy')->where('code', $request->param('key'))->where_gt('expires', time())->findOne();
    if (!$deploy) {
        $response->code(404)->body('echo "The requested URL does not exist."; exit 1')->send();
        return;
    }

    $node = ORM::forTable('nodes')->findOne($deploy->node);
    if (!$node) {
        $response->code(404)->body('404 Node')->send();
        return;
    }
    // Check Access IP here

    $response->header('Content-Type', 'text/plain');
    $response->body($core->twig->render('templates/auto-deploy.tpl', array(
                'node' => $node,
                'majorVersion' => Version::getMajor()
    )));
    $response->send();
    return;
});
