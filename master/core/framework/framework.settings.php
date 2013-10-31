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
 
class getSettings extends dbConn {

	public function __construct()
		{
		
			$this->mysql = parent::getConnection();
		
			$this->query = $this->mysql->prepare("SELECT * FROM `acp_settings`");
			$this->query->execute();
				
				while($this->row = $this->query->fetch()){
					
					$this->_data[$this->row['setting_ref']] = $this->row['setting_val'];
					
				}
		
		}
		
	public function get($setting)
		{
		
			return (array_key_exists($setting, $this->_data)) ? $this->_data[$setting] : '_notfound_';
		
		}

}

?>