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

namespace Database;

trait database {

	protected static $db;
	public static $salt;
	
	public static function buildConnection(){
	
		require('configuration.php');
		try {
			
			/*
			 * Connect to SQL Server over SSL
			 */
			if(array_key_exists('sql_ssl', $_INFO) && $_INFO['sql_ssl'] === true){
			
				self::$db = new databaseInit('mysql:host='.$_INFO['sql_h'].';dbname='.$_INFO['sql_db'], $_INFO['sql_u'], $_INFO['sql_p'],
					array(
				        \PDO::MYSQL_ATTR_SSL_KEY => $_INFO['sql_ssl_client-key'],
				        \PDO::MYSQL_ATTR_SSL_CERT => $_INFO['sql_ssl_client-cert'],
				        \PDO::MYSQL_ATTR_SSL_CA => $_INFO['sql_ssl_ca-cert'],
						\PDO::ATTR_PERSISTENT => true,
						\PDO::ATTR_DEFAULT_FETCH_MODE => \PDO::FETCH_ASSOC
					)
				);
			
			}else{
			
				self::$db = new databaseInit('mysql:host='.$_INFO['sql_h'].';dbname='.$_INFO['sql_db'], $_INFO['sql_u'], $_INFO['sql_p'], array(
		    		\PDO::ATTR_PERSISTENT => true,
		    		\PDO::ATTR_DEFAULT_FETCH_MODE => \PDO::FETCH_ASSOC
				));
			
			}
	
			self::$db->setAttribute( \PDO::ATTR_ERRMODE, \PDO::ERRMODE_EXCEPTION );
	
		}catch (PDOException $e) {
	
			echo "MySQL Connection Error: " . $e->getMessage();
	
		}
	
	}
	
	public static function connect() {
		 
		if (!self::$db) {
	
			self::buildConnection();
	
		}
	 
		return self::$db;
	
	}
	
}

class databaseInit extends \PDO {

	use database;
	
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
	
	public static function getCount(){
		return self::$queryCounter;
	}
	
	public static function returnStartTime(){
	
		return microtime(true);
	
	}

}
?>