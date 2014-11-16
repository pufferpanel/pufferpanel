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

$ports = [];

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Components\Page::redirect('../../../index.php');
}

if(!isset($_POST['add_ports_node']))
	Components\Page::redirect('../../list.php');

if(!isset($_POST['add_ports'], $_POST['add_ports_ip']))
	Components\Page::redirect('../../view.php?id='.$_POST['add_ports_node'].'&tab=allocation');

if(!preg_match('/^[\d-,]+$/', $_POST['add_ports']))
	Components\Page::redirect('../../view.php?id='.$_POST['add_ports_node'].'&disp=add_port_fail&tab=allocation');

if(preg_match('/^[\d-,]+$/', $_POST['add_ports']))
{
	$portsSplit = explode(',', $_POST['add_ports']);
	for($i = 0; $i < count($portsSplit); $i++)
	{
		if(preg_match('/^[\d-]+$/', $portsSplit[$i]) && !preg_match('/^[\d ]+$/',$portsSplit[$i]))
		{
			$a = explode('-', $portsSplit[$i]);
			$rangeOne = intval($a[0]);
			$rangeTwo = intval($a[1]);
			while($rangeOne <= $rangeTwo)
			{
				$ports[] = $rangeOne;
				$rangeOne++;
			}
		}
		else
		{
			$ports[] = $portsSplit[$i];
		}
	}
}
elseif(preg_match('/^[\d-]+$/', $_POST['add_ports']))
{
	$portsSplit = explode('-', str_replace(" ", "", $_POST['add_ports']));
	$rangeOne = intval($portsSplit[0]);
	$rangeTwo = intval($portsSplit[1]);
	while($rangeOne <= $rangeTwo)
	{
		$ports[] = $rangeOne;
		$rangeOne++;
	}
} else {
	$ports = explode(',', str_replace(" ", "", $_POST['add_ports']));
}


$node = ORM::forTable('nodes')->findOne($_POST['add_ports_node']);

$saveips = json_decode($node->ips, true);
$saveports = json_decode($node->ports, true);

foreach($ports as $id => $port)
	{

		if(strlen($port) < 6 && strlen($port) > 0 && array_key_exists($_POST['add_ports_ip'], $saveports) && !array_key_exists($port, $saveports[$_POST['add_ports_ip']])){

			$saveports[$_POST['add_ports_ip']][$port] = 1;
			$saveips[$_POST['add_ports_ip']]['ports_free']++;

		}

	}

$node->ips = json_encode($saveips);
$node->ports = json_encode($saveports);
$node->save();

Components\Page::redirect('../../view.php?id='.$_POST['add_ports_node'].'&tab=allocation');
?>
