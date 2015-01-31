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

require_once('../../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true)
	Components\Page::redirect('../../../index.php');

setcookie("__TMP_pp_admin_updateglobal", json_encode($_POST), time() + 30, '/');

if(!isset($_POST['main_url'], $_POST['master_url'], $_POST['assets_url']))
	Components\Page::redirect('../urls.php?error=main_url|master_url|assets_url');

foreach($_POST as $id => $val)
	if(!preg_match('/^((https?:){0,1})(\/\/){1}([-\d\w\/.]*)$/', $val))
		Components\Page::redirect('../urls.php?error='.$id);
	else
		$_POST[$id] = preg_replace('/^(http:|https:)(.*)?$/', '$2', $val);

$query = ORM::forTable('acp_settings')->rawExecute("
UPDATE acp_settings SET setting_val = CASE setting_ref
	WHEN 'main_website' THEN :main_url
	WHEN 'master_url' THEN :master_url
	WHEN 'assets_url' THEN :assets_url
	ELSE setting_val
END
", array(
	'main_url' => $_POST['main_url'],
	'master_url' => $_POST['master_url'],
	'assets_url' => $_POST['assets_url']
));

Components\Page::redirect('../urls.php?success=true');

?>