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
namespace PufferPanel\Core;
use \ORM, \Otp\Otp, \Base32\Base32, \PufferPanel\Core\Config\DatabaseConfig;

/**
 * PufferPanel Core Authentication Class
 */
class Authentication {

	use Components\Authentication, Components\Page;

	protected $settings;

	protected $authenticated = false;

	protected $select;

	/**
	 * Authentcation constructor class.
	 *
	 * @return void
	 */
	public function __construct(){

		$this->settings = new DatabaseConfig('acp_settings', 'setting_ref', 'setting_val');

		$this->select = (!isset($_COOKIE['pp_auth_token']) || empty($_COOKIE['pp_auth_token'])) ? false : ORM::forTable('users')->where(array('session_ip' => $_SERVER['REMOTE_ADDR'], 'session_id' => $_COOKIE['pp_auth_token']))->findOne();

		$this->authenticated = (!$this->select) ? false : true;

	}

	/**
	 * Validates a TOTP request.
	 *
	 * @todo Prevent TOTP replay attack.
	 * @param string $key The TOTP token sent in.
	 * @param string $secret The TOTP secret.
	 * @return bool
	 */
	public function validateTOTP($key, $secret){

		$otp = new Otp();

		return $otp->checkTotp(Base32::decode($secret), $key);

	}

	/**
	 * Verifys the a user is entering the correct password for their account.
	 *
	 * @param string $email
	 * @param string $raw The raw password.
	 * @return bool
	 */
	public function verifyPassword($email, $raw){

		$this->get = ORM::forTable('users')->select('password')->where('email', $email)->findOne();

		if($this->get !== false) {
			return $this->password_compare($raw, $this->get->password);
		} else {
			return false;
		}

	}

	/**
	 * Returns the authentication status of a user.
	 *
	 * @return bool
	 */
	public final function isLoggedIn() {

		return $this->authenticated;

	}

	/**
	 * Returns the admin status of a user.
	 *
	 * @return bool
	 */
	public final function isAdmin() {

		if(!$this->select) {
			return false;
		} else {
			return ($this->select->root_admin == 1) ? true : false;
		}

	}

	/**
	 * Checks if the selected server belongs to the user.
	 *
	 * @return bool
	 */
	public final function isServer() {

		if(!isset($_COOKIE['pp_server_hash']) || empty($_COOKIE['pp_server_hash'])) {
			return false;
		}

		$query = ORM::forTable('servers')->where(array('hash' => $_COOKIE['pp_server_hash'], 'active' => 1));

		if(!$this->isAdmin()) {

			$permissions = (!empty($this->select->permissions)) ? array_keys(json_decode($this->select->permissions, true)) : array();

			$query->where_raw('`owner_id` = ? OR `hash` IN(?)', array($this->select->id, join(',', $permissions)));

		}

		return (!$query->findOne()) ? false : true;

	}

}