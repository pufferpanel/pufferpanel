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
	 * Constructor class for language
	 *
	 * @return void
	 */
	public function __construct(){

		$s = new Config\DatabaseConfig('acp_settings', 'setting_ref');

		if(!$this->getData('language')) {
			$this->language = (isset($_COOKIE['pp_language']) && !empty($_COOKIE['pp_language'])) ? $_COOKIE['pp_language'] : $s->config('default_language');
		} else {
			$this->language = $this->getData('language');
		}

		if(!file_exists(dirname(__DIR__).'/lang/'.$this->language.'.json')) {
			throw new Exception('Unable to load the required language file! lang/'.$this->language.'.json');
		}

		$this->loaded = json_decode(file_get_contents(dirname(__DIR__).'/lang/'.$this->language.'.json'), true);

	}

	/**
	 * Returns the language value associated with a key. Graceful fallback to English if translated version doesn't exist.
	 *
	 * @param string $template The language key.
	 * @return string
	 */
	public function render($template){

		$template = str_replace(".", "_", $template);

		if(array_key_exists($template, $this->loaded)) {
			return $this->loaded[$template];
		}

		if($this->language == "en") {
			return $template;
		}

		$load_english = json_decode(file_get_contents(__DIR__.'/lang/en.json'), true);

		if(array_key_exists($template, $load_english)) {
			return $load_english[$template];
		}

		return $template;

	}

	/**
	 * Returns the loaded langauge as an array.
	 *
	 * @return array
	 */
	public function loadTemplates() {

		return $this->loaded;

	}

 }

 ?>