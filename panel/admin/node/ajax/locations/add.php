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
use \ORM as ORM;

require_once('../../../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Components\Page::redirect('../../../index.php');
}

setcookie("__TMP_pp_admin_newlocation", json_encode($_POST), time() + 30, '/');

if(!isset($_POST['shortcode'], $_POST['location'])) {
	Components\Page::redirect('../../locations.php?error=missing_args');
}

if(!preg_match('/^[\w-]{1,10}$/', $_POST['shortcode'])) {
	Components\Page::redirect('../../locations.php?error=shortcode');
}

if(empty($_POST['location'])) {
	Components\Page::redirect('../../locations.php?error=location');
}

$location = ORM::forTable('locations')->create();
$location->short = $_POST['shortcode'];
$location->long = $_POST['location'];
$location->save();

Components\Page::redirect('../../locations.php?success');