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

	// Update Captcha
	if($request->param('page') == "captcha" && $request->param('action') == "update") {

		ORM::forTable('acp_settings')->rawExecute(
			"UPDATE acp_settings SET setting_val = IF(setting_ref='captcha_pub', :pub, :priv) WHERE setting_ref IN ('captcha_pub', 'captcha_priv')",
			array(
				'pub' => $request->param('pub_key'),
				'priv' => $request->param('priv_key')
			)
		);

		$service->flash('<div class="alert alert-success">Your reCAPTCHA settings have been updated.</div>');
		$response->redirect('/admin/settings/captcha')->send();

	}

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

		$permissionsParam = (is_array($request->param('permissions'))) ? $request->param('permissions') : array();

		try {

			ORM::forTable('acp_settings')->rawExecute(
				"UPDATE acp_settings SET setting_val = CASE setting_ref
                    WHEN 'use_api' THEN :enable_api
                    WHEN 'https' THEN :https
                    WHEN 'allow_subusers' THEN :allow_subusers
					WHEN 'master_url' THEN :master_url
					WHEN 'assets_url' THEN :assets_url
                    ELSE setting_val
                END", array(
					'enable_api' => (!in_array('use_api', $permissionsParam)) ? 0 : 1,
					'https' => (!in_array('https', $permissionsParam)) ? 0 : 1,
					'allow_subusers' => (!in_array('allow_subusers', $permissionsParam)) ? 0 : 1,
					'master_url' => (in_array('https', $permissionsParam)) ? str_replace("http://", "https://", Settings::config()->master_url) : str_replace("https://", "http://", Settings::config()->master_url),
					'assets_url' => (in_array('https', $permissionsParam)) ? str_replace("http://", "https://", Settings::config()->assets_url) : str_replace("https://", "http://", Settings::config()->assets_url)
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

		if(!in_array($request->param('transport_method'), array('php', 'postmark', 'mandrill', 'mailgun', 'sendgrid'))) {

			$service->flash('<div class="alert alert-danger">The email method selected was not a valid choice.</div>');
			$response->redirect('/admin/settings/email')->send();
			return;

		}

		if(!filter_var($request->param('transport_email'), FILTER_VALIDATE_EMAIL)) {

			$service->flash('<div class="alert alert-danger">The email provided as the sendmail address is not valid.</div>');
			$response->redirect('/admin/settings/email')->send();
			return;

		}

		if(empty($request->param('transport_token')) && ($request->param('transport_method') != 'php')) {

			$service->flash('<div class="alert alert-danger">The API key was not provided for the selected method.</div>');
			$response->redirect('/admin/settings/email')->send();
			return;

		}

		try {

			ORM::forTable('acp_settings')->rawExecute(
			"UPDATE acp_settings SET setting_val = CASE setting_ref
                    WHEN 'transport_method' THEN :transport_method
                    WHEN 'transport_email' THEN :transport_email
                    WHEN 'transport_token' THEN :transport_token
                    ELSE setting_val
                END", array(
					'transport_method' => $request->param('transport_method'),
					'transport_email' => $request->param('transport_email'),
					'transport_token' => $request->param('transport_token')
				)
			);

			$service->flash('<div class="alert alert-success">Your email settings have been updated.</div>');

		} catch(\Exception $e) {

			Debugger::log($e);
			$service->flash('<div class="alert alert-danger">An error occured while trying to perform this MySQL command.</div>');

		}

		$response->redirect('/admin/settings/email')->send();

	}

	// Update URLs
	if($request->param('page') == "urls" && $request->param('action') == "update") {

		$urls = array();
		foreach(array(
			'main_url' => $request->param('main_url'),
			'master_url' => $request->param('master_url')
		) as $id => $val) {

			$url = parse_url($val);

			if(!isset($url['host'])) {

				$service->flash('<div class="alert alert-danger">At least one of the URLs provided was invalid and could not be processed.</div>');
				$response->redirect('/admin/settings/urls')->send();
				return;

			}

			$url['path'] = (isset($url['path'])) ? $url['path'] : null;
			$url['port'] = (isset($url['port'])) ? ':'.$url['port'] : null;

			$urls[$id] = (Settings::config()->https == 1) ? 'https://'.$url['host'].$url['port'].$url['path'] : 'http://'.$url['host'].$url['port'].$url['path'];
			$urls[$id] = rtrim($urls[$id], '/').'/';

		}

		try {

			ORM::forTable('acp_settings')->rawExecute(
				"UPDATE acp_settings SET setting_val = CASE setting_ref
                    WHEN 'main_website' THEN :main_url
                    WHEN 'master_url' THEN :master_url
                    WHEN 'assets_url' THEN :assets_url
                    ELSE setting_val
                END", array(
					'main_url' => $urls['main_url'],
					'master_url' => $urls['master_url'],
					'assets_url' => $urls['assets_url']
				)
			);

			$service->flash('<div class="alert alert-success">Your URL settings have been updated.</div>');

		} catch(\Exception $e) {

			Debugger::log($e);
			$service->flash('<div class="alert alert-danger">An error occured while trying to perform this MySQL command.</div>');

		}

		$response->redirect('/admin/settings/urls')->send();

	}

});
