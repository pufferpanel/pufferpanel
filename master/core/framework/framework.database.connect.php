<?php

/*
 *Database Connection
 */

class dbConn {
 
	protected static $db;
	 
	private function __construct() {
	 	
	 	require('master_configuration.php');
		try {

			self::$db = new PDO('mysql:host='.$_INFO['sql_h'].';dbname='.$_INFO['sql_db'], $_INFO['sql_u'], $_INFO['sql_p'], array(
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

?>