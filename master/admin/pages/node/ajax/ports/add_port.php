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
session_start();
require_once('../../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	exit('<div class="error-box round">Failed to Authenticate Account.</div>');
}

if(!isset($_POST['add_ports_node']))
	$core->framework->page->redirect('../../list.php');

if(!isset($_POST['add_ports'], $_POST['add_ports_ip']))
	$core->framework->page->redirect('../../view.php?id='.$_POST['add_ports_node']);
	
if(!preg_match('/^[\d, ]+$/', $_POST['add_ports']))
	$core->framework->page->redirect('../../view.php?id='.$_POST['add_ports_node']);
	
$ports = explode(',', str_replace(" ", "", $_POST['add_ports']));

$select = $mysql->prepare("SELECT `ips`, `ports` FROM `nodes` WHERE `id` = :nid");
$select->execute(array(':nid' => $_POST['add_ports_node']));
$data = $select->fetch();

$saveips = json_decode($data['ips'], true);
$saveports = json_decode($data['ports'], true);

foreach($ports as $id => $port)
	{

		if(strlen($port) < 6 && strlen($port) > 0 && array_key_exists($_POST['add_ports_ip'], $saveports) && !array_key_exists($port, $saveports[$_POST['add_ports_ip']])){
		
			$saveports[$_POST['add_ports_ip']][$port] = 1;
			$saveips[$_POST['add_ports_ip']]['ports_free']++;
			
		}

	}
	
$mysql->prepare("UPDATE `nodes` SET `ips` = :ips, `ports` = :ports WHERE `id` = :nid")->execute(array(
	':nid' => $_POST['add_ports_node'],
	':ips' => json_encode($saveips),
	':ports' => json_encode($saveports)
));
$core->framework->page->redirect('../../view.php?id='.$_POST['add_ports_node']);

?>