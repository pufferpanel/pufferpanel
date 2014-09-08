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
namespace Language;

/**
 * PufferPanel Core Language Class
 */
class lang {

	/**
	 * Constructor class for language
	 *
	 * @param string $language The language to load.
	 * @return void
	 */
	public function __construct($language){

		$this->_language = $language;

		if(!file_exists(__DIR__.'/lang/'.$language.'.json'))
			exit('Unable to load the required language file! lang/'.$language.'.json');

		$this->loaded_language = json_decode(file_get_contents(__DIR__.'/lang/'.$language.'.json'), true);

	}

	/**
	 * Returns the language value associated with a key. Graceful fallback to English if translated version doesn't exist.
	 *
	 * @param string $template The language key.
	 * @return string
	 */
	public function tpl($template){

		$template = str_replace(".", "_", $template);

		if(array_key_exists($template, $this->loaded_language))
			return $this->loaded_language[$template];
		else{

			if($this->_language == "en")
				return $template;
			else{

				$this->load_english = json_decode(file_get_contents(__DIR__.'/lang/en.json'), true);

					if(array_key_exists($template, $this->load_english))
						return $this->load_english[$template];
					else
						return $template;

			}

		}

	}

	/**
	 * Returns the loaded langauge as an array.
	 *
	 * @return array
	 */
	public function loadTemplates() {

		return $this->loaded_language;

	}

 }

 ?>