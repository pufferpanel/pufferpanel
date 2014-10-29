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
use \ORM as ORM;

/*
 * TOTP Class
 */
use Otp\Otp;
use Otp\GoogleAuthenticator;
use Base32\Base32;

/**
 * PufferPanel Core Authentication Class
 */
class Authentication {

	use Components\Authentication, Components\Page;

	/**
	 * Authentcation constructor class.
	 *
	 * @return void
	 */
	public function __construct()
		{

			$this->settings = new Settings();

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

		if($otp->checkTotp(Base32::decode($secret), $key))
			return true;
		else
			return false;

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

		if($this->get !== false)
			return $this->password_compare($raw, $this->get->password);
		else
			return false;

	}

	/**
	 * Checks if a user is currently logged in or if their session is expired.
	 *
	 * @param string $ip
	 * @param string $session
	 * @param string $serverhash
	 * @param int $acp
	 * @return bool
	 */
	public function isLoggedIn($ip, $session, $serverhash = null, $acp = false){

		$this->user = ORM::forTable('users')->where(array('session_ip' => $ip, 'session_id' => $session))->findOne();

		if($this->user !== false){

			if($this->user->root_admin != 1 && $acp === true)
				return false;
			else{

				if($this->user->root_admin != 1){

					if(!is_null($serverhash)){

						/*
						 * We have to do a mini-permissions building here since we can't call the user function from here
						 */
						if(!is_null($this->user->permissions) && !empty($this->user->permissions))
							$this->permissions = array_keys(json_decode($this->user->permissions, true));
						else
							$this->permissions = array("0" => "0");

						$this->server = ORM::forTable('servers')->where(array('hash' => $hash, 'active' => 1))->where_raw('`owner_id` = ? OR `hash` IN(?)', array($this->user->id, join(',', $this->permissions)))->findOne();

							if($this->server !== false)
								return true;
							else
								return false;

					}else
						return true;

				}else
					return true;

			}

		}else
			return false;

	}

}

?>