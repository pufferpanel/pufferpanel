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

/**
 * PufferPanel Core File Management Class
 */
class Files {

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
	 * @return double
	 */
	public function format($size, $precision = 0) {

		$base = log($size) / log(1024);

		return round(pow(1024, $base - floor($base)), $precision);

	}

	/**
	 * Reads a specified number of lines from a given file.
	 *
	 * @todo Remove function.
	 * @param string $filename
	 * @param int $lines
	 * @return void
	 */
	public function readLines($filename, $lines) {
	}

	/**
	 * Reads a specified number of lines from a given file beginning at the end of the file.
	 *
	 * @todo Remove function.
	 * @param string $path
	 * @param int $line_count
	 * @param int $block_size
	 * @return void
	 */
	function last_lines($path, $line_count, $block_size = 512) {
	}

}