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

/**
 * PufferPanel Core Settings Class File
 */
class settings {

	use Database\database;

    /**
     * Constructor class for building settings data.
     *
     * @return void
     */
	public function __construct()
		{

			$this->mysql = self::connect();

			$this->query = $this->mysql->prepare("SELECT * FROM `acp_settings`");
			$this->query->execute();

				while($this->row = $this->query->fetch()){

					$this->_data[$this->row['setting_ref']] = $this->row['setting_val'];

				}

			$this->queryNode = $this->mysql->prepare("SELECT `id`, `node` FROM `nodes`");
			$this->queryNode->execute();

				while($this->rowNode = $this->queryNode->fetch()){

					$this->_dataNode[$this->rowNode['id']] = $this->rowNode['node'];

				}

		}

    /**
     * Function to retrieve various panel settings.
     *
     * @param string $setting The name of the setting for which you want the value.
     * @return array|string This will return the column data for the setting, or if $setting was left blank all settings in an array.
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
