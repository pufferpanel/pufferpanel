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
	Components\Page::redirect('../../../index.php?login');
}

if(!isset($_POST['nid']) || !is_numeric($_POST['nid']))
	Components\Page::redirect('../../list.php');

if(!isset($_POST['name'], $_POST['ip']))
	Components\Page::redirect('../../view.php?id='.$_POST['nid'].'&disp=missing_args');

/*
 * Validate Node Name
 */
if(!preg_match('/^[\w.-]{1,15}$/', $_POST['name']))
	Components\Page::redirect('../../view.php?id='.$_POST['nid'].'&error=name&disp=n_fail');

/*
 * Validate IP
 */
if(!filter_var(gethostbyname($_POST['fqdn']), FILTER_VALIDATE_IP, FILTER_FLAG_NO_RES_RANGE))
	Components\Page::redirect('../../view.php?id='.$_POST['nid'].'&error=fqdn&disp=ip_fail');

/*
 * Update Record
 */
$mysql->prepare("UPDATE `nodes` SET `node` = :name, `fqdn` = :fqdn WHERE `id` = :nid")->execute(array(
	':nid' => $_POST['nid'],
	':name' => $_POST['name'],
	':fqdn' => $_POST['fqdn']
));

Components\Page::redirect('../../view.php?id='.$_POST['nid']);

?>
