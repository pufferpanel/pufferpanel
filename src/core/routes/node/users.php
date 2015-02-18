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
use \ORM, \Unirest, \PufferPanel\Core\Settings;

class Users extends \PufferPanel\Core\Email {

	use \PufferPanel\Core\Components\Authentication, \PufferPanel\Core\Components\GSD, \PufferPanel\Core\Components\Error_Handler;

	/**
	 * @param object
	 */
	protected $server;

	/**
	 * Constructor class for \PufferPanel\Core\Router\Node\Users
	 */
	public function __construct(\PufferPanel\Core\Server $server) {

		$this->server = $server;

	}

	/**
	 * Handles adding a new subuser to a server.
	 *
	 * @return bool
	 */
	public function addSubuser(\Klein\DataCollection\DataCollection $data) {

		if(is_null($data->permissions) || !$data->permissions) {
			return false;
		}

		$registerToken = $this->keygen(32);
		$subuserToken = $this->keygen(32);
		$subuserUUID = $this->generateUniqueUUID('subusers', 'uuid');
		$gsdSecret = $this->generateUniqueUUID('servers', 'gsd_secret');
		$data->permissions = self::_rebuildUserPermissions($data->permissions);

		$find = ORM::forTable('users')->where('email', $data->email)->findOne();
		if(!$find) {

			$newAccount = ORM::forTable('account_change')->create();
			$newAccount->set(array(
				'type' => 'user_register',
				'content' => $data->email,
				'key' => $registerToken,
				'time' => time()
			));
			$newAccount->save();

			$newAccount = ORM::forTable('account_change')->create();
			$newAccount->set(array(
				'type' => 'subuser',
				'content' => $subuserUUID,
				'key' => $subuserToken,
				'time' => time()
			));
			$newAccount->save();

			/*
			* Send Email
			*/
			$this->buildEmail('new_subuser_createaccount', array(
				'REGISTER_TOKEN' => $registerToken,
				'SUBUSER_TOKEN' => $subuserToken,
				'URLENCODE_TOKEN' => urlencode($registerToken),
				'SERVER' => $this->server->getData('name'),
				'EMAIL' => $data->email
			))->dispatch($data->email, Settings::config()->company_name.' - You\'ve Been Invited to Manage a Server');

		}

		$addToken = ORM::forTable('subusers')->create();
		$addToken->set(array(
			'uuid' => $subuserUUID,
			'user' => (!$find) ? "-1" : $find->id,
			'server' => $this->server->getData('id'),
			'gsd_secret' => $gsdSecret,
			'gsd_permissions' => json_encode(self::_buildGSDPermissions($data->permissions)),
			'permissions' => (!$find) ? json_encode($data->permissions) : null,
			'pending' => (!$find) ? 1 : 0,
			'pending_email' => (!$find) ? $data->email : null
		));
		$addToken->save();

		try {

			Unirest\Request::put(
				"https://".$this->server->nodeData('ip').":".$this->server->nodeData('gsd_listen')."/gameservers/".$this->server->getData('gsd_id'),
				array(
					"X-Access-Token" => $this->server->nodeData('gsd_secret')
				),
				array(
					"keys" => json_encode(array(
						$gsdSecret => self::_buildGSDPermissions($data->permissions)
					))
				)
			);

			if($find) {

				/*
				* Send Email
				*/
				$this->buildEmail('new_subuser', array(
					'SERVER' => $this->server->getData('name'),
					'EMAIL' => $data->email
				))->dispatch($data->email, Settings::config()->company_name.' - You\'ve Been Invited to Manage a Server');

			}
			return true;

		} catch(\Exception $e) {

			\Tracy\Debugger::log($e);
			self::_setError("An error occured when trying to update the server information on the daemon.");
			return false;

		}

	}

	/**
	 * Updates information for a subuser.
	 *
	 * @param array $data
	 * @return bool
	 */
	public function modifySubuser(\Klein\DataCollection\DataCollection $data) {

		if(!$this->avaliable($this->server->nodeData('ip'), $this->server->nodeData('gsd_listen'))) {
			self::_setError("Unable to access the server management daemon.");
			return false;
		}

		$select = ORM::forTable('subusers')->where('uuid', $data->uuid)->findOne();

		if(!$select) {
			self::_setError("Invalid user was provided.");
			return false;
		}

		ORM::forTable('permissions')->where(array(
			'user' => $select->user,
			'server' => $select->server
		))->deleteMany();

		foreach(self::_rebuildUserPermissions($data->permissions) as $id => $permission) {

			ORM::forTable('permissions')->create()->set(array(
				'user' => $select->user,
				'server' => $select->server,
				'permission' => $permission
			))->save();

		}

		$select->gsd_permissions = json_encode(self::_buildGSDPermissions($data->permissions));
		$select->save();

		try {

			Unirest\Request::put(
				"https://".$this->server->nodeData('ip').":".$this->server->nodeData('gsd_listen')."/gameservers/".$this->server->getData('gsd_id'),
				array(
					"X-Access-Token" => $this->server->nodeData('gsd_secret')
				),
				array(
					"keys" => json_encode(array(
						$select->gsd_secret => self::_buildGSDPermissions($data->permissions)
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

		\Tracy\Debugger::log($data);
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

			if(in_array($permission, array('manage.ftp.details', 'manage.ftp.password', 'manage.rename.jar', 'manage.ftp.view', 'manage.rename.view')) && !in_array('manage.view', $data)) {
				$data = array_merge($data, array("manage.view"));
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

		if($orm->pending == 0) {

			try {

				Unirest\Request::put(
					"https://".$this->server->nodeData('ip').":".$this->server->nodeData('gsd_listen')."/gameservers/".$this->server->getData('gsd_id'),
					array(
						"X-Access-Token" => $this->server->nodeData('gsd_secret')
					),
					array(
						"keys" => json_encode(array(
							$orm->gsd_secret => array()
						))
					)
				);

				ORM::forTable('permissions')
					->where(array(
						'user' => $orm->user,
						'server' => $this->server->getData('id')
					))
					->delete_many();

				$orm->delete();

				return true;

			} catch(\Exception $e) {

				\Tracy\Debugger::log($e);
				self::_setError("An error occured when trying to update the server information on the daemon.");
				return false;

			}

		} else {

			ORM::forTable('account_change')->where('content', $orm->uuid)->findOne()->delete();
			ORM::forTable('subusers')->where('uuid', $orm->uuid)->findOne()->delete();
			$orm->delete();
			return true;

		}

	}

}