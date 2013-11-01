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
        
        	$this->host = $host;
        	$this->port = $port;
        	
        	/*
        	 * Check Cache
        	 */
        	$fp = fopen(__DIR__.'/use_do17.php', 'r');
        	$content = fread($fp, filesize(__DIR__.'/use_do17.php'));
        	fclose($fp);
        	
        	$lines = explode("\n", $content);
        	
        		$do17 = false;
        		foreach($lines as $id => $value)
        			{
        			
        				if($id > 2){
        				
        					if($value == $this->host.':'.$this->port)
        						$do17 = true;
        				
        				}
        			
        			}
        	
			if($kill === false && $do17 === false)
				$this->do16($this->host, $this->port);
			else
				$this->do17($this->host, $this->port);
						                
        }
        
    public function do16($host, $port = 25565)
    	{
    	
    		$this->socket = @fsockopen('udp://'.$host, $port, $this->errNo, $this->errStr, $this->timeout);
    		stream_set_timeout($this->socket, $this->timeout);
    		
    		$challengePack = pack('c*', 0xFE, 0xFD, 0x09, 0x01, 0x02, 0x03, 0x04);
    		fwrite($this->socket, $challengePack, strlen($challengePack));
    		
    		$this->get(pack('N', substr(fread($this->socket, 2048), 5)));
    		
    		//Hide Code Bug because of this 1.7 Stuff
    		@fclose($this->socket);
    	
    	}
    	
    public function do17($host, $port = 25565)
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
    		foreach($data['players']['sample'] as $id => $internal)
    			{
    				
    				$pList = array_merge($pList, array($internal['name']));
    			
    			}
    		
    		$this->info = array();
    		$this->info['motd'] = $data['description'];
    		$this->info['software'] = 'Unavaliable due to Minecraft 1.7 Query Changes.';
    		$this->info['plugins'] = 'Unavaliable due to Minecraft 1.7 Query Changes. n/a';
    		$this->info['version'] = $data['version']['name'];
    		$this->info['players'] = $pList;
    		$this->info['maxplayers'] = $data['players']['max'];
    		
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
	       
	       	if(!empty($data)) {
	       	
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
					
			}else{
			
				if(!file_exists(__DIR__.'/use_do17.php')){
				
					$addContent = '<?php 
	exit();
?>
';
				
				}
				
				$fp = fopen(__DIR__.'/use_do17.php', 'a+');
				fwrite($fp, $addContent.$this->host.":".$this->port."\n");
				fclose($fp);
				
				self::getStatus($this->host, $this->port, true);
				
			}
		
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