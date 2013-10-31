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
 
/*
 * PufferPanel Page Actions Function File
 */
 
class page {

	public function redirect($url) {
		
		if(!headers_sent()){
			header('Location: '.urldecode($url));
			exit;
		}else{
			exit('<meta http-equiv="refresh" content="0;url='.urldecode($url).'"/>');
			return;
		}
		
	}
	
	public function isActive($p, $s){
	
		if($p == 'i' && $s == 'acp_index')
			return 'active';
		else if($p == 'c' && $s == 'clients')
			return 'active';
		else if($p == 's' && $s == 'servers')
			return 'active';
		else if($p =='b' && $s == 'backups')
			return 'active';
		else if($p == 'n' && $s == 'nodes')
			return 'active';
		else if($p == 'sett' && $s == 'settings')
			return 'active';
	
	}
	
}

?>