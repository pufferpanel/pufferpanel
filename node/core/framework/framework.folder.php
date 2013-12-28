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
 * PufferPanel Folder Management Class
 */

class files {

	public function getFolderSize($dir_name){
		
		$dir_size =0;
			
			if (is_dir($dir_name)) {
				
				if ($dh = opendir($dir_name)) {
					
					while (($file = readdir($dh)) !== false) {

						if($file !="." && $file != ".."){

							if(is_file($dir_name."/".$file)){

								$dir_size += filesize($dir_name."/".$file);

							}
	                             
							if(is_dir($dir_name."/".$file)){
								
								$dir_size +=  $this->getFolderSize($dir_name."/".$file);
							
							}
							
						}
					
					}
				
				}
				
			}
	    
	    if(isset($dh)){ closedir($dh); }
		return $dir_size;
	
	}
	
	public function formatSize($size) {

		$units = explode(',', 'B,KB,MB,GB,TB,PB');	
	    $mod = 1024;
	
	    for ($i = 0; $size > $mod; $i++) {
	        $size /= $mod;
	    }
	
	    $endIndex = strpos($size, ".")+3;
	
	    return substr($size, 0, $endIndex).' '.$units[$i];
	    
	}
	
	public function readLines($filename, $lines)
	{

		$r = '';
		$file = file($filename);
		
			if(count($file) < $lines){
			
				$lines = count($file);
				
			}
		
				for ($i = count($file)-$lines; $i < count($file); $i++) {
		  
		 			 $r .= $file[$i];
		
				}
				
		return $r;
	      
	}
	
	function last_lines($path, $line_count, $block_size = 512){
	    
	    if(!is_file($path)) {
	    	
	    	echo('Could not locate server.log!');
	    	return false;
	    	
	    }else{
	    
		    $lines = array();
			    $leftover = "";
		    $fh = fopen($path, 'r');
		    fseek($fh, 0, SEEK_END);
		    
		    do{
	
		        $can_read = $block_size;
		        if(ftell($fh) < $block_size){
		            $can_read = ftell($fh);
		        }
		
		        fseek($fh, -$can_read, SEEK_CUR);
		        $data = fread($fh, $can_read);
		        $data .= $leftover;
		        fseek($fh, -$can_read, SEEK_CUR);
		
		        $split_data = array_reverse(explode("\n", $data));
		        $new_lines = array_slice($split_data, 0, -1);
		        $lines = array_merge($lines, $new_lines);
		        $leftover = $split_data[count($split_data) - 1];
		        
		    }
		    while(count($lines) < $line_count && ftell($fh) != 0);
		    
		    if(ftell($fh) == 0){
		        $lines[] = $leftover;
		    }
		    
		    fclose($fh);
	
			/*
			 * Return in Readbale Format
			 */
			$a = array_reverse(array_slice($lines, 0, $line_count));
		    
		    	$output = '';
		    	$total = count($a);
		    	foreach($a as $id => $line){
		    	
		    		if($id == $total || $id == ($total - 1)){
		    			$output .= $line;
		    		}else{
		    			$output .= $line."\n";
		    		}
		    	
		    	}	    	
		    
		    return $output;
		    
		}
	    
	}
	
	function download($filename) { 
		
		$chunksize = 1*(1024*1024);
		$buffer = '';
		$handle = fopen($filename, 'rb');
		
			if ($handle === false){
			
				return false;
				
			}

			while (!feof($handle)){
			
				$buffer = fread($handle, $chunksize);
				print $buffer;
				
			}
			
		return fclose($handle);
		
	} 
	
	/*
	 * Broken Function
	 */
	public function ftp_rmdirr($path, $handle){
	    
	    if (!(@ssh2_sftp_rmdir($handle, $path) || @ssh2_sftp_unlink($handle, $path))){
	        
	        $list = scandir("ssh2.sftp://".$handle.$path);
	        
	        if (!empty($list)){
	            
	            foreach($list as $value){
					
					if($value != '.' && $value != '..'){
						
						if(!is_dir("ssh2.sftp://".$handle.$path.$value)){
							
							ssh2_sftp_unlink($handle, $path);
							echo $value." is not a Dir.\n";
							
						}else{
							
							ssh2_sftp_rmdir($handle, $value);
							echo $value." is a Dir.\n";
							
						}
						
					}
					
	            }
	            
	        }
	        
	    }
	        
	    @ssh2_sftp_rmdir($handle, $path);
	    
	}

}

?>