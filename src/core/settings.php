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
use \ORM as ORM;

/**
 * PufferPanel Core Settings Class File
 */
class Settings {

	/**
	 * Constructor class for building settings data.
	 *
	 * @return void
	 */
	public function __construct()
		{

			$this->settings = ORM::forTable('acp_settings')->findMany();

			foreach($this->settings as $this->setting)
				$this->_data[$this->setting->setting_ref] = $this->setting->setting_val;

			$this->nodes = ORM::forTable('nodes')->select('id')->select('node')->findMany();

			foreach($this->nodes as $this->node)
					$this->_dataNode[$this->node->id] = $this->node->node;

		}

	/**
	 * Function to retrieve various panel settings.
	 *
	 * @param string $setting The name of the setting for which you want the value.
	 * @return mixed This will return the column data for the setting, or if $setting was left blank all settings in an array.
	 */
	public function get($setting = null)
		{

			if(is_null($setting))
				return $this->_data;
			else
				return (array_key_exists($setting, $this->_data)) ? $this->_data[$setting] : '_notfound_';

		}

	/**
	 * Convert a node ID into a name for the node.
	 *
	 * @param int $id The ID of the node you want the name for.
	 * @return string The name of the node.
	 */
	public function nodeName($id)
		{

			return (array_key_exists($id, $this->_dataNode)) ? $this->_dataNode[$id] : 'unknown';

		}

}
?>
