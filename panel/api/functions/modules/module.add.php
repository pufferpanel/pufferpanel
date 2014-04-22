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

/* 	auth => key,
	function => add,
	data =>	(
		
	)
			
*/

trait addServer {

	use globalInit;
	
	public function validateRequest() {
	
		$data = $this->getStoredData();
		if(array_key_exists('function', $data) && $data['function'] == 'add'){
					
			/*
			 * Is all of the data here?
			 */
			$dataOptions = array(
				'server_name',
				'node',
				'modpack',
				'email',
				'server_ip',
				'server_port',
				'alloc_mem',
				'alloc_disk',
				'sftp_pass',
				'sftp_pass_2',
				'cpu_limit'
			); 
			
			foreach($dataOptions as $dataOption) {
			
				if(!array_key_exists($dataOption, $data['data']) || $data['data'][$dataOption] == "")
					$this->throwResponse('Missing required data values in API call.', false);
			
			}
		
		}else{
		
			$this->throwResponse('Accessing API in an illegal manner.', false);
		
		}
		
		$this->validateData();
	
	}
	
	private function validateData() {
	
		$data = $this->getStoredData();
		$data = $data['data'];
		
		/*
		 * Validate Server Name
		 */
		if(!preg_match('/^[\w-]{4,35}$/', $data['server_name']))
			$this->throwResponse('Function Failed to Finish: invalid server name provided.', false);
			
		/*
		 * Determine if Node (IP & Port) is Avaliable
		 */
		$select = $mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :id");
		$select->execute(array(
			':id' => $data['node']
		));
		
		if($select->rowCount() == 1)
			$node = $select->fetch();
		else
			$this->throwResponse('Function Failed to Finish: that node is not valid.', false);
		
			/*
			 * Validate IP & Port
			 */
			$ips = json_decode($node['ips'], true);
			$ports = json_decode($node['ports'], true);
		
			if(!array_key_exists($data['server_ip'], $ips))
				$this->throwResponse('Function Failed to Finish: that ip does not exist.', false);
				
			if(!array_key_exists($data['server_port'], $ports[$data['server_ip']]))
				$this->throwResponse('Function Failed to Finish: that port does not exist.', false);
				
			if($ports[$data['server_ip']][$data['server_port']] == 0)
				$this->throwResponse('Function Failed to Finish: that port is already in use.', false);
			
		/*
		 * Validate Email
		 */	
		if(!filter_var($data['email'], FILTER_VALIDATE_EMAIL))
			$this->throwResponse('Function Failed to Finish: the email provided is invalid.', false);
		
		$selectEmail = $mysql->prepare("SELECT `id` FROM `users` WHERE `email` = ?");
		$selectEmail->execute(array($_POST['email']));
		
			if($selectEmail->rowCount() != 1)
				$this->throwResponse('Function Failed to Finish: no account with that email exists in the system.', false);
			else {
				$oid = $selectEmail->fetch();
				$oid = $oid['id'];
			}
		
		/*
		 * Validate Disk & Memory
		 */	
		if(!is_numeric($data['alloc_mem']) || !is_numeric($data['alloc_disk']))
			$this->throwResponse('Function Failed to Finish: invalid amount of memory or disk space provided.', false);
		
		/*
		 * Validate CPU Limit
		 */	
		if(!is_numeric($data['cpu_limit']))
			$this->throwResponse('Function Failed to Finish: invalid CPU limit provided.', false);
		
		/*
		 * Validate SFTP Password
		 * @TODO: Modify this to send an encrypted password over HTTPS
		 */
		if($data['sftp_pass'] != $data['sftp_pass_2'] || strlen($data['sftp_pass']) < 8)
			$this->throwResponse('Function Failed to Finish: invalid SFTP passwords provided.', false);			
		
		//@TODO: This is a horrendous mess.
		$this->iv = ;
		$data['sftp_pass'] = openssl_encrypt($data['sftp_pass'], 'AES-256-CBC', file_get_contents(HASH), false, base64_decode($iv));
		
		/*
		 * Validate Modpack
		 */
		$selectPack = $mysql->prepare("SELECT `min_ram`, `server_jar` FROM `modpacks` WHERE `hash` = :hash AND `deleted` = 0");
		$selectPack->execute(array(
			':hash' => $data['modpack']
		));
		
			if($selectPack->rowCount() != 1)
				$this->throwResponse('Function Failed to Finish: that modpack does not exist.', false);
			else
				$pack = $selectPack->fetch();
				
		/*
		 * Modpack RAM Requirements
		 */
		if($pack['min_ram'] > $data['alloc_mem'])
			$this->throwResponse('Function Failed to Finish: not enough ram was allocated to use with this modpack.', false);
	
	}
	
	public function runRequest() {
	
		$data = $this->getStoredData();
		$data = $data['data'];
		
		/*
		 * Add Server to Database
		 */
		$this->ftpUser = (strlen($data['server_name']) > 6) ? substr($data['server_name'], 0, 6).'_'.$core->auth->keygen(5) : $data['server_name'].'_'.$core->auth->keygen((11 - strlen($data['server_name'])));
				
		$this->add = $this->mysql->prepare("INSERT INTO `servers` VALUES(NULL, NULL, NULL, :hash, :gsd_secret, :e_iv, :node, :sname, :modpack, :sjar, 1, :oid, :ram, :disk, :cpu, :date, :sip, :sport, :ftpuser, :ftppass)");
		$this->add->execute(array(
			':hash' => $core->auth->keygen(42),
			':gsd_secret' => $core->auth->keygen(16).$core->auth->keygen(16),
			':e_iv' => $iv,
			':node' => $data['node'],
			':sname' => $data['server_name'],
			':modpack' => $data['modpack'],
			':sjar' => $pack['server_jar'],
			':oid' => $oid,
			':ram' => $data['alloc_mem'],
			':disk' => $data['alloc_disk'],
			':cpu' => $data['cpu_limit'],
			':date' => time(),
			':sip' => $data['server_ip'],
			':sport' => $data['server_port'],
			':ftpuser' => $this->ftpUser,
			':ftppass' => $data['sftp_pass']
		));
		
		$this->lastInsert = $this->mysql->lastInsertId();
		
		/*
		 * Update IP Count
		 */
		$ips[$data['server_ip']]['ports_free']--;
		$ports[$data['server_ip']][$data['server_port']]--;
		
		$mysql->prepare("UPDATE `nodes` SET `ips` = :ips")->execute(array(':ips' => json_encode($ips)));
		$mysql->prepare("UPDATE `nodes` SET `ports` = :ports")->execute(array(':ports' => json_encode($ports)));		
	
	}

}

class apiModuleAddServer extends query {
	
	use addServer;
	
	public function __construct() {
	
		$this->mysql = parent::getConnection();
	
		$this->validateRequest();
		$this->runRequest();
		$this->runSSH();
		
		$this->finish();
	
	}
	
	private function runSSH() {
	
		apiModuleAddServer_Extended::run();
	
	}

}

class apiModuleAddServer_Extended extends ssh {

	use addServer;
	
	public function run() {
	
		/*
		 * Do Server Making Stuff Here 
		 */
		
			/*
			 * Set the Soft Limit
			 */
			$softLimit = ($data['alloc_disk'] <= 512) ? 0 : ($_POST['alloc_disk'] - 512);
			
			$core->ssh->generateSSH2Connection($node['id'], true)->executeSSH2Command('cd /srv/scripts; sudo ./create_user.sh '.$ftpUser.' '.$data['sftp_pass_2'].' '.$softLimit.' '.$_POST['alloc_disk'], false);	
	
	}

}

?>