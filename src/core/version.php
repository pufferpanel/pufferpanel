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

namespace PufferPanel\Core;

/**
 * PufferPanel versioning class
 */
class Version {

	private static $sha;

	/**
	 * Returns the version string, including the current version and sha
	 * in the format "Version (sha)"
	 */
	public static function get() {
		return trim(file_get_contents(SRC_DIR.'versions/current'));
	}

}
