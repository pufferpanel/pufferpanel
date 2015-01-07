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
namespace PufferPanel\Core\Router\Node;
use \ORM, \Unirest;

class Users extends \PufferPanel\Core\Email {

	use \PufferPanel\Core\Components\Authentication, \PufferPanel\Core\Components\GSD, \PufferPanel\Core\Components\Error_Handler;

	/**
	 * @param object
	 */
	protected $server;

	/**
	 * @param object
	 */
	protected $settings;

	/**
	 * Constructor class for \PufferPanel\Core\Router\Node\Users
	 *
	 * @return void
	 */
	public function __construct(\PufferPanel\Core\Server $server) {

		$this->server = $server;
		$this->settings = new \PufferPanel\Core\Config\InMemoryDatabaseConfig('acp_settings', 'setting_ref', 'setting_val');

	}

	/**
	 * Handles adding a new subuser to a server.
	 *
	 * @return bool
	 */
	public function addSubuser(array $data) {

		$registerToken = $this->keygen(32);

		$find = ORM::forTable('users')->where('email', $data['email'])->findOne();
		if(!$find) {

			$newAccount = ORM::forTable('account_change')->create();
			$newAccount->set(array(
				'type' => 'user_register',
				'content' => $data['email'],
				'key' => $registerToken,
				'time' => time()
			));
			$newAccount->save();

		}

		$subuserToken = $this->keygen(32);
		$addToken = ORM::forTable('account_change')->create();
		$addToken->set(array(
			'type' => 'subuser',
			'content' => json_encode(array(
				$this->server->getData('hash') => array(
					'key' => $this->generateUniqueUUID('servers', 'gsd_secret'),
					'perms' => self::_rebuildUserPermissions($data['permissions']),
					'perms_gsd' => self::_buildGSDPermissions($data['permissions'])
				)
			)),
			'key' => $subuserToken,
			'time' => time()
		));
		$addToken->save();

		$subusers = (!empty($this->server->getData('subusers'))) ? json_decode($this->server->getData('subusers'), true) : array();
		$subusers[$data['email']] = $subuserToken;

		$updateServer = ORM::forTable('servers')->findOne($this->server->getData('id'));
		$updateServer->subusers = json_encode($subusers);
		$updateServer->save();

		/*
		* Send Email
		*/
		$this->buildEmail('new_subuser', array(
			'REGISTER_TOKEN' => $registerToken,
			'SUBUSER_TOKEN' => $subuserToken,
			'URLENCODE_TOKEN' => urlencode($registerToken),
			'SERVER' => $this->server->getData('name'),
			'EMAIL' => $data['email']
		))->dispatch($data['email'], $this->settings->config('company_name').' - You\'ve Been Invited to Manage a Server');

		return true;

	}

	/**
	 * Updates information for a subuser.
	 *
	 * @param array $data
	 * @return bool
	 */
	public function modifySubuser(array $data) {

		if(!$this->avaliable($this->server->nodeData('ip'), $this->server->nodeData('gsd_listen'))) {
			self::_setError("Unable to access the server management daemon.");
			return false;
		}

		$select = ORM::forTable('users')->where('uuid', $data['uuid'])->findOne();

		$permissions = @json_decode($select->permissions, true);
		if(!$select || !array_key_exists($this->server->getData('hash'), $permissions)) {
			self::_setError("Invalid user was provided.");
			return false;
		}

		$permissions[$this->server->getData('hash')]['perms'] = self::_rebuildUserPermissions($data['permissions']);
		$permissions[$this->server->getData('hash')]['perms_gsd'] = self::_buildGSDPermissions($data['permissions']);
		$select->permissions = json_encode($permissions);
		$select->save();

		try {

			Unirest::put(
				"http://".$this->server->nodeData('ip').":".$this->server->nodeData('gsd_listen')."/gameservers/".$this->server->getData('gsd_id'),
				array(
					"X-Access-Token" => $this->server->nodeData('gsd_secret')
				),
				array(
					"keys" => json_encode(array(
						$permissions[$this->server->getData('hash')]['key'] => self::_buildGSDPermissions($data['permissions'])
					))
				)
			);

			return true;

		} catch(\Exception $e) {

			\Tracy\Debugger::log($e);
			self::_setError("An error occured when trying to update the server information on the daemon.");
			return false;

		}

	}

