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
 * PufferPanel Page Actions Function File
 */
namespace Page;

trait components {
	
	public static function redirect($url, $manual = false) {
		
		/*
		 * Generate URL
		 */
		if(!headers_sent()){
			header('Location: '.urldecode(self::appendToURL($manual).$url));
			exit();
		}else
			exit('<meta http-equiv="refresh" content="0;url='.urldecode(self::appendToURL($manual).$url).'"/>');
		
	}
    
    public static function appendToURL($manual){
        
        if($manual === false){
        
	        $https = (isset($_SERVER['HTTPS'])) ? 'https://' : 'http://';
        	return $https.$_SERVER['HTTP_HOST'].str_replace("index.php", "", $_SERVER['PHP_SELF']);
        
        }else
        	return null;
        
    }
    
    /*
     * Fixes issue with twig not accepting empty array values
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

trait CURL {

	public function makeConnection() {
	
		return false;
	
	}

}

?>