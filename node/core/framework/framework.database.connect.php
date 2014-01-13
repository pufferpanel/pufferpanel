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
 
/*
 *Database Connection
 */

class dbConn {
 
	protected static $db;
	 
	private function __construct() {
	 	
	 	require('node_configuration.php');
		try {

			self::$db = new PDOEx('mysql:host='.$_INFO['sql_h'].';dbname='.$_INFO['sql_db'], $_INFO['sql_u'], $_INFO['sql_p'], array(
	    		PDO::ATTR_PERSISTENT => true,
	    		PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
			));
	
			self::$db->setAttribute( PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION );
	
		}catch (PDOException $e) {
	
			echo "MySQL Connection Error: " . $e->getMessage();
	
		}
	 
	}
	 
	public static function getConnection() {
	 
		if (!self::$db) {
	
			new dbConn();
	
		}
	 
		return self::$db;

	}
		 
}

class PDOEx extends PDO {
	
	private static $queryCounter = 0;
		
	public function query($query)
    {
        ++self::$queryCounter;
        return parent::query($query);
    }
    
    public function prepare($statement, $options = array())
    {
        ++self::$queryCounter;
        return parent::prepare($statement, $options);
    }
    
    public function exec($statement)
    {
        ++self::$queryCounter;
        return parent::exec($statement);
    }
    
    public function getCount(){
    	return self::$queryCounter;
    }
	
}

?>