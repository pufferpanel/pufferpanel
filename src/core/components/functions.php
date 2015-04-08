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
use \Exception;

/**
 * General Functions Trait
 */
trait Functions {

	use Authentication;

	/**
	 * Generates a unique FTP username
	 *
	 * @param string $base The base text to use for the username.
	 * @return string Returns the unique username beginning with mc-.
	 * @static
	 */
	public static function generateFTPUsername($base) {

		$base = str_replace(" ", "", $base);
		$i = strlen($base);
		if($i > 6) {
			$username = substr($base, 0, 6).'_'.self::keygen(5);
		} else {
			$username = $base.'_'.self::keygen((11 - $i));
		}

		return "pp-".strtolower($username);

	}

	/**
	 * Returns an array of ports from a 'range' seperated by '-'.
	 *
	 * @param string $input String of two ports using '-' to seperate.
	 * @return array Returns an array of integers. Returns null if first port is smaller than next.
	 */
	public static function processPorts($input) {

		$port_list = [];

		if(!trim($input)) {
			return $port_list;
		}

		foreach(explode(',', $input) as $range) {

			if(strpos($range, '-')) {

				$explode = explode('-', $range);
				if(!is_numeric($explode[0]) || !is_numeric($explode[1])) {
					throw new Exception("The range provided for ports in processPorts({$input}) is invalid. The range must be numeric.");
				}

				settype($explode[0], "int");
				settype($explode[1], "int");

				if($explode[0] > $explode[1]) {
					throw new Exception("The range provided for ports in processPorts({$input}) is invalid. The start value ({$explode[0]}) can't be higher than the end value ({$explode[1]}).");
				}

				for($i = $explode[0]; $i <= $explode[1]; $i++) {
					$port_list[] = $i;
				}

			} else {
				$port_list[] = $range;
			}

		}

		return $port_list;

	}

}