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
use \Exception;

/**
 * PufferPanel Core Language Class
 */
class Language extends User {

	/**
	 * @param string $language
	 */
	protected $language;

	/**
	 * @param array $loaded
	 */
	protected $loaded;

	/**
	 * Constructor class for Language
	 */
	public function __construct() {

		if(!$this->getData('language')) {
			$this->language = (isset($_COOKIE['pp_language']) && !empty($_COOKIE['pp_language'])) ? $_COOKIE['pp_language'] : Settings::config()->default_language;
		} else {
			$this->language = $this->getData('language');
		}

		if(!file_exists(APP_DIR.'languages/'.$this->language.'.json')) {
			throw new Exception('Unable to load the required language file! '.APP_DIR.'languages/'.$this->language.'.json');
		}

		$this->loaded = json_decode(file_get_contents(APP_DIR.'languages/'.$this->language.'.json'), true);

	}

	/**
	 * Returns the language value associated with a key. Graceful fallback to English if translated version doesn't exist.
	 *
	 * @param string $template The language key.
	 * @return string
	 */
	public function render($template) {

		if(!array_key_exists($template, $this->loaded)) {

			$load_english = json_decode(file_get_contents(APP_DIR.'languages/en.json'), true);

			if(array_key_exists($template, $load_english)) {
				return $load_english[$template];
			} else {
				return "{{ $template }}";
			}

		}

		return $this->loaded[$template];

	}

	/**
	 * Returns the loaded langauge as an array.
	 *
	 * @return array
	 * @deprecated
	 */
	public function loadTemplates() {

		return $this->loaded;

	}

 }

 ?>