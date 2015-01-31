<?php
/*
    PufferPanel - A Minecraft Server Management Panel
    Copyright (c) 2013 Dane Everitt

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
use \ORM, \Tracy\Debugger;

$klein->respond('GET', '/admin/settings/[:page]', function($request, $response, $service) use ($core) {

    $response->body($core->twig->render(
        'admin/settings/'.$request->param('page').'.html',
        array(
            'flash' => $service->flashes()
        )
    ))->send();

});

$klein->respond('POST', '/admin/settings/[:page]/[:action]', function($request, $response, $service) use ($core) {

    // Set Company Name
    if($request->param('page') == "global" && $request->param('action') == "company") {

        $query = ORM::forTable('acp_settings')->where('setting_ref', 'company_name')->findOne();
        $query->setting_val = $request->param('company_name');
        $query->save();

        $service->flash('<div class="alert alert-success">Your company name has been successfully updated.</div>');
        $response->redirect('/admin/settings/global')->send();

    }

    // Update Global Settings
    if($request->param('page') == "global" && $request->param('action') == "general") {

        try {

            $query = ORM::forTable('acp_settings')->rawExecute("
                UPDATE acp_settings SET setting_val = CASE setting_ref
                    WHEN 'use_api' THEN :enable_api
                    WHEN 'force_online' THEN :force_online
                    WHEN 'https' THEN :https
                    WHEN 'allow_subusers' THEN :allow_subusers
                    ELSE setting_val
                END
                ", array(
                    'enable_api' => (!in_array('use_api', $request->param('permissions'))) ? 0 : 1,
                    'force_online' => (!in_array('force_online', $request->param('permissions'))) ? 0 : 1,
                    'https' => (!in_array('https', $request->param('permissions'))) ? 0 : 1,
                    'allow_subusers' => (!in_array('allow_subusers', $request->param('permissions'))) ? 0 : 1
                )
            );

            $service->flash('<div class="alert alert-success">Your global settings have been successfully updated.</div>');

        } catch(\Exception $e) {

            Debugger::log($e);
            $service->flash('<div class="alert alert-danger">An error occured while trying to perform this MySQL command.</div>');

        }

        $response->redirect('/admin/settings/global')->send();

    }

    if($request->param('page') == "email" && $request->param('action') == "update") {

        $response->cookie("__TMP_pp_admin_updateglobal", json_encode($request->paramsPost()), time() + 30);

        if(!in_array($request->param('smail_method'), array('php', 'postmark', 'mandrill', 'mailgun', 'sendgrid'))) {

            $service->flash('<div class="alert alert-danger">The email method selected was not a valid choice.</div>');
            $response->redirect('/admin/settings/email')->send();
            return;

        }

        if(!filter_var($request->param('sendmail_email'), FILTER_VALIDATE_EMAIL)) {

            $service->flash('<div class="alert alert-danger">The email provided as the sendmail address is not valid.</div>');
            $response->redirect('/admin/settings/email')->send();
            return;

        }

        if($request->param('smail_method') != 'php' && empty($request->param($request->param('smail_method').'_api_key'))) {

            $service->flash('<div class="alert alert-danger">The API key was not provided for the selected method.</div>');
            $response->redirect('/admin/settings/email')->send();
            return;

        }

        /*
         * Handle Sendgrid Information
         */
        $sendgrid = null;
        if(strpos($request->param('sendgrid_api_key'), '|')) {
            $iv = $core->auth->generate_iv();
            $sendgrid = $iv.'.'.$core->auth->encrypt($request->param('sendgrid_api_key'), $iv);
        }

        try {

            $query = ORM::forTable('acp_settings')->rawExecute("
                UPDATE acp_settings SET setting_val = CASE setting_ref
                    WHEN 'sendmail_method' THEN :sendmail_method
                    WHEN 'sendmail_email' THEN :sendmail_email
                    WHEN 'postmark_api_key' THEN :postmark_api_key
                    WHEN 'mandrill_api_key' THEN :mandrill_api_key
                    WHEN 'mailgun_api_key' THEN :mailgun_api_key
                    WHEN 'sendgrid_api_key' THEN :sendgrid_api_key
                    ELSE setting_val
                END
                ", array(
                    'sendmail_method' => $request->param('smail_method'),
                    'sendmail_email' => $request->param('sendmail_email'),
                    'postmark_api_key' => $request->param('postmark_api_key'),
                    'mandrill_api_key' => $request->param('mandrill_api_key'),
                    'mailgun_api_key' => $request->param('mailgun_api_key'),
                    'sendgrid_api_key' => $sendgrid,
                )
            );

            $service->flash('<div class="alert alert-success">Your email settings have been updated.</div>');

        } catch(\Exception $e) {

            Debugger::log($e);
            $service->flash('<div class="alert alert-danger">An error occured while trying to perform this MySQL command.</div>');

        }

        $response->redirect('/admin/settings/email')->send();

    }

});