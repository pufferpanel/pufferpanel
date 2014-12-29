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
use \ORM;

class Users extends \PufferPanel\Core\Email {

	use \PufferPanel\Core\Components\Authentication;

	protected static $_data;
	protected static $_gsd;
	protected static $_id;
	protected static $_permission;
	protected static $_permissionNode;
	protected $_server;

	/**
	 * Constructor class for \PufferPanel\Core\Router\Node_Users
	 *
	 * @return void
	 */
	public function __construct($server) {

		$this->_server = $server;
		$this->settings = new \PufferPanel\Core\Settings();

	}

	/**
	 * Handles adding a new subuser to a server.
	 *
	 * @return bool
	 */
	public function addSubuser($data = array()) {

		$this->registerToken = $this->keygen(32);

		$this->find = ORM::forTable('users')->where('email', $data['email'])->findOne();
		if(!$this->find) {

			$this->newAccount = ORM::forTable('account_change')->create();
			$this->newAccount->set(array(
				'type' => 'user_register',
				'content' => $data['email'],
				'key' => $this->registerToken,
				'time' => time()
			));
			$this->newAccount->save();

		}

		$this->subuserToken = $this->keygen(32);
		$this->addToken = ORM::forTable('account_change')->create();
		$this->addToken->set(array(
			'type' => 'subuser',
			'content' => json_encode(array(
				$this->_server->getData('hash') => array(
					'key' => $this->generateUniqueUUID('servers', 'gsd_secret'),
					'perms' => self::_rebuildUserPermissions($data['permissions']),
					'perms_gsd' => self::_buildGSDPermissions($data['permissions'])
				)
			)),
			'key' => $this->subuserToken,
			'time' => time()
		));
		$this->addToken->save();

		$this->subusers = (!empty($this->_server->getData('subusers'))) ? json_decode($this->_server->getData('subusers'), true) : array();
		$this->subusers[$data['email']] = $this->subuserToken;

		$this->updateServer = ORM::forTable('servers')->findOne($this->_server->getData('id'));
		$this->updateServer->subusers = json_encode($this->subusers);
		$this->updateServer->save();

		/*
		* Send Email
		*/
		$this->buildEmail('new_subuser', array(
			'REGISTER_TOKEN' => $this->registerToken,
			'SUBUSER_TOKEN' => $this->subuserToken,
			'URLENCODE_TOKEN' => urlencode($this->registerToken),
			'SERVER' => $this->_server->getData('name'),
			'EMAIL' => $data['email']
		))->dispatch($data['email'], $this->settings->get('company_name').' - You\'ve Been Invited to Manage a Server');

		return true;

	}

	/**
	 * Rebuilds user permissions to include the parent .view permission.
	 *
	 * @param array $data
	 * @return array
	 * @static
	 */
	protected final static function _rebuildUserPermissions($data = array()) {

		self::$_data = $data;

		foreach(self::$_data as self::$_id => self::$_permission) {

			if(in_array(self::$_permission, array('files.edit', 'files.save', 'files.download', 'files.delete', 'files.create', 'files.upload', 'files.zip')) && !in_array('files.view', self::$_data)) {
				self::$_data = array_merge(self::$_data, array("files.view"));
			}

			if(in_array(self::$_permission, array('manage.rename.jar')) && !in_array('manage.rename.view', self::$_data)) {
				self::$_data = array_merge(self::$_data, array("manage.rename.view"));
			}

			if(in_array(self::$_permission, array('manage.ftp.details', 'manage.ftp.password')) && !in_array('manage.ftp.view', self::$_data)) {
				self::$_data = array_merge(self::$_data, array("manage.ftp.view"));
			}

		}

		return self::$_data;

	}

	/**
	 * Builds an array of equivalent GSD permissions for each panel permission.
	 *
	 * @param array $data
	 * @return array
	 * @static
	 */
	protected final static function _buildGSDPermissions($data = array()) {

		self::$_data = $data;
		self::$_gsd = array("s:console", "s:query");

		foreach(self::$_data as self::$_id => self::$_permissionNode) {

			switch(self::$_permissionNode) {

				case "console.power":
					self::$_gsd = array_merge(self::$_gsd, array("s:power"));
					break;
				case "console.commands":
					self::$_gsd = array_merge(self::$_gsd, array("s:console:command"));
					break;
				case "files.view":
					self::$_gsd = array_merge(self::$_gsd, array("s:files"));
					break;
				case "files.edit":
					self::$_gsd = array_merge(self::$_gsd, array("s:files:get"));
					break;
				case "files.save":
					self::$_gsd = array_merge(self::$_gsd, array("s:files:put"));
					break;
				case "files.zip":
					self::$_gsd = array_merge(self::$_gsd, array("s:files:zip"));
					break;

			}

		}

		return self::$_gsd;

	}

}