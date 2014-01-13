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
class page extends dbConn {

	public function __construct($user, $settings)
		{
		
			$this->mysql = parent::getConnection();
			$this->user = $user;
			$this->settings = $settings;
		
		}

	public function redirect($url, $redirect = '') {
        
        $doRedirect = (isset($redirect) && !empty($redirect)) ? '?redirect='.$redirect : '';
        if(!headers_sent()){
			header('Location: '.urldecode($url).$doRedirect);
			exit;
		}else{
			exit('<meta http-equiv="refresh" content="0;url='.urldecode($url).$doRedirect.'"/>');
			return;
		}
		
	}
    
    public function genRedirect()
        {
        
            $https = (isset($_SERVER['HTTPS'])) ? 'https://' : 'http://';
            return $https.$_SERVER['HTTP_HOST'].$_SERVER['REQUEST_URI'];
        
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
	
	public function override_getCount(){
		return $this->getCount();
	}
	
	public function nodeRedirect($hash)
		{
		
			if($this->user->getData('root_admin') == '1'){
				$query = $this->mysql->prepare("SELECT * FROM `servers` WHERE `hash` = ? AND `active` = '1'");
				$query->execute(array($hash));
			}else{
				$query = $this->mysql->prepare("SELECT * FROM `servers` WHERE `owner_id` = :ownerid AND `hash` = :hash AND `active` = '1'");
				$query->execute(array(
					':ownerid' => $this->user->getData('id'),
					':hash' => $hash
				));
			}
			
				if($query->rowCount() == 1){
				
					$row = $query->fetch();
					
						setcookie('pp_server_hash', $row['hash'], 0, '/', $this->settings->get('cookie_website'));
					
						/*
						 * Select Node Information
						 */
						$queryNode = $this->mysql->prepare("SELECT * FROM `nodes` WHERE `id` = ? LIMIT 1");
						$queryNode->execute(array($row['node']));
						
							if($queryNode->rowCount() == 1){
						 
						 		$queryRow = $queryNode->fetch();
								$this->redirect($queryRow['node_link'].'index.php');
								
							}else{
							
								$this->redirect('servers.php?error=error&c=NODE_NOT_FOUND');
							
							}
				
				}else{
				
					$this->redirect('servers.php?error=error');
				
				}
		
		}
	
}

?>