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
	$core->framework->page->redirect('../../../../index.php');
}

if(!isset($_GET['node']) || empty($_GET['node']))
	exit('No node defined!');


$selectData = $mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :name");
$selectData->execute(array(
	':name' => $_GET['node']
));
if($selectData->rowCount() == 1)
	$node = $selectData->fetch();
else
	exit('Unknown node.');

$ips = json_decode($node['ips'], true);
$ports = json_decode($node['ports'], true);

$i = 0;
foreach($ips as $ip => $internal){

	if($internal['ports_free'] > 0){
	
		if($i != 0) echo '<br /><br />';
		echo $ip.' has '.$internal['ports_free'].' avaliable port(s).<br />';
		
		foreach($ports[$ip] as $port => $avaliable){
			
			if($avaliable == 1)
				echo ' - - '.$port.' is <span style="color:#52964f;">avaliable</span>.<br />';
			else
				echo ' - - '.$port.' is <span style="color:#cf4425;">unavaliable</span>.<br />';
					
		}
		
	}else{
		
		if($i != 0) echo '<br /><br />';
		echo $ip.' has no avaliable port(s).<br />';
		
	}

	$i++;

}


?>