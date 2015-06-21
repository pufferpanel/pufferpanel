<?php

/*
  PufferPanel - A Game Server Management Panel
  Copyright (c) 2015 PufferPanel

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
 * URL utility class
 */
trait Url {

	/**
	 * Strips http:// from a URL.
	 * If true is passed, https:// will also be stripped.
	 *
	 * @param string $source String to strip
	 * @param bool $stripHttps Strip HTTPS as well (default false)
	 */
	public static function stripHttp($source, $stripHttps = false) {

		$limit = 1;

		if(stripos($source, "http://") == 0) {
			$source = str_ireplace("http://", "", $source, $limit);
		}

		if($stripHttps && stripos($source, "https://") == 0) {
			$source = str_ireplace("https://", "", $source, $limit);
		}

		return $source;
		
	}

	/**
	 * Convience method to remove trailing / from URL.
	 *
	 * @param type $source
	 * @return type
	 */
	public static function stripTrailing($source) {
		return rtrim($source, '/ ');
	}

	/**
	 * Convience method to add trailing / to URL, if one does not exist
	 *
	 * @param type $source
	 * @return type
	 */
	public static function addTrailing($source) {
		return self::stripTrailing($source).'/';
	}

}
