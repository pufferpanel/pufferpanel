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
 * Page components trait to allow access in other classes and functions easily.
 */
trait Page {

	/**
	 * Redirects the user to a specified page automatically depending on headers.
	 *
	 * @deprecated deprecated since version 0.8.0
	 * @param string $url The URL to redirect the user to.
	 * @return void
	 * @static
	 */
	public static function redirect($url) {

		if(!headers_sent()){
			header('Location: '.urldecode($url));
			exit();
		}else
			exit('<meta http-equiv="refresh" content="0;url='.urldecode($url).'"/>');
		}

	}

	/**
	 * Generates a URL.
	 *
	 * @deprecated deprecated since version 0.8.0
	 * @return string
	 * @static
	 */
	public static function genRedirect(){

		$https = (isset($_SERVER['HTTPS'])) ? 'https://' : 'http://';
		return $https.$_SERVER['HTTP_HOST'].$_SERVER['REQUEST_URI'];

	}

	/**
	 * Fixes issue with twig not accepting empty array values
	 *
	 * @deprecated deprecated since version 0.8.0
	 * @return array
	 * @static
	 */
	public static function twigGET() {

		$vars = $_GET;
		$return = array();
		foreach($vars as $id => $value){

			if(empty($value))
				$value = true;

			$return[$id] = $value;
		}

		return $return;

	}

}

?>