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
namespace PufferPanel\Core\Components;
use \ORM as ORM;

/**
* PufferPanel Core Components Trait
*/
trait Authentication {

	/**
	* Returns a hashed version of the raw string that is passed. Use for password hashing.
	*
	* @param string $raw The raw password.
	* @return string The hashed password.
	*/
	public function hash($raw){

		return password_hash($raw, PASSWORD_BCRYPT);

	}

	/**
	* Compares a password to the hashed version to see if they match.
	*
	* @param string $raw The raw password.
	* @param string $hashed The hashed password.
	* @return bool Returns true if the password matches.
	*/
	private function password_compare($raw, $hashed){

		return password_verify($raw, $hashed);

	}

	/**
	* Generates an OpenSSL Encryption initalization vector.
	*
	* @return string
	*/
	public function generate_iv(){

		return base64_encode(mcrypt_create_iv(mcrypt_get_iv_size(MCRYPT_CAST_256, MCRYPT_MODE_CBC), MCRYPT_RAND));

	}

	/**
	* Encrypts a given string using an IV and AES-256-CBC encryption.
	*
	* @param string $raw The raw string to be encrypted.
	* @param string $iv The initalization vector to use.
	* @return string
	* @static
	*/
	public static function encrypt($raw, $iv){

		return openssl_encrypt($raw, 'AES-256-CBC', file_get_contents(HASH), false, base64_decode($iv));

	}

	/**
	* Decrypts a given string using an IV and defined method.
	*
	* @param string $encrypted The encrypted string that you want decrypted.
	* @param string $iv The initalization vector to use.
	* @param string $method Defaults to AES-256-CBC but you can define any other valid encryption method.
	* @return string
	* @static
	*/
	public static function decrypt($encrypted, $iv, $method = 'AES-256-CBC'){

		return openssl_decrypt($encrypted, $method, file_get_contents(HASH), 0, base64_decode($iv));

	}

	/**
	* Generate RFC 4122 Compliant v4 UUIDs
	*
	* @return string Returns a RFC 412 compliant UUID in the format XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
	* @static
	*/
	public static function gen_UUID(){

		return sprintf('%04x%04x-%04x-%04x-%04x-%04x%04x%04x',
			mt_rand(0, 0xffff),
			mt_rand(0, 0xffff),
			mt_rand(0, 0xffff),
			mt_rand(0, 0x0fff) | 0x4000,
			mt_rand(0, 0x3fff) | 0x8000,
			mt_rand(0, 0xffff),
			mt_rand(0, 0xffff),
			mt_rand(0, 0xffff)
		);

	}

	/**
	* Returns a valid UUID that is not currently in use.
	*
	* @param string $database
	* @param string $column
	* @return string
	*/
	public function generateUniqueUUID($database, $column){

		$uuid = self::gen_UUID();

		$check = ORM::forTable($database)->where($column, $uuid)->findOne();

		if(!$check) {
			return $uuid;
		} else {
			$this->generateUniqueUUID($database, $column);
		}

	}

	/**
	* Generates a random string of characters.
	*
	* @param int $amount
	* @return string
	* @static
	*/
	public static function keygen($amount){

		$character_set = "abcdefghijklmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ0123456789";
		$random = null;
		$max = (strlen($character_set) - 1);

		for ($i=0; $i < $amount; $i++) {
			$random .= $character_set[mt_rand(0, $max)];
		}

		return str_shuffle($random);

	}

	/**
	* Handles validating that a user's password meets the requirements for being changed.
	*
	* @param string $password
	* @param string $regex Optional parameter to define your own regex for checking password requirements.
	* @return bool
	*/
	public function validatePasswordRequirements($password, $regex = "#.*^(?=.{8,200})(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9]).*$#") {

		return (preg_match($regex, $password)) ? true : false;

	}

	/**
	* Returns the specified cookie.
	*
	* @param string $cookie
	* @return mixed
	*/
	public function getCookie($cookie){

		if(isset($_COOKIE[$cookie])){
			return $_COOKIE[$cookie];
		}else{
			return null;
		}

	}

	/**
	* Creates a XSRF Token
	*
	* @param mixed $token
	* @param mixed $identifier
	* @return mixed
	*/
	public function XSRF($token = null, $identifier = null, $reset = false){

		$cookie = "pp_xsrf_token".$identifier;

		if(!is_null($token)) {

			return (isset($_COOKIE[$cookie]) && $_COOKIE[$cookie] == $token) ? true : false;

		} else {

			if(!isset($_COOKIE[$cookie]) || $reset) {

				$xsrf = base64_encode(openssl_random_pseudo_bytes(32));
				setcookie($cookie, $xsrf, 0, '/');
				return '<input type="hidden" name="xsrf'.$identifier.'" value="'.$xsrf.'" />';

			} else {
				return '<input type="hidden" name="xsrf'.$identifier.'" value="'.$_COOKIE[$cookie].'" />';
			}

		}

	}

}
