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

		$username = (strlen($base) > 6) ? substr($base, 0, 6).'_'.self::keygen(5) : $base.'_'.self::keygen((11 - strlen($base)));
	    return "mc-".strtolower($username);

	}

	/**
	 * Returns an array of ports from a 'range' seperated by '-'.
	 *
	 * @param string $range String of two ports using '-' to seperate.
	 * @return array Returns an array of integers. Returns null if first port is smaller than next.
	 */
	public static function processPorts($ports) {

		$portList = [];

		if(!strpos($ports, ",")) {

			// Possible a Range, or a Single Port
			if(!strpos($ports, "-")) {

				$portList[] = $ports;

			} else {

				$split = explode("-",$ports);
				$rangeA = intval($split[0]);
				$rangeB = intval($split[1]);

				if($rangeA < $rangeB) {

					while($rangeA <= $rangeB) {
						$portList[] = $rangeA;
						$rangeA++;
					}

				} else {

					error_log('Ports are in wrong order! Range 1 must be a number > Range 2.',0);
					error_log('Port Range 1: ' . $rangeA . ' Port Range 2: ' . $rangeB,0);
					return null;

				}

			}

		} else {

			// Possible Mix of Ranges and Single Ports
			if(!strpos($ports, "-")) {

				// No ranges
				$portList = array_merge($portList, explode(",", $ports));

			} else {

				// Singles Mixed with Range
				foreach(explode(",", $ports) as $id => $range) {

					if(!strpos($range, "-")) {

						$portList = array_merge($portList, array($range));

					} else {

						$split = explode("-",$range);
						$rangeA = intval($split[0]);
						$rangeB = intval($split[1]);

						if($rangeA < $rangeB) {

							while($rangeA <= $rangeB) {
								$portList[] = $rangeA;
								$rangeA++;
							}

						} else {

							error_log('Ports are in wrong order! Range 1 must be a number > Range 2.',0);
							error_log('Port Range 1: ' . $rangeA . ' Port Range 2: ' . $rangeB,0);
							return null;

						}

					}

				}

			}

		}

		return $portList;
	}
}

?>
