<?php

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
						//Should be irrelevant:
						//setcookie('pp_server_node', $row['node'], 0, '/', $this->get('cookie_website'));	
					
						/*
						 * Select Node Information
						 */
						$queryNode = $this->mysql->prepare("SELECT * FROM `nodes` WHERE `node_name` = ? LIMIT 1");
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