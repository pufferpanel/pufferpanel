<?php

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