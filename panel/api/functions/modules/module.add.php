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

namespace Modules\Add;

trait server {

	use \Modules\Functions;

	private function validateData() {
	
		/*
		 * @TODO: Validate Incoming Data
		 */
		$this->data = self::getStoredData()['data'];
		
		/*
		 * Validate Server Name
		 */
		if(!preg_match('/^[\w-]{4,35}$/', $this->data['server_name']))
			self::throwResponse('Invalid server name was provided.', false);
			
		/*
		 * Validate Email
		 */
		if(!filter_var($_POST['email'], FILTER_VALIDATE_EMAIL))
			self::throwResponse('Invalid email was provided.', false);
		
		/*
		 * Validate Disk & Memory
		 */	
		if(!is_numeric($_POST['alloc_mem']) || !is_numeric($_POST['alloc_disk']))
			self::throwResponse('Invalid memory or disk allocation.', false);
		
		/*
		 * Validate CPU Limit
		 */	
		if(!is_numeric($_POST['cpu_limit']))
			self::throwResponse('Invalid CPU limit.', false);
			
		/*
		 * Validate User Exists
		 */
		$this->validateUserbyEmail();
		
	}
	
	private function runRequest() {
	
		$this->data = self::getStoredData()['data'];
		
		/*
		 * Generate Variables
		 */
		$this->ftpUser = self::generateFTPUsername($this->data['server_name']);
		$this->iv = $this->generate_iv();
		$this->password = $this->encrypt($this->data['sftp_pass'], $this->iv);
		
		/*
		 * Add Server to Database
		 */	
		$this->add = $this->mysql->prepare("INSERT INTO `servers` VALUES(NULL, NULL, NULL, :hash, :gsd_secret, :e_iv, :node, :sname, :modpack, :sjar, 1, :oid, :ram, :disk, :cpu, :date, :sip, :sport, :ftpuser, :ftppass)");
		
		$this->add->execute(array(
			':hash' => self::keygen(42),
			':gsd_secret' => self::keygen(16).self::keygen(16),
			':e_iv' => $this->iv,
			':node' => $this->data['node'],
			':sname' => $this->data['server_name'],
			':modpack' => $this->data['modpack'],
			':sjar' => "server",
			':oid' => $this->oid,
			':ram' => $this->data['alloc_mem'],
			':disk' => $this->data['alloc_disk'],
			':cpu' => $this->data['cpu_limit'],
			':date' => time(),
			':sip' => $this->data['server_ip'],
			':sport' => $this->data['server_port'],
			':ftpuser' => $this->ftpUser,
			':ftppass' => $this->password
		));
		
		$this->lastInsert = $this->mysql->lastInsertId();	
	
	}
	
	private function addToGSD() {
	
		/*
		 * Add User to GSD
		 */
	
	}
	
	private function validateUserbyEmail() {
	
		$this->email = self::getStoredData()['data']['email'];
		
		$this->select = $this->mysql->prepare("SELECT `id` FROM `users` WHERE `email` = :email");
		$this->select->execute(array(
			':email' => $this->email
		));
		
			if($this->select->rowCount() != 1)
				self::throwResponse('No user is associated with that email.', false);
			else
				$this->oid = $this->select->fetch()['id'];
	
	}
	
	private function handlePorts() {
	
	
	
	}
	
	private function handleModpack() {
	
	
	
	}
	
	private function updatePorts() {
	
		/*
		 * Update IP Count
		 */
		#$ips[$this->data['server_ip']]['ports_free']--;
		#$ports[$this->data['server_ip']][$this->data['server_port']]--;
		
		#$this->mysql->prepare("UPDATE `nodes` SET `ips` = :ips")->execute(array(':ips' => json_encode($ips)));
		#$this->mysql->prepare("UPDATE `nodes` SET `ports` = :ports")->execute(array(':ports' => json_encode($ports)));	
	
	}

}

class apiModuleAddServer extends \query {
	
	use server;
	
	public function __construct() {
	
		$this->mysql = parent::connect();
	
		$this->validateData();
		$this->runRequest();
		#$this->runSSH();
	
	}
	
	private function runSSH() {
	
		//apiModuleAddServer_Extended::run();
	
	}

}

class apiModuleAddServer_Extended extends \ssh {

	use server;
	
	public function run() {
	
		/*
		 * Run SSH Function to Add Server
		 */
	
	}

}

?>