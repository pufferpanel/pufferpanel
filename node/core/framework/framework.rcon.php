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

class rcon {
	    
	public function __construct($timeout = 3)
		{
		
			$this->timeout = $timeout;
		
		}
	
	public function online($host, $port)
		{
		
			if(!@fsockopen($host, $port, $errno, $errstr, $this->timeout))
				return false;
			else
		    	return true;
		    
		}
	
	public function getStatus($host, $port = 25565)
        {
        
        	$this->socket = @fsockopen('udp://'.$host, $port, $this->errNo, $this->errStr, $this->timeout);
            stream_set_timeout($this->socket, $this->timeout);
	        
	        $challengePack = pack('c*', 0xFE, 0xFD, 0x09, 0x01, 0x02, 0x03, 0x04);
	        fwrite($this->socket, $challengePack, strlen($challengePack));
	        
            $this->get(pack('N', substr(fread($this->socket, 2048), 5)));
            fclose($this->socket);
                
        }
	
	public function data($value)
		{
		
			return array_key_exists($value, $this->info) ? $this->info[$value] : false;	
		
		}
	
	private function get($challenge)
		{
		        
	        $sendData = pack('c*', 0xFE, 0xFD, 0x00, 0x01, 0x02, 0x03, 0x04).$challenge.pack('c*', 0x00, 0x00, 0x00, 0x00);
	        fwrite($this->socket, $sendData, strlen($sendData));
	        $data = substr(fread($this->socket, 2048), 5);
	
	        $this->info = array();
	
	        $data = explode("\x00\x00\x01player_\x00\x00", substr($data, 11));
	        $this->info['players'] = explode("\x00", substr($data[1], 0, -2));
	        $data = explode("\x00", $data[0]);
	
			foreach($data as $key => $value)
				{
				
					if($key & 1) {
					
						if($data[($key - 1)] == 'plugins')
							{
							
								list($software, $plugins) = explode(':', $value, 2);
								$this->info['software'] = $software;
								$this->info['plugins'] = str_replace(';', ',', str_replace(',', '-', trim($plugins)));
							
							}
						else
							$this->info[$data[($key - 1)]] = $value;
						
					
					}
				
				}
		
		}
		      
}

?>