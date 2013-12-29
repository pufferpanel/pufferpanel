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
	    
	public function __construct($timeout = 1)
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
	
	public function getStatus($host, $port = 25565, $kill = false)
        {
        
        	    	
    		$this->socket = socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
    		socket_set_option($this->socket, SOL_SOCKET, SO_SNDTIMEO, array( 'sec' => $this->timeout, 'usec' => 0 ));
    		socket_set_option($this->socket, SOL_SOCKET, SO_RCVTIMEO, array( 'sec' => $this->timeout, 'usec' => 0 ));
    		@socket_connect($this->socket, $host, $port);
    		
    		$data = Pack('cccca*', hexdec(strlen($host)), 0, 0x04, strlen($host), $host).pack('nc', $port, 0x01);
    		
    		socket_send($this->socket, $data, strlen($data), 0);
    		socket_send($this->socket, "\x01\x00", 2, 0);
    		    		
    		if($this->readInt($this->socket) < 10 )
    			{
    		        
    		        socket_close($this->socket);
    		        return false;
    			
    			}
    		
    		socket_read($this->socket, 1);
    		$data = socket_read($this->socket, $this->readInt($this->socket), PHP_NORMAL_READ);
    		socket_close($this->socket);
    		
    		$data = json_decode($data, true);
    		
    		
    		$pList = array();
    		
    		die(print_r($data));
    		foreach($data['players']['sample'] as $id => $internal)
    			{
    				
    				$pList = array_merge($pList, array($internal['name']));
    			
    			}
    		
    		$this->info = array();
    		$this->info['motd'] = $data['description'];
    		$this->info['version'] = $data['version']['name'];
    		$this->info['players'] = $pList;
    		$this->info['maxplayers'] = $data['players']['max'];
    		
    	}
	
	public function data($value)
		{
		
			return array_key_exists($value, $this->info) ? $this->info[$value] : false;
		
		}
	
	private function readInt($socket)
		{
		
			$i = 0;
			$j = 0;
			
			while(true)
				{
			        $k = ord(socket_read($socket, 1));
			        
			        $i |= ($k & 0x7F) << $j++ * 7;
			        
			        if($j > 5)
				        exit('VarInt Is Too Big');
			        
			        if(($k & 0x80) != 128)
				        break;
				}
			
			return $i;
			
		}
		      
}

?>