	/**
	 * Rebuilds user permissions to include the parent .view permission.
	 *
	 * @param array $data
	 * @return array
	 * @static
	 */
	protected final static function _rebuildUserPermissions(array $data) {

		foreach($data as $permission) {

			if(in_array($permission, array('files.edit', 'files.save', 'files.download', 'files.delete', 'files.create', 'files.upload', 'files.zip')) && !in_array('files.view', $data)) {
				$data = array_merge($data, array("files.view"));
			}

			if(in_array($permission, array('manage.rename.jar')) && !in_array('manage.rename.view', $data)) {
				$data = array_merge($data, array("manage.rename.view"));
			}

			if(in_array($permission, array('manage.ftp.details', 'manage.ftp.password')) && !in_array('manage.ftp.view', $data)) {
				$data = array_merge($data, array("manage.ftp.view"));
			}

		}

		return $data;

	}

	/**
	 * Builds an array of equivalent GSD permissions for each panel permission.
	 *
	 * @param array $data
	 * @return array
	 * @static
	 */
	protected final static function _buildGSDPermissions(array $data) {

		$gsd = array("s:console", "s:query");

		foreach($data as $permissionNode) {

			switch($permissionNode) {

				case "console.power":
					$gsd = array_merge($gsd, array("s:power"));
					break;
				case "console.commands":
					$gsd = array_merge($gsd, array("s:console:send"));
					break;
				case "files.view":
					$gsd = array_merge($gsd, array("s:files"));
					break;
				case "files.edit":
					$gsd = array_merge($gsd, array("s:files:get"));
					break;
				case "files.save":
					$gsd = array_merge($gsd, array("s:files:put"));
					break;
				case "files.zip":
					$gsd = array_merge($gsd, array("s:files:zip"));
					break;

			}

		}

		return $gsd;

	}

	/**
	 * Revokes subuser permissions for a given user that has an active account on the panel.
	 *
	 * @param object $orm Database query object.
	 * @return bool
	 */
	public function revokeActiveUserPermissions(ORM $orm) {

		if(!$this->avaliable($this->server->nodeData('ip'), $this->server->nodeData('gsd_listen'))) {
			self::_setError("Unable to access the server management daemon.");
			return false;
		}

		if(!array_key_exists($this->server->getData('hash'), json_decode($orm->permissions, true))) {
			self::_setError("Unable to locate server in user information.");
			return false;
		}

		$permissions = json_decode($this->server->getData('subusers'), true);
		unset($permissions[$orm->id]);

		$server = ORM::forTable('servers')->findOne($this->server->getData('id'));
		$server->subusers = json_encode($permissions);

		$userPermissions = json_decode($orm->permissions, true);

		try {

			Unirest::put(
				"http://".$this->server->nodeData('ip').":".$this->server->nodeData('gsd_listen')."/gameservers/".$this->server->getData('gsd_id'),
				array(
					"X-Access-Token" => $this->server->nodeData('gsd_secret')
				),
				array(
					"keys" => json_encode(array(
						$userPermissions[$this->server->getData('hash')]['key'] => array()
					))
				)
			);

			unset($userPermissions[$this->server->getData('hash')]);

			$orm->permissions = json_encode($userPermissions);
			$orm->save();
			$server->save();

			return true;

		} catch(\Exception $e) {

			\Tracy\Debugger::log($e);
			self::_setError("An error occured when trying to update the server information on the daemon.");
			return false;

		}

	}

}