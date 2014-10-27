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

	use Components\Database, Components\Authentication, Components\Page;

	/**
	 * Authentcation constructor class.
	 *
	 * @return void
	 */
	public function __construct()
		{

			$this->mysql = self::connect();
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

		$this->get = $this->mysql->prepare("SELECT `password` FROM `users` WHERE `email` = :email");
		$this->get->execute(array(
			':email' => $email
		));

			if($this->get->rowCount() == 1){

				$this->row = $this->get->fetch();
				return $this->password_compare($raw, $this->row['password']);

			}else
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

		$this->query = $this->mysql->prepare("SELECT * FROM `users` WHERE `session_ip` = :sessip AND `session_id` = :session");
		$this->query->execute(array(
			':sessip' => $ip,
			':session' => $session
		));

		if($this->query->rowCount() == 1){

			$this->row = $this->query->fetch();

			if($this->row['root_admin'] != 1 && $acp === true)
				return false;
			else{

				if($this->row['root_admin'] != '1'){

					if(!is_null($serverhash)){

						/*
						 * We have to do a mini-permissions building here since we can't call the user function from here
						 */
						if(!is_null($this->row['permissions']) && !empty($this->row['permissions']))
							$this->row['permissions'] = array_keys(json_decode($this->row['permissions'], true));
						else
							$this->row['permissions'] = array("0" => "0");

						$this->hashes = array_map(array($this->mysql, 'quote'), $this->row['permissions']);
						$this->_validateServer = $this->mysql->prepare("SELECT * FROM `servers` WHERE `hash` = :hash AND `owner_id` = :oid OR `hash` IN(".join(',', $this->hashes).") AND `active` = 1");
						$this->_validateServer->execute(array(
							':oid' => $this->row['id'],
							':hash' => $serverhash
						));

							if($this->_validateServer->rowCount() == 1)
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
