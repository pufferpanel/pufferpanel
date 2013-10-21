<?php

class MCOnline {

	public function isOnline($server, $port){
	
		$fp = @fsockopen($server, $port, $errno, $errstr, 3);
		
			if(!$fp){
			
				return false;
				
			}else{
			
				return true;
			
			}
	
	}

}

?>