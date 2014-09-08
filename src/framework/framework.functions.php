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

namespace Functions;

/**
 * General Functions Trait
 */
trait general {

	use \Auth\components;

	/**
	 * Generates a unique FTP username
	 *
	 * @param string $base The base text to use for the username.
	 * @return string Returns the unique username beginning with mc-.
	 */
	public static function generateFTPUsername($base) {

		$username = (strlen($base) > 6) ? substr($base, 0, 6).'_'.self::keygen(5) : $base.'_'.self::keygen((11 - strlen($base)));
	    return "mc-".$username;

	}

}

?>
