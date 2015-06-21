<?php
/*
	PufferPanel - A Game Server Management Panel
	Copyright (c) 2015 Dane Everitt

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
 * PufferPanel Core File Management Class
 */
class Files {

	/**
	 * Returns an array of files that are allowed to be edited through the panel.
	 *
	 * @return array
	 */
	public final function editable() {

		return array(
			'txt',
			'yml',
			'yaml',
			'log',
			'conf',
			'config',
			'html',
			'json',
			'properties',
			'props',
			'cfg',
			'lang',
			'ini',
			'cmd',
			'sh',
			'lua',
			'0' // Supports BungeeCord Files
		);

	}

	/**
	 * Converts from bytes into the largest possible size that is still readable.
	 *
	 * @param int $bytes
	 * @param int $decimals Defaults to 2 decimal places.
	 * @return string
	 */
	public function formatSize($bytes, $decimals = 2) {

		$sz = explode(',', 'B,KB,MB,GB');
		$factor = floor((strlen($bytes) - 1) / 3);

		return sprintf("%.{$decimals}f", $bytes / pow(1024, $factor)).' '.$sz[$factor];

	}

	/**
	 * Formats the size to a certain precision.
	 *
	 * @param int $size
	 * @param int $precision
	 * @return int
	 */
	public function format($size, $precision = 0) {

		$base = log($size) / log(1024);

		return round(pow(1024, $base - floor($base)), $precision);

	}

}