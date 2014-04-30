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
require_once('../../../../core/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) === true){

	if(isset($_POST['ip'])){
		
		/*
		 * Query Servers
		 */
		if(!@fsockopen($_POST['ip'], 8003, $num, $error, 3))
			exit('<span class="label label-danger">Offline</span>');
		else
			exit('<span class="label label-success">Online</span>');
		
	}else{
	
		exit('<span class="label label-danger">No IP#</span>');
	
	}

}else{

	exit('<span class="label label-danger">Auth Error</span>');

}

?>